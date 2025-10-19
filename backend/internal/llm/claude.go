package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClaudeProvider implements the Provider interface for Anthropic Claude
type ClaudeProvider struct {
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey, model string, maxTokens int) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the provider name
func (c *ClaudeProvider) Name() string {
	return "claude"
}

// Call makes an API call to Claude
func (c *ClaudeProvider) Call(ctx context.Context, prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model":      c.model,
		"max_tokens": c.maxTokens,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return "", fmt.Errorf("Claude API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return response.Content[0].Text, nil
}

// buildClassificationPrompt builds a prompt for claim classification
func (c *ClaudeProvider) buildClassificationPrompt(text string) string {
	return fmt.Sprintf(`You are an expert fact-checker and claim analyzer. Your task is to classify the following text as either "opinion", "fact", "mixed", or "unclear".

Instructions:
- "opinion": Subjective statements, personal beliefs, value judgments
- "fact": Objective, verifiable claims
- "mixed": Contains both factual and opinion elements
- "unclear": Ambiguous or insufficient information

Text to classify: "%s"

Respond ONLY with valid JSON in this exact format:
{
  "type": "opinion|fact|mixed|unclear",
  "confidence": 0.0-1.0,
  "reasoning": "brief explanation"
}`, text)
}

// ClassificationResult represents the parsed classification result
type ClassificationResult struct {
	Type       string  `json:"type"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning"`
}

// parseClassificationResponse parses the JSON response from Claude
func (c *ClaudeProvider) parseClassificationResponse(response string) (*ClassificationResult, error) {
	var result ClassificationResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse classification response: %w", err)
	}

	// バリデーション
	validTypes := map[string]bool{"opinion": true, "fact": true, "mixed": true, "unclear": true}
	if !validTypes[result.Type] {
		return nil, fmt.Errorf("invalid type: %s", result.Type)
	}

	if result.Confidence < 0.0 || result.Confidence > 1.0 {
		return nil, fmt.Errorf("confidence out of range: %f", result.Confidence)
	}

	return &result, nil
}

// Classify uses Claude to classify a claim
func (c *ClaudeProvider) Classify(ctx context.Context, text string) (*ClassificationResult, error) {
	prompt := c.buildClassificationPrompt(text)
	response, err := c.Call(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return c.parseClassificationResponse(response)
}
