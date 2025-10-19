package llm

import (
	"testing"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// 🔴 Red: Test DeepSeek provider creation
func TestDeepSeekProvider_Creation(t *testing.T) {
	provider := NewDeepSeekProvider("test_key")
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}

	if provider.apiKey != "test_key" {
		t.Errorf("Expected API key 'test_key', got %s", provider.apiKey)
	}
}

// 🔴 Red: Test DeepSeek prompt building
func TestDeepSeekProvider_BuildPrompt(t *testing.T) {
	provider := NewDeepSeekProvider("test_key")
	prompt := provider.buildClassificationPrompt("Python is the best")

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if !contains(prompt, "Python is the best") {
		t.Error("Prompt should contain the input text")
	}

	if !contains(prompt, "opinion") && !contains(prompt, "fact") {
		t.Error("Prompt should mention opinion and fact")
	}
}

// 🔴 Red: Test DeepSeek response parsing
func TestDeepSeekProvider_ParseResponse(t *testing.T) {
	provider := NewDeepSeekProvider("test_key")

	tests := []struct {
		name         string
		response     string
		expectedType domain.ClaimType
		expectedConf float64
		shouldError  bool
	}{
		{
			name: "valid opinion response",
			response: `{
				"choices": [{
					"message": {
						"content": "{\"type\": \"opinion\", \"confidence\": 0.95}"
					}
				}]
			}`,
			expectedType: domain.Opinion,
			expectedConf: 0.95,
			shouldError:  false,
		},
		{
			name: "valid fact response",
			response: `{
				"choices": [{
					"message": {
						"content": "{\"type\": \"fact\", \"confidence\": 0.88}"
					}
				}]
			}`,
			expectedType: domain.Fact,
			expectedConf: 0.88,
			shouldError:  false,
		},
		{
			name:        "invalid json",
			response:    `invalid json`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claim, err := provider.parseResponse([]byte(tt.response))

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if claim.Type != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, claim.Type)
			}

			if claim.Confidence != tt.expectedConf {
				t.Errorf("Expected confidence %f, got %f", tt.expectedConf, claim.Confidence)
			}
		})
	}
}
