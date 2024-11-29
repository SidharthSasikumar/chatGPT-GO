package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiKeyEnvVar = "OPENAI_API_KEY"
	apiBaseURL   = "https://api.openai.com/v1"
	chatURL      = apiBaseURL + "/engines/davinci/completions"
)

type CompletionRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Created int64  `json:"created"`
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	apiKey := os.Getenv(apiKeyEnvVar)
	if apiKey == "" {
		log.Fatal("API key not found in environment variable:", apiKeyEnvVar)
	}

	prompt := "Once upon a time"
	maxTokens := 100
	temperature := 0.7

	response, err := getCompletion(apiKey, prompt, maxTokens, temperature)
	if err != nil {
		log.Fatalf("Error getting completion: %v", err)
	}

	if len(response.Choices) > 0 {
		log.Println("Response:", response.Choices[0].Text)
	} else {
		log.Println("No completion choices found.")
	}
}

func getCompletion(apiKey, prompt string, maxTokens int, temperature float64) (CompletionResponse, error) {
	requestData := CompletionRequest{
		Prompt:      prompt,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return CompletionResponse{}, errors.New("failed to marshal request data: " + err.Error())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", chatURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return CompletionResponse{}, errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return CompletionResponse{}, errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return CompletionResponse{}, errors.New("API error: " + string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CompletionResponse{}, errors.New("failed to read response body: " + err.Error())
	}

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return CompletionResponse{}, errors.New("failed to unmarshal response: " + err.Error())
	}

	return completionResponse, nil
}
