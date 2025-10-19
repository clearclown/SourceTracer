package search

import (
	"testing"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// 🔴 Red: Semantic Scholar client test
func TestSemanticScholarClient_Creation(t *testing.T) {
	client := NewSemanticScholarClient("")

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.Name() != "semantic_scholar" {
		t.Errorf("Expected name 'semantic_scholar', got %q", client.Name())
	}
}

func TestSemanticScholarClient_BuildSearchURL(t *testing.T) {
	client := NewSemanticScholarClient("")

	query := "climate change acceleration"
	limit := 10

	url := client.buildSearchURL(query, limit)

	// URLにクエリが含まれているか
	if url == "" {
		t.Error("Expected non-empty URL")
	}

	// 基本的なURL構造チェック
	if !contains(url, "api.semanticscholar.org") {
		t.Error("URL should contain Semantic Scholar API domain")
	}
}

func TestSemanticScholarClient_ParseResponse(t *testing.T) {
	client := NewSemanticScholarClient("")

	// Mock response from Semantic Scholar API
	mockResponse := `{
		"data": [
			{
				"paperId": "abc123",
				"title": "Climate Change Acceleration Study",
				"authors": [
					{"name": "John Doe"},
					{"name": "Jane Smith"}
				],
				"url": "https://www.semanticscholar.org/paper/abc123",
				"abstract": "This paper investigates the acceleration of climate change...",
				"citationCount": 47,
				"year": 2024,
				"publicationDate": "2024-03-15"
			}
		]
	}`

	evidences, err := client.parseResponse([]byte(mockResponse))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(evidences) != 1 {
		t.Fatalf("Expected 1 evidence, got %d", len(evidences))
	}

	ev := evidences[0]
	if ev.Title != "Climate Change Acceleration Study" {
		t.Errorf("Expected title 'Climate Change Acceleration Study', got %q", ev.Title)
	}
	if len(ev.Authors) != 2 {
		t.Errorf("Expected 2 authors, got %d", len(ev.Authors))
	}
	if ev.CitationCount != 47 {
		t.Errorf("Expected citation count 47, got %d", ev.CitationCount)
	}
	if ev.SourceType != domain.SourceTypePeerReviewedJournal {
		t.Errorf("Expected source type peer_reviewed_journal, got %v", ev.SourceType)
	}
}

func TestSemanticScholarClient_Search(t *testing.T) {
	// このテストは統合テストで実施
	// ユニットテストでは外部APIを呼ばない
	t.Skip("Integration test - requires real API call")
}

func TestSemanticScholarClient_CalculateCredibility(t *testing.T) {
	client := NewSemanticScholarClient("")

	tests := []struct {
		name          string
		citationCount int
		year          int
		minScore      float64
		maxScore      float64
	}{
		{
			name:          "high citations recent",
			citationCount: 100,
			year:          2024,
			minScore:      0.85,
			maxScore:      1.0,
		},
		{
			name:          "low citations old",
			citationCount: 5,
			year:          2010,
			minScore:      0.4,
			maxScore:      0.8,
		},
		{
			name:          "no citations recent",
			citationCount: 0,
			year:          2024,
			minScore:      0.5,
			maxScore:      0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := client.calculateCredibility(tt.citationCount, tt.year)

			if score < tt.minScore {
				t.Errorf("Score %f is below minimum %f", score, tt.minScore)
			}
			if score > tt.maxScore {
				t.Errorf("Score %f is above maximum %f", score, tt.maxScore)
			}
		})
	}
}
