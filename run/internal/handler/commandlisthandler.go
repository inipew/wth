package handler

import (
	"encoding/json"
	"net/http"
	cfg "run/internal/config"

	"github.com/sirupsen/logrus"
)

// ListCommandsResponse represents the structure of the API response for the list of commands.
type ListCommandsResponse struct {
	Commands []cfg.Command `json:"commands"`
}

// ListCommandsHandler handles HTTP requests to get the list of commands.
func ListCommandsHandler(config *cfg.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := ListCommandsResponse{Commands: config.Commands}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// findCommandByValue searches for a command by its value in the Config.
func findCommandByValue(cfg *cfg.Config, value string) (*cfg.Command, bool) {
	for _, cmd := range cfg.Commands {
		if cmd.Value == value {
			// Kembalikan salinan dari cmd
			cmdCopy := cmd // Membuat salinan
			return &cmdCopy, true
		}
	}
	return nil, false
}