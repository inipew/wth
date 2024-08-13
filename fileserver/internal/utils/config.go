package utils

import (
	"os"
	"strconv"
)

const defaultMaxUploadSize = 10 << 20 // 10 MB

// MaxUploadSize returns the maximum upload size from environment variable or default value
func MaxUploadSize() int64 {
    sizeStr := os.Getenv("MAX_UPLOAD_SIZE")
    if sizeStr != "" {
        if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
            return size
        }
    }
    return defaultMaxUploadSize
}
