package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	AppEnv             string
	AppPort            int
	AppDebug           bool
	DefaultLLMProvider string

	OpenAI    OpenAIConfig
	Anthropic AnthropicConfig
	DeepSeek  DeepSeekConfig

	Database DatabaseConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
}

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
}

// AnthropicConfig holds Anthropic Claude-specific configuration
type AnthropicConfig struct {
	APIKey    string
	Model     string
	MaxTokens int
}

// DeepSeekConfig holds DeepSeek-specific configuration
type DeepSeekConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnvAsInt("APP_PORT", 8080),
		AppDebug:           getEnvAsBool("APP_DEBUG", true),
		DefaultLLMProvider: getEnv("DEFAULT_LLM_PROVIDER", "claude"),

		OpenAI: OpenAIConfig{
			APIKey:      getEnv("OPENAI_API_KEY", ""),
			Model:       getEnv("OPENAI_MODEL", "gpt-4-turbo-preview"),
			MaxTokens:   getEnvAsInt("OPENAI_MAX_TOKENS", 4000),
			Temperature: getEnvAsFloat("OPENAI_TEMPERATURE", 0.3),
		},

		Anthropic: AnthropicConfig{
			APIKey:    getEnv("ANTHROPIC_API_KEY", ""),
			Model:     getEnv("ANTHROPIC_MODEL", "claude-sonnet-4-20250514"),
			MaxTokens: getEnvAsInt("ANTHROPIC_MAX_TOKENS", 4000),
		},

		DeepSeek: DeepSeekConfig{
			APIKey:  getEnv("DEEPSEEK_API_KEY", ""),
			Model:   getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
			BaseURL: getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),
		},

		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvAsInt("DATABASE_PORT", 5432),
			Name:     getEnv("DATABASE_NAME", "sourcetracer"),
			User:     getEnv("DATABASE_USER", "sourcetracer_user"),
			Password: getEnv("DATABASE_PASSWORD", ""),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},

		MongoDB: MongoDBConfig{
			Host:     getEnv("MONGODB_HOST", "localhost"),
			Port:     getEnvAsInt("MONGODB_PORT", 27017),
			Database: getEnv("MONGODB_DATABASE", "sourcetracer_raw"),
			User:     getEnv("MONGODB_USER", ""),
			Password: getEnv("MONGODB_PASSWORD", ""),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	return cfg, nil
}

// getEnv gets environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets environment variable as integer with default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool gets environment variable as boolean with default
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsFloat gets environment variable as float64 with default
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// 本番環境では必須チェック
	if c.AppEnv == "production" {
		if c.DefaultLLMProvider == "openai" && c.OpenAI.APIKey == "" {
			return fmt.Errorf("OPENAI_API_KEY is required in production")
		}
		if c.DefaultLLMProvider == "claude" && c.Anthropic.APIKey == "" {
			return fmt.Errorf("ANTHROPIC_API_KEY is required in production")
		}
	}
	return nil
}
