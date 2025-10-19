package search

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// ArxivClient implements search for arXiv papers
type ArxivClient struct {
	baseURL string
	client  *http.Client
}

// NewArxivClient creates a new arXiv search client
func NewArxivClient() *ArxivClient {
	return &ArxivClient{
		baseURL: "http://export.arxiv.org/api/query",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ArxivFeed represents arXiv API XML response
type ArxivFeed struct {
	XMLName xml.Name     `xml:"feed"`
	Entries []ArxivEntry `xml:"entry"`
}

// ArxivEntry represents a single arXiv paper
type ArxivEntry struct {
	ID        string        `xml:"id"`
	Title     string        `xml:"title"`
	Summary   string        `xml:"summary"`
	Published string        `xml:"published"`
	Authors   []ArxivAuthor `xml:"author"`
}

// ArxivAuthor represents an author
type ArxivAuthor struct {
	Name string `xml:"name"`
}

// Search searches arXiv for papers
func (a *ArxivClient) Search(ctx context.Context, query string, maxResults int) ([]*domain.Evidence, error) {
	searchURL := a.buildSearchURL(query, maxResults)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call arXiv API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck
		return nil, fmt.Errorf("arXiv API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return a.parseResponse(body)
}

// buildSearchURL constructs arXiv API search URL
func (a *ArxivClient) buildSearchURL(query string, maxResults int) string {
	params := url.Values{}
	params.Set("search_query", fmt.Sprintf("all:%s", query))
	params.Set("max_results", fmt.Sprintf("%d", maxResults))
	params.Set("sortBy", "relevance")
	params.Set("sortOrder", "descending")

	return fmt.Sprintf("%s?%s", a.baseURL, params.Encode())
}

// parseResponse parses arXiv XML response
func (a *ArxivClient) parseResponse(body []byte) ([]*domain.Evidence, error) {
	var feed ArxivFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	evidences := make([]*domain.Evidence, 0, len(feed.Entries))

	for _, entry := range feed.Entries {
		// Extract authors
		authors := make([]string, len(entry.Authors))
		for i, author := range entry.Authors {
			authors[i] = author.Name
		}

		// Parse published date
		publishedDate, err := time.Parse("2006-01-02T15:04:05Z", entry.Published)
		if err != nil {
			publishedDate = time.Now() // Fallback
		}

		// Calculate credibility based on recency
		credibility := a.calculateCredibility(publishedDate)

		evidence := &domain.Evidence{
			Source:        "arXiv",
			Title:         entry.Title,
			Authors:       authors,
			URL:           entry.ID,
			Snippet:       entry.Summary,
			Credibility:   credibility,
			Relevance:     0.7, // Default relevance
			PublishedDate: &publishedDate,
			SourceType:    domain.SourceTypePreprint,
			CitationCount: 0, // arXiv API doesn't provide citation count
			Metadata: map[string]string{
				"arxiv_id": entry.ID,
			},
		}

		evidences = append(evidences, evidence)
	}

	return evidences, nil
}

// calculateCredibility calculates credibility score based on publication date
func (a *ArxivClient) calculateCredibility(publishedDate time.Time) float64 {
	// Base credibility for preprints
	baseScore := 0.75

	// Boost recent papers (within 2 years)
	yearsOld := time.Since(publishedDate).Hours() / (24 * 365)
	if yearsOld < 2 {
		return baseScore + 0.1
	}

	// Penalize old papers
	if yearsOld > 5 {
		return baseScore - 0.1
	}

	return baseScore
}
