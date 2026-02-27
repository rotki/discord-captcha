package api

import (
	"encoding/json"
	"net/http"

	"github.com/rotki/discord-captcha/internal/version"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version.Version,
		"git_sha": version.GitSHA,
	})
}
