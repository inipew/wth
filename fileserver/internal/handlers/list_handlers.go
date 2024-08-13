package handlers

import (
	"fileserver/templates"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// DirectoryListing contains data for rendering the directory view
type DirectoryListing struct {
	CurrentPath string
	PrevPath    string
	Breadcrumbs []Breadcrumb
	Files       []FileInfo
}

// Breadcrumb holds the breadcrumb navigation data
type Breadcrumb struct {
	Name string
	Path string
}

// IndexFileManagerHandler handles requests to manage files in a directory
func IndexFileManagerHandler(w http.ResponseWriter, r *http.Request) {
	dirPath, err := getDirectoryPath(r)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	fileInfos, err := prepareFileInfo(files, dirPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	breadcrumbs := generateBreadcrumbs(dirPath)
	prevPath := getPreviousPath(dirPath)

	data := DirectoryListing{
		CurrentPath: filepath.ToSlash(filepath.Clean(dirPath)),
		PrevPath:    filepath.ToSlash(prevPath),
		Breadcrumbs: breadcrumbs,
		Files:       fileInfos,
	}

	if err := renderTemplate(w, "list_directory.html", data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}

// getDirectoryPath retrieves and validates the directory path from the request
func getDirectoryPath(r *http.Request) (string, error) {
	dir := r.URL.Query().Get("dir")
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

// prepareFileInfo prepares a slice of FileInfo from the given directory entries
func prepareFileInfo(files []os.DirEntry, dirPath string) ([]FileInfo, error) {
	var fileInfos []FileInfo

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		formattedSize := formatFileSize(info, file)
		lastModified := info.ModTime().Format("2006-01-02 15:04:05")

		fileInfos = append(fileInfos, FileInfo{
			Name:          file.Name(),
			Path:          filepath.ToSlash(filepath.Join(dirPath, file.Name())),
			IsDir:         file.IsDir(),
			FileSize:      info.Size(),
			FormattedSize: formattedSize,
			LastModified:  lastModified,
			IsEditable:    IsFileEditable(file.Name()),
		})
	}

	sortFileInfos(fileInfos)
	return fileInfos, nil
}

// generateBreadcrumbs generates breadcrumbs for the current directory path
func generateBreadcrumbs(path string) []Breadcrumb {
	var breadcrumbs []Breadcrumb
	breadcrumbs = append(breadcrumbs, Breadcrumb{Name: "", Path: ""})

	if path == "" {
		return breadcrumbs
	}

	parts := strings.Split(filepath.ToSlash(path), "/")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		breadcrumbPath := strings.Join(parts[:i+1], "/")
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Name: parts[i],
			Path: breadcrumbPath,
		})
	}
	return breadcrumbs
}

// getPreviousPath returns the previous directory path
func getPreviousPath(currentPath string) string {
	prevPath := filepath.Dir(filepath.Clean(currentPath))
	if prevPath == currentPath {
		return ""
	}
	return prevPath
}

// handleError sends an HTTP error response
func handleError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}

// renderTemplate renders the specified template with data
func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
		"ext":     filepath.Ext,
	}

	tmpl, err := templates.GetTemplate(templateName, funcMap)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
