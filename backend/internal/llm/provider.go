package llm

import (
	"context"
	"fmt"
)

// Provider is the interface that all LLM providers must implement
type Provider interface {
	Call(ctx context.Context, prompt string) (string, error)
	Name() string
}

// MockProvider is a mock implementation for testing
// これはテスト用のモック - 外部依存のモックとして使用OK
type MockProvider struct {
	response string
	err      error
}

// NewMockProvider creates a new mock provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// SetResponse sets the mock response
func (m *MockProvider) SetResponse(response string) {
	m.response = response
}

// SetError sets the mock error
func (m *MockProvider) SetError(errMsg string) {
	m.err = fmt.Errorf("%s", errMsg)
}

// Call simulates an LLM API call
func (m *MockProvider) Call(ctx context.Context, prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

// Name returns the provider name
func (m *MockProvider) Name() string {
	return "mock"
}
