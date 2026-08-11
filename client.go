package ebecasv2client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is a client for interacting with the eBECAS V2 API.
type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
	pageSize    int
}

// Config contains the configuration required to create a eBECAS V2 API client.
type Config struct {
	// BaseURL is the eBECAS V2 API URL and must include the /api/v2 path.
	// Example: https://college.ap2.ebecas.app/api/v2
	BaseURL string

	// AccessToken is the eBECAS V2 API access token used for authentication.
	AccessToken string

	// HTTPClient is an optional custom HTTP client.
	// If nil, a client with a 15-second timeout is used.
	HTTPClient *http.Client

	// PageSize is the number of records requested per eBECAS V2 API page.
	// If zero, the default page size of 10 is used.
	// The value must be between 1 and 100.
	PageSize int
}

// NewClient creates a new eBECAS V2 API client using the provided configuration.
func NewClient(config Config) (*Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, errors.New("base URL is required")
	}

	if strings.TrimSpace(config.AccessToken) == "" {
		return nil, errors.New("access token is required")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	pageSize := config.PageSize

	if pageSize < 0 {
		return nil, errors.New("page size must not be negative")
	}

	if pageSize == 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		return nil, errors.New("page size must not exceed 100")
	}

	return &Client{
		baseURL:     strings.TrimRight(config.BaseURL, "/"),
		accessToken: config.AccessToken,
		httpClient:  httpClient,
		pageSize:    pageSize,
	}, nil
}

// do executes an authenticated eBECAS V2 API request.
//
// It returns the response body, HTTP status code, and any request or eBECAS V2 API error.
func (c *Client) do(req *http.Request) ([]byte, int, error) {
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute eBECAS V2 API request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("read eBECAS V2 API response body: %w", err)
	}

	return data, res.StatusCode, nil
}
