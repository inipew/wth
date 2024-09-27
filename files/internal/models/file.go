package models

// DirectoryInfo menyimpan informasi direktori
type DirectoryInfo struct {
	CurrentPath   string     `json:"current_path,omitempty"`
	PreviousPath  string     `json:"previous_path,omitempty"`
	Files         []FileInfo `json:"files"`
}

// FileInfo menyimpan informasi file
type FileInfo struct {
	Name           string `json:"name,omitempty"`
	Path           string `json:"path,omitempty"`
	IsDir          bool   `json:"is_dir"`
	FileSize       string `json:"file_size,omitempty"`
	Size		   int64  `json:"size"`
	LastModified   string `json:"last_modified,omitempty"`
	IsEditable     bool   `json:"is_editable"`
	Permissions    string `json:"permissions,omitempty"`
	FileType       string `json:"file_type,omitempty"`
	Owner          string `json:"owner,omitempty"`
	Group          string `json:"group,omitempty"`
	CreationDate   string `json:"creation_date,omitempty"`
}

type ArchiveInfo struct{
	Name	string 				`json:"name"`
	Path    string 				`json:"path"`
	Files	[]ArchiveFileInfo	`json:"files"`
}

type ArchiveFileInfo struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	IsDir          bool   `json:"is_dir"`
	FileSize       string `json:"file_size"`
	LastModified   string `json:"last_modified"`
	CreationDate   string `json:"creation_date,omitempty"`
}
