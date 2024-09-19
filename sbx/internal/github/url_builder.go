package github

import (
	"fmt"
	"runtime"
)

// URL format templates for different repository names
const (
	caddyFormat   = "caddy_%s_%s_%s.tar.gz"
	singBoxFormat = "sing-box-%s-%s-%s.tar.gz"
)

// Error messages
const (
	errEmptyInput      = "repoOwner, repoName, and version must be non-empty"
	errUnsupportedRepo = "unsupported repoName: %s"
)

// repositories maps repository names to their corresponding file name formats
var repositories = map[string]string{
	"caddy":     caddyFormat,
	"sing-box":  singBoxFormat,
}

// BuildDownloadURL constructs a download URL based on repository details, version, OS, and architecture.
//
// Parameters:
//   - repoOwner: The owner of the repository on GitHub.
//   - repoName: The name of the repository.
//   - version: The version of the release.
//
// Returns:
//   - The constructed download URL or an error if any parameter is invalid or unsupported.
func BuildDownloadURL(repoOwner, repoName, version string) (string, error) {
	if err := validateInputs(repoOwner, repoName, version); err != nil {
		return "", err
	}

	fileNameFormat, exists := repositories[repoName]
	if !exists {
		return "", fmt.Errorf(errUnsupportedRepo, repoName)
	}

	fileName := formatFileName(fileNameFormat, version)
	downloadURL := formatDownloadURL(repoOwner, repoName, version, fileName)

	return downloadURL, nil
}

// validateInputs checks if the provided parameters are valid.
func validateInputs(repoOwner, repoName, version string) error {
	if repoOwner == "" || repoName == "" || version == "" {
		return fmt.Errorf(errEmptyInput)
	}
	return nil
}

// formatFileName creates a file name based on the format and parameters.
func formatFileName(format, version string) string {
	return fmt.Sprintf(format, version, runtime.GOOS, runtime.GOARCH)
}

// formatDownloadURL constructs the complete download URL for the release.
func formatDownloadURL(repoOwner, repoName, version, fileName string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/v%s/%s", repoOwner, repoName, version, fileName)
}
