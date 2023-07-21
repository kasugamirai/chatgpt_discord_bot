package bard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Content string `json:"content"`
}

type Prompt struct {
	Messages []Message `json:"messages"`
}

type RequestBody struct {
	Model          string  `json:"model"`
	Temperature    float32 `json:"temperature"`
	CandidateCount int     `json:"candidate_count"`
	TopK           int     `json:"top_k"`
	TopP           float32 `json:"top_p"`
	Prompt         Prompt  `json:"prompt"`
}

type Candidate struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Candidates []Candidate `json:"candidates"`
	Messages   []Message   `json:"messages"`
}

func GenerateChatResponse(input string) (string, error) {
	apiKey := os.Getenv("BARD_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("BARD_API_KEY environment variable not set")
	}

	url := "https://generativelanguage.googleapis.com/v1beta2/models/chat-bison-001:generateMessage?key=" + apiKey

	messages := []Message{{Content: input}}

	requestBody := &RequestBody{
		Temperature:    0.25,
		CandidateCount: 1,
		TopK:           40,
		TopP:           0.95,
		Prompt:         Prompt{Messages: messages},
	}

	jsonValue, _ := json.Marshal(requestBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", fmt.Errorf("The HTTP request failed with error %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read the response body: %v", err)
	}

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal the response body: %v", err)
	}

	if len(responseBody.Candidates) > 0 {
		return responseBody.Candidates[0].Content, nil
	} else {
		return "", fmt.Errorf("No candidates found in the response")
	}
}
