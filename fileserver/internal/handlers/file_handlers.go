package handlers

import (
	"fileserver/internal/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileInfo stores information about files and directories
type FileInfo struct {
	Name          string
	Path          string
	IsDir         bool
	FileSize      int64
	FormattedSize string
	LastModified  string
	IsEditable    bool
}

// CreateUploadDir creates the directory for file uploads if it doesn't exist
func CreateUploadDir() error {
    return os.MkdirAll(utils.UploadPath, os.ModePerm)
}

func IsFileEditable(filename string) bool {
	editableExtensions := map[string]bool{
		".txt":  true,
		".md":   true,
		".log":  true,
		".html": true,
		".css":  true,
		".js":   true,
		".json": true,
		".xml":  true,
		".yml":  true,
		".yaml": true,
		".php":  true,
		".py":   true,
		".rb":   true,
		".java": true,
		".cpp":  true,
		".c":    true,
		".cs":   true,
		".go":   true,
		".rs":   true,
		".sh":   true,
		".bat":  true,
		".ini":  true,
		".cfg":  true,
		".conf": true,
		".env":  true,
		".pl":   true,
		".ps1":  true,
		".lua":  true,
		".r":    true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return editableExtensions[ext]
}

// sortFileInfos sorts the file information by directory and name
func sortFileInfos(fileInfos []FileInfo) {
	sort.SliceStable(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir != fileInfos[j].IsDir {
			return fileInfos[i].IsDir
		}
		return strings.ToLower(fileInfos[i].Name) < strings.ToLower(fileInfos[j].Name)
	})
}

// formatFileSize formats the file size for display
func formatFileSize(info os.FileInfo, file os.DirEntry) string {
	if file.IsDir() {
		return ""
	}
	size := utils.ByteSize(info.Size())
	return size.String()
}