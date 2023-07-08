// Package chatGPT provides a simple interface to interact with GPT
package chatGPT

// Import required packages
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Define constants
const OpenAIAPIURL = "https://api.openai.com/v1/chat/completions"

// Define struct types for API request and response
type Choice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

// Message is a structure representing a message sent to the API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the structure of the request sent to the OpenAI API.
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

// ChatWithGPT sends a chat completion request to the OpenAI API
func ChatWithGPT(prompt string, output chan string) {
	// Get the OpenAI API key from the environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Printf("error: OPENAI_API_KEY environment variable not set")
	}
	// Initialize message structure for the API request
	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// Create the chat completion request object
	requestBody := &ChatCompletionRequest{
		Model:    "gpt-4",
		Messages: messages,
		Stream:   true,
	}
	// Convert the request object to JSON
	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		fmt.Print(err)
	}

	// Create a new HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", OpenAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Print(err)
	}

	// Set headers for the HTTP request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the HTTP request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()

	// Process the API response
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

// Package main provides the main application to run the Discord
