package models

// Response digunakan untuk format balasan API
type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
