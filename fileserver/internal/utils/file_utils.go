package utils

import (
	"os"
	"path/filepath"
	"strings"
)

const UploadPath = "./uploads"
const maxPathDepth = 10


func IsValidPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return !strings.Contains(absPath, "..") && len(strings.Split(absPath, string(os.PathSeparator))) <= maxPathDepth
}