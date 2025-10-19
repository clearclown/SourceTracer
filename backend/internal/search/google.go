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

// GoogleClient implements search using Google Custom Search API
type GoogleClient struct {
	apiKey  string
	cx      string // Custom Search Engine ID
	baseURL string
	client  *http.Client
}

// NewGoogleClient creates a new Google Custom Search client
func NewGoogleClient(apiKey, cx string) *GoogleClient {
	return &GoogleClient{
		apiKey:  apiKey,
		cx:      cx,
		baseURL: "https://www.googleapis.com/customsearch/v1",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GoogleSearchResponse represents the Google Custom Search API response
type GoogleSearchResponse struct {
	Items []GoogleSearchItem `json:"items"`
}

// GoogleSearchItem represents a single search result
type GoogleSearchItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Snippet     string `json:"snippet"`
	DisplayLink string `json:"displayLink"`
}

// Search searches using Google Custom Search API
func (g *GoogleClient) Search(ctx context.Context, query string, maxResults int) ([]*domain.Evidence, error) {
	searchURL := g.buildSearchURL(query, maxResults)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("Google API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return g.parseResponse(body)
}

// buildSearchURL constructs Google Custom Search API URL
func (g *GoogleClient) buildSearchURL(query string, maxResults int) string {
	params := url.Values{}
	params.Set("key", g.apiKey)
	params.Set("cx", g.cx)
	params.Set("q", query)
	params.Set("num", fmt.Sprintf("%d", maxResults))

	return fmt.Sprintf("%s?%s", g.baseURL, params.Encode())
}

// parseResponse parses Google Custom Search JSON response
func (g *GoogleClient) parseResponse(body []byte) ([]*domain.Evidence, error) {
	var response GoogleSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	evidences := make([]*domain.Evidence, 0, len(response.Items))

	for _, item := range response.Items {
		// Determine source type based on domain
		sourceType := g.determineSourceType(item.DisplayLink)

		// Calculate credibility based on source type
		credibility := g.calculateCredibility(sourceType)

		evidence := &domain.Evidence{
			Source:        "Google",
			Title:         item.Title,
			Authors:       []string{}, // Google doesn't provide authors
			URL:           item.Link,
			Snippet:       item.Snippet,
			Credibility:   credibility,
			Relevance:     0.7, // Default relevance
			PublishedDate: nil, // Google doesn't provide publish date in basic API
			SourceType:    sourceType,
			CitationCount: 0,
			Metadata: map[string]string{
				"domain": item.DisplayLink,
			},
		}

		evidences = append(evidences, evidence)
	}

	return evidences, nil
}

// determineSourceType determines source type from domain
func (g *GoogleClient) determineSourceType(domainStr string) domain.SourceType {
	// Simple heuristics based on domain
	if contains(domainStr, ".gov") {
		return domain.SourceTypeGovernment
	}
	if contains(domainStr, ".edu") || contains(domainStr, "scholar") {
		return domain.SourceTypePeerReviewedJournal
	}
	if contains(domainStr, "news") || contains(domainStr, "times") || contains(domainStr, "post") {
		return domain.SourceTypeNews
	}
	if contains(domainStr, "blog") || contains(domainStr, "medium") {
		return domain.SourceTypeBlog
	}
	return domain.SourceTypeUnknown
}

// calculateCredibility calculates credibility based on source type
func (g *GoogleClient) calculateCredibility(sourceType domain.SourceType) float64 {
	switch sourceType {
	case domain.SourceTypePeerReviewedJournal:
		return 0.9
	case domain.SourceTypeGovernment:
		return 0.85
	case domain.SourceTypeNews:
		return 0.7
	case domain.SourceTypePreprint:
		return 0.75
	case domain.SourceTypeBlog:
		return 0.5
	default:
		return 0.6
	}
}
