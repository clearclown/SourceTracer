package search

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// SearchClient interface for all search providers
type SearchClient interface {
	Search(ctx context.Context, query string, maxResults int) ([]*domain.Evidence, error)
}

// Aggregator aggregates results from multiple search providers
type Aggregator struct {
	clients map[string]SearchClient
}

// NewAggregator creates a new search aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		clients: make(map[string]SearchClient),
	}
}

// AddClient adds a search client with a name
func (a *Aggregator) AddClient(name string, client SearchClient) {
	a.clients[name] = client
}

// SearchAll searches across all registered providers in parallel
func (a *Aggregator) SearchAll(ctx context.Context, query string, maxResultsPerProvider int) ([]*domain.Evidence, error) {
	if len(a.clients) == 0 {
		return []*domain.Evidence{}, nil
	}

	// Channel to collect results
	type result struct {
		evidences []*domain.Evidence
		err       error
		source    string
	}

	results := make(chan result, len(a.clients))
	var wg sync.WaitGroup

	// Launch parallel searches
	for name, client := range a.clients {
		wg.Add(1)
		go func(n string, c SearchClient) {
			defer wg.Done()
			evidences, err := c.Search(ctx, query, maxResultsPerProvider)
			results <- result{
				evidences: evidences,
				err:       err,
				source:    n,
			}
		}(name, client)
	}

	// Wait for all searches to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	var allEvidences []*domain.Evidence
	var errors []error

	for res := range results {
		if res.err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", res.source, res.err))
			continue
		}
		allEvidences = append(allEvidences, res.evidences...)
	}

	// Sort by credibility (highest first)
	sort.Slice(allEvidences, func(i, j int) bool {
		return allEvidences[i].Credibility > allEvidences[j].Credibility
	})

	// If all searches failed, return error
	if len(allEvidences) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all searches failed: %v", errors)
	}

	return allEvidences, nil
}

// SearchWithProviders searches using specific providers
func (a *Aggregator) SearchWithProviders(ctx context.Context, query string, maxResultsPerProvider int, providers []string) ([]*domain.Evidence, error) {
	if len(providers) == 0 {
		return a.SearchAll(ctx, query, maxResultsPerProvider)
	}

	// Validate providers exist
	selectedClients := make(map[string]SearchClient)
	for _, name := range providers {
		if client, ok := a.clients[name]; ok {
			selectedClients[name] = client
		}
	}

	if len(selectedClients) == 0 {
		return nil, fmt.Errorf("no valid providers found: %v", providers)
	}

	// Temporarily replace clients
	originalClients := a.clients
	a.clients = selectedClients
	defer func() { a.clients = originalClients }()

	return a.SearchAll(ctx, query, maxResultsPerProvider)
}
