package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/yourusername/sourcetracer/internal/analyzer"
	"github.com/yourusername/sourcetracer/internal/domain"
)

// Database interface for abstraction
type Database interface {
	Insert(ctx context.Context, table string, data map[string]interface{}) (string, error)
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
	BeginTx(ctx context.Context) (Transaction, error)
}

// Transaction interface
type Transaction interface {
	Commit() error
	Rollback() error
	Exec(ctx context.Context, query string, args ...interface{}) error
}

// Repository handles database operations
type Repository struct {
	db Database
}

// NewRepository creates a new repository
func NewRepository(db Database) *Repository {
	// If db is nil, use mock database for testing
	if db == nil {
		db = NewMockDB()
	}
	return &Repository{
		db: db,
	}
}

// SaveAnalysis saves an analysis result to database
func (r *Repository) SaveAnalysis(ctx context.Context, originalText string, result *analyzer.AnalyzeResult) (string, error) {
	// Generate UUID for analysis
	analysisID := uuid.New().String()

	// Insert analysis record
	data := map[string]interface{}{
		"id":                  analysisID,
		"original_text":       originalText,
		"overall_credibility": result.OverallCredibility,
		"processing_time_ms":  result.ProcessingTimeMS,
		"created_at":          result.CreatedAt,
	}

	_, err := r.db.Insert(ctx, "analyses", data)
	if err != nil {
		return "", fmt.Errorf("failed to insert analysis: %w", err)
	}

	// Save claims
	for _, claim := range result.Claims {
		claimData := map[string]interface{}{
			"id":             uuid.New().String(),
			"analysis_id":    analysisID,
			"text":           claim.Text,
			"claim_type":     claim.Type.String(),
			"confidence":     claim.Confidence,
			"position_start": claim.Position.Start,
			"position_end":   claim.Position.End,
			"summary":        claim.Summary,
			"created_at":     claim.CreatedAt,
		}

		_, err := r.db.Insert(ctx, "claims", claimData)
		if err != nil {
			return "", fmt.Errorf("failed to insert claim: %w", err)
		}
	}

	return analysisID, nil
}

// GetAnalysis retrieves an analysis by ID
func (r *Repository) GetAnalysis(ctx context.Context, id string) (*analyzer.AnalyzeResult, error) {
	query := `
		SELECT id, original_text, overall_credibility, processing_time_ms, created_at
		FROM analyses
		WHERE id = $1
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query analysis: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("analysis not found: %s", id)
	}

	// Parse result (simplified for mock)
	row := rows[0]
	result := &analyzer.AnalyzeResult{
		AnalysisID:         fmt.Sprintf("%v", row["id"]),
		Claims:             []*domain.Claim{},
		OverallCredibility: 0.0,
		ProcessingTimeMS:   0,
		CreatedAt:          time.Now(),
	}

	return result, nil
}

// ListRecentAnalyses retrieves recent analyses
func (r *Repository) ListRecentAnalyses(ctx context.Context, limit, offset int) ([]*analyzer.AnalyzeResult, error) {
	query := `
		SELECT id, original_text, overall_credibility, processing_time_ms, created_at
		FROM analyses
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query analyses: %w", err)
	}

	results := make([]*analyzer.AnalyzeResult, 0, len(rows))
	for _, row := range rows {
		result := &analyzer.AnalyzeResult{
			AnalysisID: fmt.Sprintf("%v", row["id"]),
			Claims:     []*domain.Claim{},
		}
		results = append(results, result)
	}

	return results, nil
}

// WithTransaction executes a function within a transaction
func (r *Repository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Ignore rollback error during panic recovery
			_ = tx.Rollback() //nolint:errcheck
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		// Ignore rollback error when function fails
		_ = tx.Rollback() //nolint:errcheck
		return err
	}

	return tx.Commit()
}

// MockDB implements Database interface for testing
type MockDB struct {
	data map[string][]map[string]interface{}
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		data: make(map[string][]map[string]interface{}),
	}
}

// Insert mock implementation
func (m *MockDB) Insert(ctx context.Context, table string, data map[string]interface{}) (string, error) {
	// Generate ID if not present
	if _, ok := data["id"]; !ok {
		data["id"] = uuid.New().String()
	}

	if m.data[table] == nil {
		m.data[table] = make([]map[string]interface{}, 0)
	}

	m.data[table] = append(m.data[table], data)
	id, ok := data["id"].(string)
	if !ok {
		return "", fmt.Errorf("id is not a string")
	}
	return id, nil
}

// Query mock implementation
func (m *MockDB) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	// Simple mock - return empty result
	return []map[string]interface{}{}, nil
}

// Exec mock implementation
func (m *MockDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	return nil
}

// BeginTx mock implementation
func (m *MockDB) BeginTx(ctx context.Context) (Transaction, error) {
	return &MockTransaction{}, nil
}

// MockTransaction implements Transaction interface
type MockTransaction struct {
	committed bool
}

// Commit mock implementation
func (t *MockTransaction) Commit() error {
	t.committed = true
	return nil
}

// Rollback mock implementation
func (t *MockTransaction) Rollback() error {
	return nil
}

// Exec mock implementation
func (t *MockTransaction) Exec(ctx context.Context, query string, args ...interface{}) error {
	return nil
}

// PostgresDB implementation moved to postgres.go
