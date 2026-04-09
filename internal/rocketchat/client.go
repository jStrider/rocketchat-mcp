package rocketchat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is an immutable HTTP client for the Rocket.Chat REST API v1.
type Client struct {
	baseURL    string
	userID     string
	authToken  string
	httpClient *http.Client
	maxBody    int64
}

// NewClient creates a new Rocket.Chat API client.
func NewClient(baseURL, userID, authToken string, maxBody int64) *Client {
	return &Client{
		baseURL:   strings.TrimRight(baseURL, "/") + "/api/v1",
		userID:    userID,
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxBody: maxBody,
	}
}

// get performs an authenticated GET request and returns the raw JSON body.
func (c *Client) get(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	path := endpoint
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	return c.do(ctx, http.MethodGet, path, nil)
}

// post performs an authenticated POST request with a JSON body.
func (c *Client) post(ctx context.Context, endpoint string, body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body for %s: %w", endpoint, err)
	}
	return c.do(ctx, http.MethodPost, endpoint, data)
}

// do executes the HTTP request with auth headers, body limit, and error handling.
func (c *Client) do(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = strings.NewReader(string(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request %s %s: %w", method, path, err)
	}

	req.Header.Set("X-Auth-Token", c.authToken)
	req.Header.Set("X-User-Id", c.userID)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, c.maxBody)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read response from %s: %w", path, err)
	}

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp.StatusCode, path, raw)
	}

	return raw, nil
}

// parseError extracts a meaningful error from an HTTP error response.
func (c *Client) parseError(status int, path string, body []byte) error {
	var apiErr struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
		return fmt.Errorf("API error %d on %s: %s", status, path, apiErr.Error)
	}
	if apiErr.Message != "" {
		return fmt.Errorf("API error %d on %s: %s", status, path, apiErr.Message)
	}

	switch status {
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized: check ROCKETCHAT_AUTH_TOKEN and ROCKETCHAT_USER_ID")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden: insufficient permissions for %s", path)
	case http.StatusNotFound:
		return fmt.Errorf("not found: %s", path)
	default:
		return fmt.Errorf("HTTP %d on %s", status, path)
	}
}

// pagingParams builds count/offset query parameters from ListOptions.
func pagingParams(opts ListOptions, defaultCount int) url.Values {
	count := opts.Count
	if count <= 0 {
		count = defaultCount
	}
	if count > 200 {
		count = 200
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}
	return url.Values{
		"count":  {fmt.Sprintf("%d", count)},
		"offset": {fmt.Sprintf("%d", offset)},
	}
}
