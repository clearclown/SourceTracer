package search

import (
	"context"
	"testing"
)

// 🔴 Red: Test arXiv client creation
func TestArxivClient_Creation(t *testing.T) {
	client := NewArxivClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// 🔴 Red: Test arXiv search URL building
func TestArxivClient_BuildSearchURL(t *testing.T) {
	client := NewArxivClient()
	url := client.buildSearchURL("quantum computing", 10)

	if url == "" {
		t.Error("Expected non-empty URL")
	}

	// Should contain query parameter
	if !contains(url, "search_query") {
		t.Error("URL should contain search_query parameter")
	}

	// Should contain max_results
	if !contains(url, "max_results") {
		t.Error("URL should contain max_results parameter")
	}
}

// 🔴 Red: Test arXiv API response parsing
func TestArxivClient_ParseResponse(t *testing.T) {
	client := NewArxivClient()

	// Sample arXiv XML response
	xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <entry>
    <id>http://arxiv.org/abs/2301.00001v1</id>
    <title>Sample Paper on Quantum Computing</title>
    <author><name>John Doe</name></author>
    <author><name>Jane Smith</name></author>
    <summary>This is a sample abstract about quantum computing.</summary>
    <published>2023-01-01T00:00:00Z</published>
  </entry>
</feed>`

	evidences, err := client.parseResponse([]byte(xmlResponse))
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(evidences) == 0 {
		t.Fatal("Expected at least one evidence")
	}

	evidence := evidences[0]

	if evidence.Title != "Sample Paper on Quantum Computing" {
		t.Errorf("Expected title 'Sample Paper on Quantum Computing', got %s", evidence.Title)
	}

	if len(evidence.Authors) != 2 {
		t.Errorf("Expected 2 authors, got %d", len(evidence.Authors))
	}

	if evidence.SourceType != "preprint" {
		t.Errorf("Expected source type 'preprint', got %s", evidence.SourceType)
	}
}

// 🔴 Red: Test arXiv search (integration test)
func TestArxivClient_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewArxivClient()
	ctx := context.Background()

	evidences, err := client.Search(ctx, "machine learning", 3)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(evidences) == 0 {
		t.Error("Expected at least one result")
	}

	// Check first result has required fields
	if len(evidences) > 0 {
		e := evidences[0]
		if e.Title == "" {
			t.Error("Expected non-empty title")
		}
		if e.URL == "" {
			t.Error("Expected non-empty URL")
		}
		if e.Source != "arXiv" {
			t.Errorf("Expected source 'arXiv', got %s", e.Source)
		}
	}
}
