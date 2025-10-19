package analyzer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/sourcetracer/internal/classifier"
	"github.com/yourusername/sourcetracer/internal/domain"
)

// Analyzer analyzes text and extracts claims with evidences
type Analyzer struct {
	classifier    *classifier.Classifier
	searchClients []SearchClient
}

// SearchClient is an interface for search clients
type SearchClient interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.Evidence, error)
	Name() string
}

// AnalyzeOptions contains options for analysis
type AnalyzeOptions struct {
	IncludeEvidences bool
	MaxClaims        int
	MaxEvidences     int
	SearchEngines    []string
}

// AnalyzeResult contains the analysis result
type AnalyzeResult struct {
	AnalysisID         string           `json:"analysis_id"`
	Claims             []*domain.Claim  `json:"claims"`
	OverallCredibility float64          `json:"overall_credibility"`
	ProcessingTimeMS   int64            `json:"processing_time_ms"`
	CreatedAt          time.Time        `json:"created_at"`
}

// NewAnalyzer creates a new Analyzer
func NewAnalyzer(c *classifier.Classifier, searchClients []SearchClient) *Analyzer {
	return &Analyzer{
		classifier:    c,
		searchClients: searchClients,
	}
}

// Analyze analyzes text and extracts claims
func (a *Analyzer) Analyze(ctx context.Context, text string, opts AnalyzeOptions) (*AnalyzeResult, error) {
	startTime := time.Now()

	// Extract claims from text
	claims := a.extractClaims(text)

	// Limit claims if specified
	if opts.MaxClaims > 0 && len(claims) > opts.MaxClaims {
		claims = claims[:opts.MaxClaims]
	}

	// Classify each claim
	if err := a.classifyClaims(ctx, claims); err != nil {
		return nil, fmt.Errorf("failed to classify claims: %w", err)
	}

	// Search for evidences if requested
	if opts.IncludeEvidences && len(a.searchClients) > 0 {
		if err := a.searchEvidences(ctx, claims, opts); err != nil {
			// Log error but continue - evidence search is optional
			// In production, we'd use a proper logger
		}
	}

	// Calculate overall credibility
	overallCredibility := a.calculateOverallCredibility(claims)

	processingTime := time.Since(startTime).Milliseconds()

	return &AnalyzeResult{
		AnalysisID:         generateAnalysisID(),
		Claims:             claims,
		OverallCredibility: overallCredibility,
		ProcessingTimeMS:   processingTime,
		CreatedAt:          time.Now(),
	}, nil
}

// extractClaims extracts individual claims from text
// Simple implementation: split by sentence
func (a *Analyzer) extractClaims(text string) []*domain.Claim {
	// 簡単な実装: 句点で分割
	// 将来的にはNLPライブラリを使用
	sentences := splitSentences(text)

	claims := make([]*domain.Claim, 0, len(sentences))
	position := 0

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		claim := domain.NewClaim(sentence, position, position+len(sentence))
		claims = append(claims, claim)

		position += len(sentence) + 1 // +1 for separator
	}

	return claims
}

// splitSentences splits text into sentences
func splitSentences(text string) []string {
	// シンプルな実装: ピリオド、感嘆符、疑問符で分割
	replacer := strings.NewReplacer(".", ".|", "!", "!|", "?", "?|")
	marked := replacer.Replace(text)
	sentences := strings.Split(marked, "|")

	result := make([]string, 0, len(sentences))
	for _, s := range sentences {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// classifyClaims classifies all claims
func (a *Analyzer) classifyClaims(ctx context.Context, claims []*domain.Claim) error {
	for _, claim := range claims {
		result := a.classifier.Classify(claim.Text)
		claim.Type = result.Type
		claim.Confidence = result.Confidence
	}
	return nil
}

// searchEvidences searches for evidences for each claim
func (a *Analyzer) searchEvidences(ctx context.Context, claims []*domain.Claim, opts AnalyzeOptions) error {
	for _, claim := range claims {
		// Skip unclear claims
		if claim.Type == domain.Unclear {
			continue
		}

		// Search with each client
		for _, client := range a.searchClients {
			limit := 5
			if opts.MaxEvidences > 0 {
				limit = opts.MaxEvidences
			}

			evidences, err := client.Search(ctx, claim.Text, limit)
			if err != nil {
				// Log error but continue with other clients
				continue
			}

			for _, ev := range evidences {
				claim.Evidences = append(claim.Evidences, *ev)
			}
		}
	}

	return nil
}

// calculateOverallCredibility calculates overall credibility from claims
func (a *Analyzer) calculateOverallCredibility(claims []*domain.Claim) float64 {
	if len(claims) == 0 {
		return 0.0
	}

	totalConfidence := 0.0
	for _, claim := range claims {
		totalConfidence += claim.Confidence
	}

	return totalConfidence / float64(len(claims))
}

// generateAnalysisID generates a unique analysis ID
func generateAnalysisID() string {
	return fmt.Sprintf("anl_%d", time.Now().UnixNano())
}

// WithSearchClient adds a search client to the analyzer
func (a *Analyzer) WithSearchClient(client SearchClient) *Analyzer {
	a.searchClients = append(a.searchClients, client)
	return a
}
