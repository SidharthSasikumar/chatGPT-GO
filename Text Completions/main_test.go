package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetCompletion_Success(t *testing.T) {
	mockResponse := `{
		"id": "test-id",
		"object": "text_completion",
		"model": "davinci",
		"created": 1234567890,
		"choices": [{"text": "This is a test response."}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	chatURL = server.URL // Override the chatURL for testing

	response, err := getCompletion("test-key", "Test prompt", 10, 0.7)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Text != "This is a test response." {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestGetCompletion_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	chatURL = server.URL // Override the chatURL for testing

	_, err := getCompletion("wrong-key", "Test prompt", 10, 0.7)
	if err == nil || err.Error() != "API error: Unauthorized\n" {
		t.Errorf("Expected unauthorized error, got %v", err)
	}
}

func TestMain_MissingAPIKey(t *testing.T) {
	originalEnv := os.Getenv(apiKeyEnvVar)
	defer os.Setenv(apiKeyEnvVar, originalEnv)

	os.Unsetenv(apiKeyEnvVar)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing API key")
		}
	}()
	main()
}
