package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/rotki/discord-captcha/internal/config"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	cfg := &config.Config{
		Port:             "4000",
		RecaptchaSecret:  "test-secret",
		DiscordChannelID: "test-channel",
		DiscordToken:     "test-token",
	}

	staticFS := fstest.MapFS{
		"index.html":    &fstest.MapFile{Data: []byte("<html>test</html>")},
		"assets/app.js": &fstest.MapFile{Data: []byte("console.log('test')")},
	}

	return NewServer(cfg, nil, staticFS)
}

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Ok" {
		t.Errorf("expected body 'Ok', got %q", string(body))
	}
}

func TestStaticFileServing(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest("GET", "/assets/app.js", http.NoBody)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "console.log('test')" {
		t.Errorf("unexpected body: %q", string(body))
	}
}

func TestSPAFallback(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest("GET", "/some/unknown/route", http.NoBody)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for SPA fallback, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "<html>test</html>" {
		t.Errorf("expected index.html content, got %q", string(body))
	}
}

func TestSecurityHeaders(t *testing.T) {
	srv := newTestServer(t)
	handler := securityHeadersMiddleware(cspMiddleware(srv.mux))

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	expected := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Permissions-Policy":     "camera=(), microphone=(), geolocation=()",
	}

	for header, want := range expected {
		got := resp.Header.Get(header)
		if got != want {
			t.Errorf("header %s: expected %q, got %q", header, want, got)
		}
	}
}

func TestCSPHeader(t *testing.T) {
	srv := newTestServer(t)
	handler := cspMiddleware(srv.mux)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	csp := resp.Header.Get("Content-Security-Policy")
	if csp == "" {
		t.Error("expected Content-Security-Policy header to be set")
	}

	// Verify a few key directives are present
	directives := []string{"default-src", "script-src", "frame-src", "style-src"}
	for _, d := range directives {
		if !containsDirective(csp, d) {
			t.Errorf("CSP missing directive: %s", d)
		}
	}
}

func containsDirective(csp, directive string) bool {
	for _, part := range splitCSP(csp) {
		if part != "" && part[:len(directive)] == directive {
			return true
		}
	}
	return false
}

func splitCSP(csp string) []string {
	var parts []string
	current := ""
	for _, c := range csp {
		if c == ';' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func TestShouldGzip(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/index.html", true},
		{"/app.js", true},
		{"/style.css", true},
		{"/api/discord-invite", true},
		{"/image.png", false},
		{"/image.jpg", false},
		{"/image.jpeg", false},
		{"/image.webp", false},
		{"/font.woff", false},
		{"/font.woff2", false},
		{"/archive.gz", false},
		{"/archive.br", false},
		{"/archive.zip", false},
		{"/IMAGE.PNG", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := shouldGzip(tt.path)
			if got != tt.want {
				t.Errorf("shouldGzip(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestInviteEndpointMissingBody(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest("POST", "/api/discord-invite", http.NoBody)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestInviteEndpointEmptyCaptcha(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest("POST", "/api/discord-invite", strings.NewReader(`{"captcha":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
