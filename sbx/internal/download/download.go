package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

// Client is a wrapper around http.Client that provides configurable download capabilities.
type Client struct {
	HTTPClient       *http.Client
	RetryCount       int
	RetryDelay       time.Duration
	ConcurrentChunks int
	ChunkSize        int64
}

// NewClient creates a new Client with customizable settings.
func NewClient(timeout time.Duration, retryCount int, retryDelay time.Duration, concurrentChunks int, chunkSize int64) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		RetryCount:       retryCount,
		RetryDelay:       retryDelay,
		ConcurrentChunks: concurrentChunks,
		ChunkSize:        chunkSize,
	}
}

// ProgressBar manages the progress bar display.
type ProgressBar struct {
	bar  *progressbar.ProgressBar
	mu   sync.Mutex
	size int64
}

// NewProgressBar creates and initializes a new progress bar.
func NewProgressBar(totalSize int64) *ProgressBar {
	bar := progressbar.NewOptions64(totalSize,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan]Downloading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return &ProgressBar{bar: bar, size: totalSize}
}

// Add increments the progress by n bytes.
func (p *ProgressBar) Add(n int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bar.Add64(n)
}

// DownloadFile downloads a file from the given URL to the specified filepath, with retry and concurrent chunk logic.
func (c *Client) DownloadFile(ctx context.Context, url, filePath string) error {
	if err := ensureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fileSize, err := c.getFileSize(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to get file size: %w", err)
	}

	progress := NewProgressBar(fileSize)
	chunks := c.calculateChunks(fileSize)

	tempDir, err := os.MkdirTemp("", "download-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	eg, ctx := errgroup.WithContext(ctx)
	for i, chunk := range chunks {
		i, chunk := i, chunk
		eg.Go(func() error {
			return c.downloadChunk(ctx, url, tempDir, i, chunk, progress)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to download all chunks: %w", err)
	}

	return c.mergeChunks(tempDir, filePath, len(chunks))
}

func (c *Client) calculateChunks(fileSize int64) []struct{ Start, End int64 } {
	chunks := make([]struct{ Start, End int64 }, 0)
	for start := int64(0); start < fileSize; start += c.ChunkSize {
		end := start + c.ChunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}
		chunks = append(chunks, struct{ Start, End int64 }{start, end})
	}
	return chunks
}

func (c *Client) downloadChunk(ctx context.Context, url, tempDir string, index int, chunk struct{ Start, End int64 }, progress *ProgressBar) error {
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", index))
	for attempt := 1; attempt <= c.RetryCount; attempt++ {
		err := c.tryDownloadChunk(ctx, url, chunkPath, chunk, progress)
		if err == nil {
			return nil
		}
		if attempt == c.RetryCount {
			return fmt.Errorf("failed to download chunk after %d attempts: %w", c.RetryCount, err)
		}
		time.Sleep(c.RetryDelay)
	}
	return nil
}

func (c *Client) tryDownloadChunk(ctx context.Context, url, chunkPath string, chunk struct{ Start, End int64 }, progress *ProgressBar) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", chunk.Start, chunk.End))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading chunk: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server returned unexpected status: %s", resp.Status)
	}

	out, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("error creating chunk file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, io.TeeReader(resp.Body, &progressWriter{progress: progress}))
	if err != nil {
		return fmt.Errorf("error writing chunk: %w", err)
	}

	if written != chunk.End-chunk.Start+1 {
		return fmt.Errorf("incomplete chunk download")
	}

	return nil
}

func (c *Client) mergeChunks(tempDir, filePath string, numChunks int) error {
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	for i := 0; i < numChunks; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", i))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("error opening chunk file: %w", err)
		}
		_, err = io.Copy(out, chunk)
		chunk.Close()
		if err != nil {
			return fmt.Errorf("error merging chunk: %w", err)
		}
	}

	return nil
}

func (c *Client) getFileSize(ctx context.Context, url string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned: %s", resp.Status)
	}

	size := resp.ContentLength
	if size <= 0 {
		return 0, errors.New("could not determine file size")
	}

	return size, nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

type progressWriter struct {
	progress *ProgressBar
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.progress.Add(int64(n))
	return n, nil
}