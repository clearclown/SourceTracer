package llm

import (
	"testing"
)

// 🔴 Red: Claude Provider test
func TestClaudeProvider_Creation(t *testing.T) {
	apiKey := "test-api-key"
	model := "claude-sonnet-4-20250514"

	provider := NewClaudeProvider(apiKey, model, 4000)

	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
	if provider.Name() != "claude" {
		t.Errorf("Expected name 'claude', got %q", provider.Name())
	}
}

func TestClaudeProvider_BuildPrompt(t *testing.T) {
	provider := NewClaudeProvider("test-key", "test-model", 4000)

	text := "Python is the best language"
	prompt := provider.buildClassificationPrompt(text)

	// プロンプトに必要な要素が含まれているか確認
	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	// プロンプトに入力テキストが含まれているか
	if !contains(prompt, text) {
		t.Error("Prompt should contain input text")
	}

	// JSON出力を要求しているか
	if !contains(prompt, "JSON") && !contains(prompt, "json") {
		t.Error("Prompt should request JSON output")
	}
}

func TestClaudeProvider_ParseResponse(t *testing.T) {
	provider := NewClaudeProvider("test-key", "test-model", 4000)

	tests := []struct {
		name       string
		response   string
		wantType   string
		wantConf   float64
		wantErr    bool
	}{
		{
			name:     "valid opinion response",
			response: `{"type": "opinion", "confidence": 0.95, "reasoning": "Contains subjective judgment"}`,
			wantType: "opinion",
			wantConf: 0.95,
			wantErr:  false,
		},
		{
			name:     "valid fact response",
			response: `{"type": "fact", "confidence": 0.88, "reasoning": "Verifiable claim"}`,
			wantType: "fact",
			wantConf: 0.88,
			wantErr:  false,
		},
		{
			name:     "invalid json",
			response: `not a json`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := provider.parseClassificationResponse(tt.response)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result.Type != tt.wantType {
				t.Errorf("Expected type %q, got %q", tt.wantType, result.Type)
			}
			if result.Confidence != tt.wantConf {
				t.Errorf("Expected confidence %f, got %f", tt.wantConf, result.Confidence)
			}
		})
	}
}

// 実際のAPI呼び出しはモックを使ってテスト
func TestClaudeProvider_CallWithMock(t *testing.T) {
	// このテストは統合テストで実施
	// ユニットテストでは外部APIを呼ばない
	t.Skip("Integration test - requires real API key")
}

// ヘルパー関数
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsInMiddle(s, substr)))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
