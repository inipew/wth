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

// Client wraps the GitHub client to provide custom configurations.
type Client struct {
	githubClient *github.Client
}

// NewClient creates and returns a new GitHub client with the specified timeout.
func NewClient(timeout time.Duration) *Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		githubClient: github.NewClient(httpClient),
	}
}

// GetLatestRelease fetches the latest release version from a GitHub repository.
// preRelease determines whether to fetch pre-releases (true) or stable releases (false).
func (c *Client) GetLatestRelease(ctx context.Context, repoOwner, repoName string, preRelease bool) (string, error) {
	releases, err := c.fetchReleases(ctx, repoOwner, repoName)
	if err != nil {
		return "", err
	}

	filteredReleases := filterReleases(releases, preRelease)
	if len(filteredReleases) == 0 {
		return "", c.noReleasesError(repoOwner, repoName, preRelease)
	}

	latestRelease := filteredReleases[0]
	return strings.TrimPrefix(*latestRelease.TagName, "v"), nil
}

// fetchReleases retrieves the releases from the specified GitHub repository.
func (c *Client) fetchReleases(ctx context.Context, repoOwner, repoName string) ([]*github.RepositoryRelease, error) {
	releases, _, err := c.githubClient.Repositories.ListReleases(ctx, repoOwner, repoName, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving releases for %s/%s: %w", repoOwner, repoName, err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", repoOwner, repoName)
	}

	return releases, nil
}

// noReleasesError constructs an error message based on the release type.
func (c *Client) noReleasesError(repoOwner, repoName string, preRelease bool) error {
	if preRelease {
		return fmt.Errorf("no pre-releases available for %s/%s", repoOwner, repoName)
	}
	return fmt.Errorf("no stable releases available for %s/%s", repoOwner, repoName)
}

// filterReleases filters and sorts releases based on whether they are pre-releases or stable releases.
func filterReleases(releases []*github.RepositoryRelease, preRelease bool) []*github.RepositoryRelease {
	var filteredReleases []*github.RepositoryRelease

	for _, release := range releases {
		if release.GetPublishedAt().IsZero() {
			continue // Skip releases without a published date
		}
		if release.GetPrerelease() == preRelease {
			filteredReleases = append(filteredReleases, release)
		}
	}

	// Sort releases by publication date in descending order
	sort.Slice(filteredReleases, func(i, j int) bool {
		return filteredReleases[i].GetPublishedAt().Time.After(filteredReleases[j].GetPublishedAt().Time)
	})

	return filteredReleases
}
