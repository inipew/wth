package utils

import (
	"files/internal/models"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// formatFileSize memformat ukuran file untuk tampilan
func FormatFileSize(info os.FileInfo, file os.DirEntry) string {
	if file.IsDir() {
		return ""
	}
	size := ByteSize(info.Size())
	return size.String()
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

const maxPathDepth = 10

// isValidPath memvalidasi apakah path aman dan tidak melewati kedalaman maksimum
func IsValidPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return !strings.Contains(absPath, "..") && len(strings.Split(absPath, string(os.PathSeparator))) <= maxPathDepth
}

// sortFileInfos mengurutkan slice FileInfo berdasarkan direktori dan nama
func SortFileInfos(fileInfos []models.FileInfo) {
	sort.SliceStable(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir != fileInfos[j].IsDir {
			return fileInfos[i].IsDir
		}
		return strings.ToLower(fileInfos[i].Name) < strings.ToLower(fileInfos[j].Name)
	})
}