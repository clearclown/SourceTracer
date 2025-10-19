package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// SemanticScholarClient implements the SearchClient interface for Semantic Scholar
type SemanticScholarClient struct {
	apiKey string
	client *http.Client
}

// NewSemanticScholarClient creates a new Semantic Scholar client
func NewSemanticScholarClient(apiKey string) *SemanticScholarClient {
	return &SemanticScholarClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the client name
func (s *SemanticScholarClient) Name() string {
	return "semantic_scholar"
}

// buildSearchURL builds the search URL for Semantic Scholar API
func (s *SemanticScholarClient) buildSearchURL(query string, limit int) string {
	baseURL := "https://api.semanticscholar.org/graph/v1/paper/search"
	params := url.Values{}
	params.Add("query", query)
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("fields", "paperId,title,authors,url,abstract,citationCount,year,publicationDate,venue")

	return baseURL + "?" + params.Encode()
}

// SemanticScholarResponse represents the API response
type SemanticScholarResponse struct {
	Data []struct {
		PaperID string `json:"paperId"`
		Title   string `json:"title"`
		Authors []struct {
			Name string `json:"name"`
		} `json:"authors"`
		URL             string `json:"url"`
		Abstract        string `json:"abstract"`
		CitationCount   int    `json:"citationCount"`
		Year            int    `json:"year"`
		PublicationDate string `json:"publicationDate"`
		Venue           string `json:"venue"`
	} `json:"data"`
}

// Search searches for papers on Semantic Scholar
func (s *SemanticScholarClient) Search(ctx context.Context, query string, limit int) ([]*domain.Evidence, error) {
	searchURL := s.buildSearchURL(query, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// APIキーがある場合はヘッダーに追加
	if s.apiKey != "" {
		req.Header.Set("x-api-key", s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Semantic Scholar API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("Semantic Scholar API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return s.parseResponse(body)
}

// parseResponse parses the Semantic Scholar API response
func (s *SemanticScholarClient) parseResponse(data []byte) ([]*domain.Evidence, error) {
	var response SemanticScholarResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	evidences := make([]*domain.Evidence, 0, len(response.Data))

	for _, paper := range response.Data {
		// Authors
		authors := make([]string, len(paper.Authors))
		for i, author := range paper.Authors {
			authors[i] = author.Name
		}

		// Published date
		var publishedDate *time.Time
		if paper.PublicationDate != "" {
			if t, err := time.Parse("2006-01-02", paper.PublicationDate); err == nil {
				publishedDate = &t
			}
		}

		// Credibility calculation
		credibility := s.calculateCredibility(paper.CitationCount, paper.Year)

		// Snippet (abstract)
		snippet := paper.Abstract
		if len(snippet) > 300 {
			snippet = snippet[:297] + "..."
		}

		evidence := &domain.Evidence{
			ID:            fmt.Sprintf("evd_ss_%s", paper.PaperID),
			Source:        "Semantic Scholar",
			Title:         paper.Title,
			Authors:       authors,
			URL:           paper.URL,
			Snippet:       snippet,
			Credibility:   credibility,
			Relevance:     0.5, // デフォルト値 - 後でLLMで計算
			PublishedDate: publishedDate,
			SourceType:    domain.SourceTypePeerReviewedJournal,
			CitationCount: paper.CitationCount,
			Metadata: map[string]string{
				"venue":   paper.Venue,
				"year":    fmt.Sprintf("%d", paper.Year),
				"paperId": paper.PaperID,
			},
			CreatedAt: time.Now(),
		}

		evidences = append(evidences, evidence)
	}

	return evidences, nil
}

// calculateCredibility calculates credibility score based on citations and recency
func (s *SemanticScholarClient) calculateCredibility(citationCount, year int) float64 {
	// Base score from source type (peer-reviewed)
	baseScore := 0.7

	// Citation score (0.0 - 0.2)
	citationScore := 0.0
	if citationCount > 0 {
		// Logarithmic scale: 1-10 citations = 0.05, 10-100 = 0.15, 100+ = 0.2
		if citationCount >= 100 {
			citationScore = 0.2
		} else if citationCount >= 10 {
			citationScore = 0.15
		} else {
			citationScore = 0.05
		}
	}

	// Recency score (0.0 - 0.1)
	currentYear := time.Now().Year()
	age := currentYear - year
	recencyScore := 0.0
	if age <= 2 {
		recencyScore = 0.1
	} else if age <= 5 {
		recencyScore = 0.05
	}

	totalScore := baseScore + citationScore + recencyScore

	// Cap at 1.0
	if totalScore > 1.0 {
		totalScore = 1.0
	}

	return totalScore
}
