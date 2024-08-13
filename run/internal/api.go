package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// APIHandler handles HTTP requests for executing commands.
func CommandHandler(config *Config, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Mengizinkan semua domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Metode yang diizinkan
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Header yang diizinkan

		// Tangani permintaan preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // Mengirim respons 204 No Content
			return
		}
		
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		value := r.URL.Query().Get("value")
		if value == "" {
			http.Error(w, "Missing 'value' parameter", http.StatusBadRequest)
			return
		}

		var command *Command
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
			command = &Command{Command: customCommand}
		} else {
			// Pastikan fungsi ini ada dan mengembalikan command yang sesuai
			var found bool
			command, found = findCommandByValue(config, value) // command adalah *Command
			if !found {
				http.Error(w, fmt.Sprintf("Command with Value '%s' not found", value), http.StatusNotFound)
				return
			}

			// Jika command adalah pointer, Anda harus dereference saat mengakses properti
			command.Command = ReplacePaths(command.Command, config.Paths) // Pastikan ini ada
		}

		output, err := RunCommand(ctx, command.Command, timeout)
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

// ListCommandsHandler handles HTTP requests to get the list of commands.
func ListCommandsHandler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := ListCommandsResponse{Commands: config.Commands}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// ListCommandsResponse represents the structure of the API response for the list of commands.
type ListCommandsResponse struct {
	Commands []Command `json:"commands"`
}

// CommandResponse represents the structure of the API response for command output.
type CommandResponse struct {
	Output string `json:"output"`
}

// findCommandByValue searches for a command by its value in the Config.
func findCommandByValue(cfg *Config, value string) (*Command, bool) {
	for _, cmd := range cfg.Commands {
		if cmd.Value == value {
			// Kembalikan salinan dari cmd
			cmdCopy := cmd // Membuat salinan
			return &cmdCopy, true
		}
	}
	return nil, false
}

