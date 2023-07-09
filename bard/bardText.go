package bard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type text struct {
	Text string `json:"text"`
}

type safetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type candidate struct {
	Output        string         `json:"output"`
	SafetyRatings []safetyRating `json:"safetyRatings"`
}

type responseBody struct {
	Candidates []candidate `json:"candidates"`
}

type requestBody struct {
	Temperature    float32 `json:"temperature"`
	CandidateCount int     `json:"candidate_count"`
	TopK           int     `json:"top_k"`
	TopP           float32 `json:"top_p"`
	Prompt         text    `json:"prompt"`
}

func GenerateTextResponse(input string) (string, error) {
	apiKey := os.Getenv("BARD_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta2/models/text-bison-001:generateText?key=%s", apiKey)

	requestBody := &requestBody{
		Temperature:    0.25,
		CandidateCount: 1,
		TopK:           40,
		TopP:           0.95,
		Prompt:         text{Text: input},
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
