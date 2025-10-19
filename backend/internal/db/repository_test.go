package db

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/sourcetracer/internal/analyzer"
	"github.com/yourusername/sourcetracer/internal/domain"
)

// 🔴 Red: Test repository creation
func TestRepository_New(t *testing.T) {
	// This test will be skipped if no database is available
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// For now, we'll test with nil and ensure it doesn't panic
	repo := NewRepository(nil)
	if repo == nil {
		t.Fatal("Expected non-nil repository")
	}
}

// 🔴 Red: Test saving analysis
func TestRepository_SaveAnalysis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Mock database connection for testing
	repo := NewRepository(nil)

	result := &analyzer.AnalyzeResult{
		AnalysisID:         "test_123",
		Claims:             []*domain.Claim{},
		OverallCredibility: 0.85,
		ProcessingTimeMS:   100,
		CreatedAt:          time.Now(),
	}

	ctx := context.Background()
	id, err := repo.SaveAnalysis(ctx, "Test text", result)

	// For now, we expect this to work with mock
	if err != nil && !testing.Short() {
		t.Errorf("Expected no error in mock mode, got %v", err)
	}

	if id == "" && !testing.Short() {
		t.Error("Expected non-empty ID")
	}
}

// 🔴 Red: Test retrieving analysis
func TestRepository_GetAnalysis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	repo := NewRepository(nil)
	ctx := context.Background()

	_, err := repo.GetAnalysis(ctx, "test_id")

	// We expect an error with mock database
	if err == nil && !testing.Short() {
		t.Error("Expected error with non-existent ID")
	}
}

// 🔴 Red: Test listing recent analyses
func TestRepository_ListRecentAnalyses(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	repo := NewRepository(nil)
	ctx := context.Background()

	analyses, err := repo.ListRecentAnalyses(ctx, 10, 0)

	if err != nil && !testing.Short() {
		t.Errorf("Expected no error, got %v", err)
	}

	if analyses == nil {
		t.Error("Expected non-nil slice")
	}
}

// 🔴 Red: Test mock database implementation
func TestMockDB_Insert(t *testing.T) {
	db := NewMockDB()

	if db == nil {
		t.Fatal("Expected non-nil mock database")
	}

	// Test basic operations
	ctx := context.Background()
	id, err := db.Insert(ctx, "test_table", map[string]interface{}{
		"name": "test",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id == "" {
		t.Error("Expected non-empty ID")
	}
}

// 🔴 Red: Test transaction support
func TestRepository_Transaction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	repo := NewRepository(nil)
	ctx := context.Background()

	err := repo.WithTransaction(ctx, func(ctx context.Context) error {
		// Test transaction operations
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
