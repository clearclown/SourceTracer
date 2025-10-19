package config

import (
	"os"
	"testing"
)

// 🔴 Red: Config loading test
func TestConfig_Load(t *testing.T) {
	// Arrange
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DEFAULT_LLM_PROVIDER", "claude")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DEFAULT_LLM_PROVIDER")
	}()

	// Act
	cfg, err := Load()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cfg.AppEnv != "test" {
		t.Errorf("Expected AppEnv 'test', got %q", cfg.AppEnv)
	}
	if cfg.AppPort != 9090 {
		t.Errorf("Expected AppPort 9090, got %d", cfg.AppPort)
	}
	if cfg.DefaultLLMProvider != "claude" {
		t.Errorf("Expected DefaultLLMProvider 'claude', got %q", cfg.DefaultLLMProvider)
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	// 環境変数をクリア
	os.Clearenv()

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Expected no error with defaults, got %v", err)
	}

	// デフォルト値のチェック
	if cfg.AppEnv != "development" {
		t.Errorf("Expected default AppEnv 'development', got %q", cfg.AppEnv)
	}
	if cfg.AppPort != 8080 {
		t.Errorf("Expected default AppPort 8080, got %d", cfg.AppPort)
	}
	if cfg.DefaultLLMProvider != "claude" {
		t.Errorf("Expected default LLM provider 'claude', got %q", cfg.DefaultLLMProvider)
	}
}

func TestConfig_LLMProviders(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "sk-test-openai")
	os.Setenv("ANTHROPIC_API_KEY", "sk-test-anthropic")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("ANTHROPIC_API_KEY")
	}()

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.OpenAI.APIKey != "sk-test-openai" {
		t.Errorf("Expected OpenAI API key, got %q", cfg.OpenAI.APIKey)
	}
	if cfg.Anthropic.APIKey != "sk-test-anthropic" {
		t.Errorf("Expected Anthropic API key, got %q", cfg.Anthropic.APIKey)
	}
}
