package fetch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/lets-cli/lets/internal/set"
)

var allowedContentTypes = set.NewSet(
	"text/plain",
	"text/yaml",
	"text/x-yaml",
	"application/yaml",
	"application/x-yaml",
)

var httpClient = newHTTPClient()

func newHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Keep Content-Length usable for known-size download progress bars.
	transport.DisableCompression = true

	// Timeout guards against hung connections when callers pass context.Background().
	// 5 minutes matches the previous per-request context timeout used by RemoteMixin downloads.
	return &http.Client{
		Timeout:   5 * 60 * time.Second,
		Transport: transport,
	}
}

type SourceKind string

const (
	SourceRemoteConfig SourceKind = "remote config"
	SourceRemoteMixin  SourceKind = "remote mixin"
)

type ProgressInfo struct {
	Kind       SourceKind
	URL        string
	TotalBytes int64
}

type ProgressObserver interface {
	Start(info ProgressInfo) ProgressTracker
}

type ProgressTracker interface {
	Add(n int64)
	Done(err error)
}

type downloadOptions struct {
	progress ProgressObserver
	kind     SourceKind
}

type Option func(*downloadOptions)

func WithProgress(kind SourceKind, progress ProgressObserver) Option {
	return func(opts *downloadOptions) {
		opts.kind = kind
		opts.progress = progress
	}
}

func newDownloadOptions(options []Option) downloadOptions {
	opts := downloadOptions{}
	for _, option := range options {
		option(&opts)
	}

	return opts
}

// Download fetches the content at url, validates the Content-Type is a YAML variant,
// and returns the raw bytes.
func Download(ctx context.Context, url string, options ...Option) ([]byte, error) {
	opts := newDownloadOptions(options)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no such file at: %s", url)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("network error for %s: %s", url, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")

	mediaType, _, parseErr := mime.ParseMediaType(contentType)
	if parseErr != nil {
		mediaType = contentType
	}

	if !allowedContentTypes.Contains(mediaType) {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	var tracker ProgressTracker
	if opts.progress != nil {
		tracker = opts.progress.Start(ProgressInfo{
			Kind:       opts.kind,
			URL:        url,
			TotalBytes: resp.ContentLength,
		})
	}

	data, err := readAll(resp.Body, tracker, resp.ContentLength)
	if tracker != nil {
		tracker.Done(err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}

func readAll(reader io.Reader, tracker ProgressTracker, contentLength int64) ([]byte, error) {
	var buf bytes.Buffer
	if contentLength > 0 && contentLength <= int64(^uint(0)>>1) {
		buf.Grow(int(contentLength))
	}

	if tracker == nil {
		_, err := io.Copy(&buf, reader)
		return buf.Bytes(), err
	}

	_, err := io.Copy(&buf, progressReader{reader: reader, tracker: tracker})

	return buf.Bytes(), err
}

type progressReader struct {
	reader  io.Reader
	tracker ProgressTracker
}

func (r progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.tracker.Add(int64(n))
	}

	return n, err
}
