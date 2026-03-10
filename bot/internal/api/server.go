package api

import (
	"compress/gzip"
	"context"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rotki/discord-captcha/internal/config"
	"github.com/rotki/discord-captcha/internal/store"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
)

type Server struct {
	cfg    *config.Config
	store  store.InviteStore
	mux    *http.ServeMux
	server *http.Server
}

func NewServer(cfg *config.Config, s store.InviteStore, staticFS fs.FS) *Server {
	srv := &Server{
		cfg:   cfg,
		store: s,
		mux:   http.NewServeMux(),
	}

	srv.mux.HandleFunc("POST /api/discord-invite", srv.handleDiscordInvite)
	srv.mux.HandleFunc("GET /health", srv.handleHealth)

	fileServer := http.FileServerFS(staticFS)
	srv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists in the static FS
		f, err := staticFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			_ = f.Close()
			// Vite hashed assets are immutable — cache forever
			if strings.HasPrefix(path, "/assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for non-file routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	handler := accessLogMiddleware(securityHeadersMiddleware(cspMiddleware(gzipMiddleware(srv.mux))))

	srv.server = &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	return srv
}

func (s *Server) ListenAndServe() error {
	slog.Info("HTTP server starting", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// First entry is the original client IP
		if i := strings.IndexByte(xff, ','); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// Skip logging internal health checks from loopback
		if r.URL.Path == "/health" && isLoopback(realIP(r)) {
			return
		}

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start).String(),
			"remote", realIP(r),
		)
	})
}

func isLoopback(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil && parsed.IsLoopback()
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

func cspMiddleware(next http.Handler) http.Handler {
	csp := buildCSP()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", csp)
		next.ServeHTTP(w, r)
	})
}

func buildCSP() string {
	policies := map[string][]string{
		"base-uri":                {"'self'"},
		"block-all-mixed-content": {},
		"child-src":               {"'none'"},
		"connect-src":             {"'self'", "www.google.com/recaptcha/"},
		"default-src":             {"'self'"},
		"font-src":                {"'self'", "fonts.gstatic.com"},
		"form-action":             {"'self'"},
		"frame-ancestors":         {"'self'"},
		"frame-src":               {"*.recaptcha.net", "recaptcha.net", "https://www.google.com/recaptcha/", "https://recaptcha.google.com"},
		"img-src":                 {"'self'", "'unsafe-inline'", "data:", "www.gstatic.com/recaptcha"},
		"object-src":              {"'none'"},
		"script-src":              {"'self'", "'unsafe-inline'", "'unsafe-eval'", "https://www.recaptcha.net", "https://recaptcha.net", "https://www.gstatic.com/recaptcha/", "https://www.gstatic.cn/recaptcha/", "https://www.google.com/recaptcha/"},
		"script-src-elem":         {"'self'", "'unsafe-inline'", "www.google.com/recaptcha/", "www.gstatic.com/recaptcha/"},
		"style-src":               {"'self'", "'unsafe-inline'", "fonts.googleapis.com"},
		"worker-src":              {"'self'", "www.recaptcha.net"},
	}

	// Use a fixed order for deterministic output
	order := []string{
		"base-uri", "block-all-mixed-content", "child-src", "connect-src",
		"default-src", "font-src", "form-action", "frame-ancestors",
		"frame-src", "img-src", "object-src", "script-src",
		"script-src-elem", "style-src", "worker-src",
	}

	var sb strings.Builder
	for _, key := range order {
		values := policies[key]
		sb.WriteString(key)
		for _, v := range values {
			sb.WriteByte(' ')
			sb.WriteString(v)
		}
		sb.WriteByte(';')
	}

	return sb.String()
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func shouldGzip(path string) bool {
	skip := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".woff", ".woff2", ".gz", ".br", ".zip"}
	lower := strings.ToLower(path)
	for _, ext := range skip {
		if strings.HasSuffix(lower, ext) {
			return false
		}
	}
	return true
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || !shouldGzip(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer func() { _ = gz.Close() }()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
