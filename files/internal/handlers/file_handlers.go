package handlers

import (
	"encoding/json"
	"files/internal/models"
	"files/internal/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Response digunakan untuk format balasan API
type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// FileInfo menyimpan informasi file
type FileInfo struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	IsDir          bool   `json:"is_dir"`
	FileSize       int64  `json:"file_size"`
	FormattedSize  string `json:"formatted_size"`
	LastModified   string `json:"last_modified"`
	IsEditable     bool   `json:"is_editable"`
	Permissions    string `json:"permissions,omitempty"`
	FileType       string `json:"file_type,omitempty"`
	Owner          string `json:"owner,omitempty"`
	Group          string `json:"group,omitempty"`
	CreationDate   string `json:"creation_date,omitempty"`
}

// DirectoryInfo menyimpan informasi direktori
type DirectoryInfo struct {
	CurrentPath   string     `json:"current_path"`
	PreviousPath  string     `json:"previous_path,omitempty"`
	Files         []FileInfo `json:"files"`
}

// FileHandler menangani permintaan untuk daftar file dalam direktori
func FileHandler(w http.ResponseWriter, r *http.Request) {
	currentPath, err := getDirectoryPath(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsValidPath(currentPath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(currentPath)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to read directory")
		return
	}

	fileInfos, err := prepareFileInfo(files, currentPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get file info")
		return
	}

	previousPath := filepath.Dir(currentPath)
	if currentPath == previousPath {
		previousPath = ""
	}

	respondWithJSON(w, http.StatusOK, models.DirectoryInfo{
		CurrentPath:  currentPath,
		PreviousPath: previousPath,
		Files:        fileInfos,
	})
}

// getDirectoryPath mendapatkan dan memvalidasi path direktori dari permintaan
func getDirectoryPath(r *http.Request) (string, error) {
	dir := r.URL.Query().Get("path")
	if dir == "" {
		return os.Getwd()
	}

	decodedDir, err := url.QueryUnescape(dir)
	if err != nil {
		return "", err
	}
	dirClean := filepath.Clean(decodedDir)

	absDirPath, err := filepath.Abs(dirClean)
	if err != nil {
		return "", err
	}

	return absDirPath, nil
}

// prepareFileInfo menyiapkan slice FileInfo dari entri direktori yang diberikan
func prepareFileInfo(files []os.DirEntry, dirPath string) ([]models.FileInfo, error) {
	var fileInfos []models.FileInfo

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info: %v", err)
			continue
		}

		formattedSize := utils.FormatFileSize(info, file)
		lastModified := info.ModTime().Format("2006-01-02 15:04:05")
		permissions := getFilePermissions(info)
		fileType := getFileType(info)
		owner, group, _ := getFileOwnerGroup(filepath.Join(dirPath, file.Name()))
		creationDate := getCreationDate(info)

		fileInfos = append(fileInfos, models.FileInfo{
			Name:          file.Name(),
			Path:          filepath.ToSlash(filepath.Join(dirPath, file.Name())),
			IsDir:         file.IsDir(),
			FileSize:      info.Size(),
			FormattedSize: formattedSize,
			LastModified:  lastModified,
			IsEditable:    utils.IsFileEditable(file.Name()),
			Permissions:   permissions,
			FileType:      fileType,
			Owner:         owner,
			Group:         group,
			CreationDate:  creationDate,
		})
	}
	utils.SortFileInfos(fileInfos)
	return fileInfos, nil
}

// respondWithJSON mengirimkan balasan dalam format JSON
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// respondWithError mengirimkan balasan kesalahan dalam format JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{Message: message})
}

// getFilePermissions mengembalikan string yang mewakili hak akses file
func getFilePermissions(fi os.FileInfo) string {
	return fi.Mode().String()
}

// getFileType mengembalikan jenis file
func getFileType(fi os.FileInfo) string {
	if fi.IsDir() {
		return "directory"
	}
	return "file"
}

// getFileOwnerGroup mendapatkan pemilik dan grup file untuk sistem Unix
func getFileOwnerGroup(path string) (string, string, error) {
	userCmd := exec.Command("stat", "-c", "%U", path)
	groupCmd := exec.Command("stat", "-c", "%G", path)

	userBytes, err := userCmd.Output()
	if err != nil {
		return "", "", err
	}
	groupBytes, err := groupCmd.Output()
	if err != nil {
		return "", "", err
	}

	return strings.TrimSpace(string(userBytes)), strings.TrimSpace(string(groupBytes)), nil
}

// getCreationDate mendapatkan tanggal pembuatan file
func getCreationDate(fi os.FileInfo) string {
	return fi.ModTime().Format(time.RFC3339) // Pseudonym creation date for Unix files
}