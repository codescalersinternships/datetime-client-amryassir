package pkg

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetDateTime(t *testing.T) {
	t.Run("json response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			jsonResponse := DateTimeResponse{DateTime: "Sun, 07 Jul 2024 16:36:50 UTC"}
			err := json.NewEncoder(w).Encode(jsonResponse)
			if err != nil {
				t.Fatalf("Failed to encode JSON response: %v", err)
			}
		}))
		defer ts.Close()

		config := Config{URL: ts.URL}
		client := NewClient(config)

		dateTime, err := client.GetDateTime()

		assert.NoError(t, err)
		assert.Equal(t, "Sun, 07 Jul 2024 16:36:50 UTC", dateTime)
	})
	t.Run("plain text response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			_, err := w.Write([]byte("Sun, 07 Jul 2024 16:36:50 UTC"))
			if err != nil {
				t.Fatalf("Failed to write plain text response: %v", err)
			}
		}))
		defer ts.Close()

		os.Setenv("URL", ts.URL)
		defer os.Unsetenv("URL")

		config := LoadConfig()
		client := NewClient(config)

		dateTime, err := client.GetDateTime()

		assert.NoError(t, err)
		assert.Equal(t, "Sun, 07 Jul 2024 16:36:50 UTC", dateTime)
	})
	t.Run("unsupported format", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			_, err := w.Write([]byte("<datetime>Sun, 07 Jul 2024 16:36:50 UTC</datetime>"))
			if err != nil {
				t.Fatalf("Failed to write XML response: %v", err)
			}
		}))
		defer ts.Close()

		config := Config{URL: ts.URL}
		client := NewClient(config)

		dateTime, err := client.GetDateTime()

		assert.Error(t, err)
		assert.Equal(t, "", dateTime)
	})
}
