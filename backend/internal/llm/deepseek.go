package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// DeepSeekProvider implements LLM provider for DeepSeek
type DeepSeekProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(apiKey string) *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1/chat/completions",
		model:   "deepseek-chat",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// DeepSeekRequest represents the request to DeepSeek API
type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

// DeepSeekMessage represents a message in the request
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekResponse represents the response from DeepSeek API
type DeepSeekResponse struct {
	Choices []DeepSeekChoice `json:"choices"`
}

// DeepSeekChoice represents a choice in the response
type DeepSeekChoice struct {
	Message DeepSeekMessage `json:"message"`
}

// Classify classifies text using DeepSeek
func (d *DeepSeekProvider) Classify(ctx context.Context, text string) (*domain.Claim, error) {
	prompt := d.buildClassificationPrompt(text)

	reqBody := DeepSeekRequest{
		Model: d.model,
		Messages: []DeepSeekMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.apiKey))

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call DeepSeek API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("DeepSeek API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return d.parseResponse(body)
}

// buildClassificationPrompt builds the prompt for classification
func (d *DeepSeekProvider) buildClassificationPrompt(text string) string {
	return fmt.Sprintf(`You are an expert fact-checker. Classify the following text as either "opinion" or "fact".

Rules:
- "opinion": Contains subjective judgments, preferences, or beliefs (e.g., "best", "worst", "should")
- "fact": Contains objective, verifiable statements

Text to classify: "%s"

Respond ONLY with valid JSON in this exact format:
{
  "type": "opinion" or "fact",
  "confidence": 0.0 to 1.0
}`, text)
}

// parseResponse parses the DeepSeek API response
func (d *DeepSeekProvider) parseResponse(body []byte) (*domain.Claim, error) {
	var response DeepSeekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := response.Choices[0].Message.Content

	// Parse the JSON response from the model
	var result struct {
		Type       string  `json:"type"`
		Confidence float64 `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse classification result: %w", err)
	}

	// Convert string type to domain.ClaimType
	var claimType domain.ClaimType
	switch result.Type {
	case "opinion":
		claimType = domain.Opinion
	case "fact":
		claimType = domain.Fact
	default:
		claimType = domain.Unclear
	}

	claim := &domain.Claim{
		Type:       claimType,
		Confidence: result.Confidence,
	}

	return claim, nil
}
