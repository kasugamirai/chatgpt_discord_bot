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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", fmt.Errorf("creating new request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status OK, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body failed: %v", err)
	}

	var responseBody responseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", fmt.Errorf("unmarshalling response body failed: %v", err)
	}

	if len(responseBody.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	return responseBody.Candidates[0].Output, nil
}
