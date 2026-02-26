package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const (
	maxRequestBodySize  = 4096
	maxExternalRespSize = 4096
	httpClientTimeout   = 10 * time.Second
	inviteMaxAge        = 1800
	inviteMaxUses       = 1
)

var httpClient = &http.Client{Timeout: httpClientTimeout}

type inviteRequest struct {
	Captcha string `json:"captcha"`
}

type captchaVerification struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes,omitempty"`
}

type discordInviteResponse struct {
	Code      string `json:"code"`
	ExpiresAt string `json:"expires_at"`
}

func (s *Server) handleDiscordInvite(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxRequestBodySize {
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	var req inviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		slog.Debug("failed to decode invite request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Captcha == "" {
		http.Error(w, "captcha token is required", http.StatusBadRequest)
		return
	}

	verification, err := verifyCaptcha(req.Captcha, s.cfg.RecaptchaSecret)
	if err != nil {
		slog.Error("captcha verification failed", "error", err)
		http.Error(w, "captcha verification failed", http.StatusInternalServerError)
		return
	}

	if !verification.Success {
		slog.Warn("captcha verification rejected", "errors", verification.ErrorCodes)
		http.Error(w, "captcha verification failed", http.StatusBadRequest)
		return
	}

	invite, err := createDiscordInvite(s.cfg.DiscordChannelID, s.cfg.DiscordToken)
	if err != nil {
		slog.Error("discord invite creation failed", "error", err)
		http.Error(w, "invite creation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(invite)
}

func verifyCaptcha(token, secret string) (*captchaVerification, error) {
	data := url.Values{
		"response": {token},
		"secret":   {secret},
	}

	resp, err := httpClient.PostForm("https://www.google.com/recaptcha/api/siteverify", data)
	if err != nil {
		return nil, fmt.Errorf("recaptcha request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("recaptcha API returned status %d", resp.StatusCode)
	}

	var result captchaVerification
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxExternalRespSize)).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode recaptcha response: %w", err)
	}

	return &result, nil
}

func createDiscordInvite(channelID, token string) (*discordInviteResponse, error) {
	body := map[string]interface{}{
		"max_age":  inviteMaxAge,
		"max_uses": inviteMaxUses,
		"unique":   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/invites", channelID)
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "RotkiBot")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discord request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxExternalRespSize))
		slog.Debug("discord API error response", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("discord API returned status %d", resp.StatusCode)
	}

	var invite discordInviteResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxExternalRespSize)).Decode(&invite); err != nil {
		return nil, fmt.Errorf("failed to decode discord response: %w", err)
	}

	return &invite, nil
}
