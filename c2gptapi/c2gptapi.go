package c2gptapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const OpenAIAPIURL = "https://api.openai.com/v1/chat/completions"

type Choice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Stream      bool      `json:"stream"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

func ChatWithGPT(prompt string, output chan string) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Printf("error: OPENAI_API_KEY environment variable not set")
	}
	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := &ChatCompletionRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Stream:      true,
		Temperature: 1,    // Adjust the temperature
		MaxTokens:   1000, // Adjust the max_tokens
	}
	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		fmt.Print(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", OpenAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Print(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var response Response

	for scanner.Scan() {
		s := scanner.Bytes()
		if len(s) > 6 {
			err = json.Unmarshal(s[6:], &response)
			if err != nil {
				fmt.Printf("error unmarshalling JSON data: %v", err)
			}

			if response.Choices[0].FinishReason == "stop" {
				break
			}
			output <- response.Choices[0].Delta.Content
		}
	}
	close(output)
}
