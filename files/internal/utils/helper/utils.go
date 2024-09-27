package helper

import (
	"files/internal/models"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// formatFileSize memformat ukuran file untuk tampilan
func FormatFileSize(info os.FileInfo, file os.DirEntry) string {
	if file.IsDir() {
		return ""
	}
	size := ByteSize(info.Size())
	return size.String()
}

func IsText(filename string) bool {
	// Periksa apakah path adalah direktori
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fileInfo.IsDir() {
		return false
	}

	// Buka file
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	// Tentukan ukuran buffer
	const bufferSize = 1024
	buf := make([]byte, bufferSize)

	// Baca data dari file
	n, err := file.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return false
	}

	// Periksa apakah data tersebut adalah teks
	for _, b := range buf[:n] {
		if b > 127 && !unicode.IsPrint(rune(b)) && b != 0x0A && b != 0x0D && b != 0x09 {
			return false
		}
	}
	return true
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

// func IsValidPath(baseDir, newPath string) bool {
// 	// Dapatkan path absolut dari baseDir
// 	absBaseDir, err := filepath.Abs(baseDir)
// 	if err != nil {
// 		log.Logger.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
// 		return false
// 	}

// 	// Dapatkan path absolut dari newPath
// 	absNewPath, err := filepath.Abs(newPath)
// 	if err != nil {
// 		log.Logger.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
// 		return false
// 	}

// 	// Periksa apakah absNewPath adalah subdirektori dari absBaseDir
// 	if !strings.HasPrefix(absNewPath, absBaseDir) {
// 		return false
// 	}

// 	// Periksa kedalaman path
// 	pathDepth := len(strings.Split(absNewPath[len(absBaseDir):], string(os.PathSeparator)))
// 	return pathDepth <= maxPathDepth
// }

// sortFileInfos mengurutkan slice FileInfo berdasarkan direktori dan nama
func SortFileInfos(fileInfos []models.FileInfo) {
	sort.SliceStable(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir != fileInfos[j].IsDir {
			return fileInfos[i].IsDir
		}
		return strings.ToLower(fileInfos[i].Name) < strings.ToLower(fileInfos[j].Name)
	})
}