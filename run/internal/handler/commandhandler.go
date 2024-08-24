package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	cfg "run/internal/config"
	"run/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
)

// CommandResponse represents the structure of the API response for command output.
type CommandResponse struct {
	Output string `json:"output"`
}

// APIHandler handles HTTP requests for executing commands.
func CommandHandler(config *cfg.Config, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		value := r.URL.Query().Get("value")
		if value == "" {
			http.Error(w, "Missing 'value' parameter", http.StatusBadRequest)
			return
		}

		var command *cfg.Command
		if value == "custom" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form data", http.StatusBadRequest)
				return
			}
			customCommand := r.FormValue("custom_command")
			if customCommand == "" {
				http.Error(w, "Missing 'custom_command' parameter", http.StatusBadRequest)
				return
			}
			command = &cfg.Command{Command: customCommand}
		} else {
			// Pastikan fungsi ini ada dan mengembalikan command yang sesuai
			var found bool
			command, found = findCommandByValue(config, value) // command adalah *Command
			if !found {
				http.Error(w, fmt.Sprintf("Command with Value '%s' not found", value), http.StatusNotFound)
				return
			}

			// Jika command adalah pointer, Anda harus dereference saat mengakses properti
			command.Command = cfg.ReplacePaths(command.Command, config.Paths) // Pastikan ini ada
		}

		output, err := utils.RunCommand(ctx, command.Command, timeout)
		if err != nil {
			logrus.Errorf("Failed to execute command '%s': %v", command.Command, err)
			http.Error(w, "Failed to execute command", http.StatusInternalServerError)
			return
		}

		response := CommandResponse{Output: output}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}



