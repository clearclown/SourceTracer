package search

import (
	"context"
	"testing"
)

// 🔴 Red: Test Google Custom Search client creation
func TestGoogleClient_Creation(t *testing.T) {
	client := NewGoogleClient("test_api_key", "test_cx")
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// 🔴 Red: Test Google search URL building
func TestGoogleClient_BuildSearchURL(t *testing.T) {
	client := NewGoogleClient("test_api_key", "test_cx")
	url := client.buildSearchURL("machine learning", 10)

	if url == "" {
		t.Error("Expected non-empty URL")
	}

	if !contains(url, "key=") {
		t.Error("URL should contain API key parameter")
	}

	if !contains(url, "cx=") {
		t.Error("URL should contain custom search engine ID")
	}

	if !contains(url, "q=") {
		t.Error("URL should contain query parameter")
	}
}

// 🔴 Red: Test Google API response parsing
func TestGoogleClient_ParseResponse(t *testing.T) {
	client := NewGoogleClient("test_api_key", "test_cx")

	// Sample Google Custom Search JSON response
	jsonResponse := `{
		"items": [
			{
				"title": "Understanding Machine Learning",
				"link": "https://example.com/ml-article",
				"snippet": "A comprehensive guide to machine learning algorithms...",
				"displayLink": "example.com"
			}
		]
	}`

	evidences, err := client.parseResponse([]byte(jsonResponse))
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(evidences) == 0 {
		t.Fatal("Expected at least one evidence")
	}

	evidence := evidences[0]

	if evidence.Title != "Understanding Machine Learning" {
		t.Errorf("Expected title 'Understanding Machine Learning', got %s", evidence.Title)
	}

	if evidence.Source != "Google" {
		t.Errorf("Expected source 'Google', got %s", evidence.Source)
	}

	if evidence.SourceType != "unknown" {
		t.Errorf("Expected source type 'unknown', got %v", evidence.SourceType)
	}
}

// 🔴 Red: Test Google search (integration test)
func TestGoogleClient_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if no API key provided
	apiKey := "test_key"
	cx := "test_cx"

	if apiKey == "test_key" {
		t.Skip("Skipping test - no real API key provided")
	}

	client := NewGoogleClient(apiKey, cx)
	ctx := context.Background()

	evidences, err := client.Search(ctx, "test query", 3)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(evidences) == 0 {
		t.Error("Expected at least one result")
	}
}
