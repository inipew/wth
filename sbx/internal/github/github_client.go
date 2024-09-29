package github

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v63/github"
)

// ReleaseManager handles operations related to GitHub releases.
type ReleaseManager struct {
	client *github.Client
}

// NewReleaseManager creates and returns a new ReleaseManager with the specified timeout.
func NewReleaseManager(timeout time.Duration) *ReleaseManager {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &ReleaseManager{
		client: github.NewClient(httpClient),
	}
}

// GetLatestRelease fetches the latest release version from a GitHub repository.
func (rm *ReleaseManager) GetLatestRelease(ctx context.Context, repoOwner, repoName string, releaseType bool) (string, error) {
	releases, err := rm.fetchReleases(ctx, repoOwner, repoName)
	if err != nil {
		return "", err
	}

	filteredReleases := rm.filterReleases(releases, releaseType)
	if len(filteredReleases) == 0 {
		return "", rm.noReleasesError(repoOwner, repoName, releaseType)
	}

	latestRelease := filteredReleases[0]
	return strings.TrimPrefix(*latestRelease.TagName, "v"), nil
}

// fetchReleases retrieves the releases from the specified GitHub repository.
func (rm *ReleaseManager) fetchReleases(ctx context.Context, repoOwner, repoName string) ([]*github.RepositoryRelease, error) {
	releases, _, err := rm.client.Repositories.ListReleases(ctx, repoOwner, repoName, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving releases for %s/%s: %w", repoOwner, repoName, err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", repoOwner, repoName)
	}

	return releases, nil
}

// noReleasesError constructs an error message based on the release type.
func (rm *ReleaseManager) noReleasesError(repoOwner, repoName string, releaseType bool) error {
	releaseTypeStr := "stable releases"
	if releaseType {
		releaseTypeStr = "pre-releases"
	}
	return fmt.Errorf("no %s available for %s/%s", releaseTypeStr, repoOwner, repoName)
}

// filterReleases filters and sorts releases based on the specified release type.
func (rm *ReleaseManager) filterReleases(releases []*github.RepositoryRelease, releaseType bool) []*github.RepositoryRelease {
	var filteredReleases []*github.RepositoryRelease

	for _, release := range releases {
		if release.GetPublishedAt().IsZero() {
			continue // Skip releases without a published date
		}
		if (releaseType) == release.GetPrerelease() {
			filteredReleases = append(filteredReleases, release)
		}
	}

	// Sort releases by publication date in descending order
	sort.Slice(filteredReleases, func(i, j int) bool {
		return filteredReleases[i].GetPublishedAt().Time.After(filteredReleases[j].GetPublishedAt().Time)
	})

	return filteredReleases
}