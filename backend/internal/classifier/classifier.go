package classifier

import (
	"strings"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// OpinionKeywords are words that typically indicate an opinion
var OpinionKeywords = []string{
	"best", "worst", "should", "must", "better", "worse",
	"excellent", "terrible", "amazing", "awful", "prefer",
}

// Classifier classifies claims as Opinion, Fact, Mixed, or Unclear
type Classifier struct {
	// 将来的にLLM clientを追加
}

// NewClassifier creates a new Classifier instance
func NewClassifier() *Classifier {
	return &Classifier{}
}

// ClassificationResult holds the result of classification
type ClassificationResult struct {
	Type       domain.ClaimType
	Confidence float64
	Reasoning  string
}

// Classify classifies a text as Opinion, Fact, Mixed, or Unclear
// これは仮実装 - 将来的にLLMベースに置き換える
func (c *Classifier) Classify(text string) *ClassificationResult {
	// 空文字列の処理
	if strings.TrimSpace(text) == "" {
		return &ClassificationResult{
			Type:       domain.Unclear,
			Confidence: 0.0,
			Reasoning:  "Empty text",
		}
	}

	// 正規化
	normalized := strings.ToLower(text)

	// Opinion キーワードチェック
	for _, keyword := range OpinionKeywords {
		if strings.Contains(normalized, keyword) {
			return &ClassificationResult{
				Type:       domain.Opinion,
				Confidence: 0.9,
				Reasoning:  "Contains opinion keyword: " + keyword,
			}
		}
	}

	// デフォルトではFactと判定（仮実装）
	return &ClassificationResult{
		Type:       domain.Fact,
		Confidence: 0.7,
		Reasoning:  "No opinion keywords found",
	}
}
