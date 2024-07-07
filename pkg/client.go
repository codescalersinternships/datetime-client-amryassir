package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
)

// Client - represents a client for making HTTP requests
type Client struct {
	httpClient *http.Client
	config     Config
}

// DateTimeResponse - defines the structure of a response containing date and time
type DateTimeResponse struct {
	DateTime string `json:"datetime"`
}

// GetDateTime - performs an HTTP GET request to fetch the current date and time
func (c *Client) GetDateTime() (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.config.URL+"/datetime", nil)
	if err != nil {
		return "request failed", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "response failed", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get datetime: status code %d", resp.StatusCode)
	}

	responseType := resp.Header.Get("Content-Type")
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "can't read response body", err
	}

	var dateTime string
	switch {
	case responseType == "application/json":
		var jsonResponse DateTimeResponse
		err = json.Unmarshal(data, &jsonResponse)
		if err != nil {
			return "failed to unmarshal the json response", err
		}
		dateTime = jsonResponse.DateTime
	case strings.HasPrefix(responseType, "text/plain"):
		dateTime = string(data)
	default:
		return "", fmt.Errorf("unsupported content type: %s", responseType)
	}

	return dateTime, nil
}

// NewClient - initializes and returns a new instance of the Client struct
func NewClient(config Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: config,
	}
}

// Config - represents a configuration object
type Config struct {
	URL string
}

// LoadConfig - initializes and returns a Config struct with configuration values
func LoadConfig() Config {
	return Config{
		URL: getEnv("URL", "http://localhost:8001"),
	}
}

// getEnv - retrieves the value of an environment variable
func getEnv(key, defaultvalue string) string {
	value, err := os.LookupEnv(key)
	if !err {
		return defaultvalue
	}
	return value
}

// Retry - retry an operation that returns an error using an exponential backoff strategy
func Retry(o func() error) error {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 500 * time.Millisecond
	expBackoff.MaxInterval = 5 * time.Second
	expBackoff.MaxElapsedTime = 30 * time.Second

	return backoff.Retry(o, expBackoff)
}
