package llm

import (
	"context"
	"testing"
)

// 🔴 Red: LLM Provider interface test
func TestMockProvider_Classify(t *testing.T) {
	// Arrange
	provider := NewMockProvider()
	provider.SetResponse(`{"type": "opinion", "confidence": 0.95}`)

	ctx := context.Background()
	prompt := "Classify: Python is the best language"

	// Act
	response, err := provider.Call(ctx, prompt)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response == "" {
		t.Error("Expected non-empty response")
	}
	if response != `{"type": "opinion", "confidence": 0.95}` {
		t.Errorf("Expected mock response, got %q", response)
	}
}

func TestMockProvider_Error(t *testing.T) {
	// Arrange
	provider := NewMockProvider()
	provider.SetError("API rate limit exceeded")

	ctx := context.Background()

	// Act
	_, err := provider.Call(ctx, "test")

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "API rate limit exceeded" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestProvider_Interface(t *testing.T) {
	// Arrange
	var _ Provider = (*MockProvider)(nil)

	// This test just ensures MockProvider implements Provider interface
	// コンパイルが通ればOK
}
