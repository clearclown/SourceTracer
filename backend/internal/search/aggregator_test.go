package search

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/sourcetracer/internal/domain"
)

// MockSearchClient for testing
type MockSearchClient struct {
	results []*domain.Evidence
	err     error
}

func (m *MockSearchClient) Search(ctx context.Context, query string, maxResults int) ([]*domain.Evidence, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

// 🔴 Red: Test aggregator creation
func TestAggregator_Creation(t *testing.T) {
	agg := NewAggregator()
	if agg == nil {
		t.Fatal("Expected non-nil aggregator")
	}
}

// 🔴 Red: Test adding clients
func TestAggregator_AddClient(t *testing.T) {
	agg := NewAggregator()
	mock := &MockSearchClient{}

	agg.AddClient("test", mock)

	if len(agg.clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(agg.clients))
	}
}

// 🔴 Red: Test searching with no clients
func TestAggregator_SearchAll_NoClients(t *testing.T) {
	agg := NewAggregator()
	ctx := context.Background()

	results, err := agg.SearchAll(ctx, "test", 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// 🔴 Red: Test searching with single client
func TestAggregator_SearchAll_SingleClient(t *testing.T) {
	agg := NewAggregator()
	mock := &MockSearchClient{
		results: []*domain.Evidence{
			{
				Title:       "Test Evidence",
				Credibility: 0.8,
			},
		},
	}

	agg.AddClient("test", mock)
	ctx := context.Background()

	results, err := agg.SearchAll(ctx, "query", 5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Title != "Test Evidence" {
		t.Errorf("Expected title 'Test Evidence', got %s", results[0].Title)
	}
}

// 🔴 Red: Test searching with multiple clients
func TestAggregator_SearchAll_MultipleClients(t *testing.T) {
	agg := NewAggregator()

	mock1 := &MockSearchClient{
		results: []*domain.Evidence{
			{Title: "Result 1", Credibility: 0.9},
		},
	}

	mock2 := &MockSearchClient{
		results: []*domain.Evidence{
			{Title: "Result 2", Credibility: 0.7},
		},
	}

	agg.AddClient("client1", mock1)
	agg.AddClient("client2", mock2)

	ctx := context.Background()
	results, err := agg.SearchAll(ctx, "query", 5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Should be sorted by credibility (highest first)
	if results[0].Credibility < results[1].Credibility {
		t.Error("Results not sorted by credibility")
	}
}

// 🔴 Red: Test handling errors
func TestAggregator_SearchAll_WithErrors(t *testing.T) {
	agg := NewAggregator()

	mock1 := &MockSearchClient{
		err: errors.New("search failed"),
	}

	mock2 := &MockSearchClient{
		results: []*domain.Evidence{
			{Title: "Result 2", Credibility: 0.7},
		},
	}

	agg.AddClient("failing", mock1)
	agg.AddClient("working", mock2)

	ctx := context.Background()
	results, err := agg.SearchAll(ctx, "query", 5)

	// Should succeed with partial results
	if err != nil {
		t.Errorf("Expected no error with partial results, got %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result from working client, got %d", len(results))
	}
}

// 🔴 Red: Test all clients failing
func TestAggregator_SearchAll_AllFail(t *testing.T) {
	agg := NewAggregator()

	mock1 := &MockSearchClient{
		err: errors.New("search failed 1"),
	}

	mock2 := &MockSearchClient{
		err: errors.New("search failed 2"),
	}

	agg.AddClient("client1", mock1)
	agg.AddClient("client2", mock2)

	ctx := context.Background()
	results, err := agg.SearchAll(ctx, "query", 5)

	if err == nil {
		t.Error("Expected error when all clients fail")
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// 🔴 Red: Test searching with specific providers
func TestAggregator_SearchWithProviders(t *testing.T) {
	agg := NewAggregator()

	mock1 := &MockSearchClient{
		results: []*domain.Evidence{{Title: "Result 1"}},
	}

	mock2 := &MockSearchClient{
		results: []*domain.Evidence{{Title: "Result 2"}},
	}

	agg.AddClient("arxiv", mock1)
	agg.AddClient("google", mock2)

	ctx := context.Background()

	// Search only with arxiv
	results, err := agg.SearchWithProviders(ctx, "query", 5, []string{"arxiv"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Title != "Result 1" {
		t.Errorf("Expected 'Result 1', got %s", results[0].Title)
	}
}
