package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	apiKey     = "API_KEY"
	apiBaseURL = "https://api.openai.com/v1"
	chatURL    = apiBaseURL + "/engines/davinci/completions"
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
	prompt := "Once upon a time"
	maxTokens := 1000
	temperature := 0.7

	completion := getCompletion(prompt, maxTokens, temperature)
	if len(completion.Choices) > 0 {
		response := completion.Choices[0].Text
		log.Println(response)
	} else {
		log.Println("No completion choices found.")
	}
}

func getCompletion(prompt string, maxTokens int, temperature float64) CompletionResponse {
	requestData := CompletionRequest{
		Prompt:      prompt,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Fatal("Failed to marshal request data:", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", chatURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		log.Fatal("Failed to unmarshal response:", err)
	}

	return completionResponse
}
