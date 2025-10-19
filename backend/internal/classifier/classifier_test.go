package classifier

import (
	"testing"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// 🔴 Red: 失敗するテストを書く - Opinion判定
func TestClassifier_ClassifyOpinion(t *testing.T) {
	// Arrange
	classifier := NewClassifier()
	text := "Python is the best programming language"

	// Act
	result := classifier.Classify(text)

	// Assert
	if result.Type != domain.Opinion {
		t.Errorf("Expected Opinion, got %v", result.Type)
	}
	if result.Confidence <= 0.8 {
		t.Errorf("Expected confidence > 0.8 for clear opinion, got %f", result.Confidence)
	}
}

// 🔴 Red: 失敗するテストを書く - Fact判定
func TestClassifier_ClassifyFact(t *testing.T) {
	// Arrange
	classifier := NewClassifier()
	text := "The Earth orbits the Sun"

	// Act
	result := classifier.Classify(text)

	// Assert
	if result.Type != domain.Fact {
		t.Errorf("Expected Fact, got %v", result.Type)
	}
	if result.Confidence <= 0.6 {
		t.Errorf("Expected confidence > 0.6, got %f", result.Confidence)
	}
}

// 🔴 Red: 失敗するテストを書く - 空文字列
func TestClassifier_ClassifyEmpty(t *testing.T) {
	classifier := NewClassifier()
	text := ""

	result := classifier.Classify(text)

	if result.Type != domain.Unclear {
		t.Errorf("Expected Unclear for empty text, got %v", result.Type)
	}
	if result.Confidence != 0.0 {
		t.Errorf("Expected confidence 0.0 for empty text, got %f", result.Confidence)
	}
}

// 🔴 Red: 失敗するテストを書く - Opinion keywords detection
func TestClassifier_OpinionKeywords(t *testing.T) {
	classifier := NewClassifier()

	tests := []struct {
		name string
		text string
		want domain.ClaimType
	}{
		{"best", "This is the best approach", domain.Opinion},
		{"worst", "That's the worst idea ever", domain.Opinion},
		{"should", "You should try this method", domain.Opinion},
		{"must", "We must implement this feature", domain.Opinion},
		{"better", "Go is better than Java", domain.Opinion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.Classify(tt.text)
			if result.Type != tt.want {
				t.Errorf("For %q: expected %v, got %v", tt.text, tt.want, result.Type)
			}
		})
	}
}
