package models

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
