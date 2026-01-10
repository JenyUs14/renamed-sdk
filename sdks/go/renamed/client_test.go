package renamed

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("creates client with API key", func(t *testing.T) {
		client := NewClient("rt_test123")
		if client == nil {
			t.Error("expected client to be created")
		}
		if client.apiKey != "rt_test123" {
			t.Errorf("expected apiKey to be rt_test123, got %s", client.apiKey)
		}
	})

	t.Run("uses default base URL", func(t *testing.T) {
		client := NewClient("rt_test123")
		if client.baseURL != defaultBaseURL {
			t.Errorf("expected baseURL to be %s, got %s", defaultBaseURL, client.baseURL)
		}
	})

	t.Run("accepts custom options", func(t *testing.T) {
		client := NewClient("rt_test123",
			WithBaseURL("https://custom.api.com"),
			WithMaxRetries(5),
		)
		if client.baseURL != "https://custom.api.com" {
			t.Errorf("expected baseURL to be https://custom.api.com, got %s", client.baseURL)
		}
		if client.maxRetries != 5 {
			t.Errorf("expected maxRetries to be 5, got %d", client.maxRetries)
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("returns user data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify Authorization header
			auth := r.Header.Get("Authorization")
			if auth != "Bearer rt_test123" {
				t.Errorf("expected Authorization header, got %s", auth)
			}

			user := User{
				ID:      "user123",
				Email:   "test@example.com",
				Name:    "Test User",
				Credits: 100,
			}
			json.NewEncoder(w).Encode(user)
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))
		user, err := client.GetUser(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != "user123" {
			t.Errorf("expected ID user123, got %s", user.ID)
		}
		if user.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", user.Email)
		}
		if user.Credits != 100 {
			t.Errorf("expected credits 100, got %d", user.Credits)
		}
	})

	t.Run("returns AuthenticationError on 401", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid API key"})
		}))
		defer server.Close()

		client := NewClient("rt_invalid", WithBaseURL(server.URL))
		_, err := client.GetUser(context.Background())

		if err == nil {
			t.Error("expected error")
		}
		if _, ok := err.(*AuthenticationError); !ok {
			t.Errorf("expected AuthenticationError, got %T", err)
		}
	})

	t.Run("returns InsufficientCreditsError on 402", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]string{"error": "Insufficient credits"})
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))
		_, err := client.GetUser(context.Background())

		if err == nil {
			t.Error("expected error")
		}
		if _, ok := err.(*InsufficientCreditsError); !ok {
			t.Errorf("expected InsufficientCreditsError, got %T", err)
		}
	})

	t.Run("returns RateLimitError on 429", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"error": "Rate limit exceeded", "retryAfter": 60.0})
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))
		_, err := client.GetUser(context.Background())

		if err == nil {
			t.Error("expected error")
		}
		rateLimitErr, ok := err.(*RateLimitError)
		if !ok {
			t.Errorf("expected RateLimitError, got %T", err)
		}
		if rateLimitErr.RetryAfter != 60 {
			t.Errorf("expected RetryAfter 60, got %d", rateLimitErr.RetryAfter)
		}
	})

	t.Run("returns ValidationError on 400", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))
		_, err := client.GetUser(context.Background())

		if err == nil {
			t.Error("expected error")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})
}

func TestRename(t *testing.T) {
	t.Run("uploads file and returns result", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}

			// Verify multipart form
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				t.Errorf("failed to parse multipart form: %v", err)
			}

			file, header, err := r.FormFile("file")
			if err != nil {
				t.Errorf("failed to get file: %v", err)
			}
			defer file.Close()

			if header.Filename == "" {
				t.Error("expected filename")
			}

			result := RenameResult{
				OriginalFilename:  header.Filename,
				SuggestedFilename: "2025-01-15_Invoice.pdf",
				FolderPath:        "2025/Invoices",
				Confidence:        0.95,
			}
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))

		// Create temp file for testing
		result, err := client.RenameReader(
			context.Background(),
			&mockReader{data: []byte("fake pdf content")},
			"test.pdf",
			nil,
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.SuggestedFilename != "2025-01-15_Invoice.pdf" {
			t.Errorf("expected suggested filename, got %s", result.SuggestedFilename)
		}
		if result.Confidence != 0.95 {
			t.Errorf("expected confidence 0.95, got %f", result.Confidence)
		}
	})
}

func TestPDFSplit(t *testing.T) {
	t.Run("returns async job", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(pdfSplitResponse{
				StatusURL: "https://api.example.com/status/job123",
			})
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))

		job, err := client.PDFSplitReader(
			context.Background(),
			&mockReader{data: []byte("fake pdf content")},
			"test.pdf",
			nil,
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if job == nil {
			t.Error("expected job")
		}
	})
}

func TestAsyncJob(t *testing.T) {
	t.Run("polls until completed", func(t *testing.T) {
		callCount := 0
		mockResult := PdfSplitResult{
			OriginalFilename: "multi.pdf",
			Documents: []SplitDocument{
				{
					Index:       0,
					Filename:    "doc1.pdf",
					Pages:       "1-5",
					DownloadURL: "https://...",
					Size:        1000,
				},
			},
			TotalPages: 10,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			var status JobStatusResponse
			if callCount < 3 {
				status = JobStatusResponse{
					JobID:    "job123",
					Status:   JobStatusProcessing,
					Progress: callCount * 33,
				}
			} else {
				status = JobStatusResponse{
					JobID:    "job123",
					Status:   JobStatusCompleted,
					Progress: 100,
					Result:   &mockResult,
				}
			}
			json.NewEncoder(w).Encode(status)
		}))
		defer server.Close()

		client := NewClient("rt_test123", WithBaseURL(server.URL))

		job := &AsyncJob{
			client:       client,
			statusURL:    server.URL + "/status/job123",
			pollInterval: 1, // 1ms for fast testing
			maxAttempts:  10,
		}

		progressUpdates := []int{}
		result, err := job.Wait(context.Background(), func(s *JobStatusResponse) {
			progressUpdates = append(progressUpdates, s.Progress)
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.OriginalFilename != "multi.pdf" {
			t.Errorf("expected original filename multi.pdf, got %s", result.OriginalFilename)
		}
		if len(result.Documents) != 1 {
			t.Errorf("expected 1 document, got %d", len(result.Documents))
		}

		// Check progress was tracked
		found33 := false
		found66 := false
		for _, p := range progressUpdates {
			if p == 33 {
				found33 = true
			}
			if p == 66 {
				found66 = true
			}
		}
		if !found33 || !found66 {
			t.Errorf("expected progress updates 33 and 66, got %v", progressUpdates)
		}
	})
}

// mockReader is a simple io.Reader for testing
type mockReader struct {
	data []byte
	pos  int
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
