package fetch

import (
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

// httpClient backstop timeout guards against hung connections when callers pass context.Background().
// 5 minutes matches the previous per-request context timeout used by RemoteMixin downloads.
var httpClient = &http.Client{
	Timeout: 5 * 60 * time.Second,
}

// Download fetches the content at url, validates the Content-Type is a YAML variant,
// and returns the raw bytes.
func Download(ctx context.Context, url string) ([]byte, error) {
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
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("network error for %s: %s", url, resp.Status)
	}
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("network error: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	mediaType, _, parseErr := mime.ParseMediaType(contentType)
	if parseErr != nil {
		mediaType = contentType
	}

	if !allowedContentTypes.Contains(mediaType) {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}
