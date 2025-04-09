package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenRouterClient представляет клиент для работы с OpenRouter API
type OpenRouterClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewOpenRouterClient создает новый клиент для OpenRouter API
func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		APIKey:  apiKey,
		BaseURL: "https://openrouter.ai/api/v1",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateDescriptionRequest для OpenRouter API
type GenerateDescriptionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateDescriptionResponse от OpenRouter API
type GenerateDescriptionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateCompanyDescription с использованием OpenRouter
func (c *OpenRouterClient) GenerateCompanyDescription(name, businessSphere string) (string, error) {
	prompt := fmt.Sprintf("Напиши краткое, но информативное описание для компании '%s', которая работает в сфере: %s. "+
		"Описание должно быть на русском языке, длиной 2-3 предложения, привлекательным для клиентов и отражать суть бизнеса.",
		name, businessSphere)

	requestBody := GenerateDescriptionRequest{
		Model: "mistralai/mistral-7b-instruct",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   150,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("HTTP-Referer", "https://morozovdesign.art/") // Укажите ваш домен
	req.Header.Set("X-Title", "Company Description Generator")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var response GenerateDescriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if response.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
