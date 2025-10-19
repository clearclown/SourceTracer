package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresDB implements Database interface for PostgreSQL
type PostgresDB struct {
	conn *sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(connString string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{conn: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.conn.Close()
}

// Insert implements Database.Insert
func (p *PostgresDB) Insert(ctx context.Context, table string, data map[string]interface{}) (string, error) {
	// For analyses table with UUID return
	if table == "analyses" {
		var id string
		err := p.conn.QueryRowContext(ctx,
			`INSERT INTO analyses (original_text, overall_credibility, processing_time_ms, created_at)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			data["original_text"],
			data["overall_credibility"],
			data["processing_time_ms"],
			data["created_at"],
		).Scan(&id)
		if err != nil {
			return "", fmt.Errorf("failed to insert into analyses: %w", err)
		}
		return id, nil
	}

	// For claims table
	if table == "claims" {
		var id string
		err := p.conn.QueryRowContext(ctx,
			`INSERT INTO claims (analysis_id, text, claim_type, confidence, position_start, position_end, summary, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 RETURNING id`,
			data["analysis_id"],
			data["text"],
			data["claim_type"],
			data["confidence"],
			data["position_start"],
			data["position_end"],
			data["summary"],
			data["created_at"],
		).Scan(&id)
		if err != nil {
			return "", fmt.Errorf("failed to insert into claims: %w", err)
		}
		return id, nil
	}

	return "", fmt.Errorf("unsupported table: %s", table)
}

// Query implements Database.Query
func (p *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// Exec implements Database.Exec
func (p *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	return nil
}

// PostgresTx implements Transaction interface
type PostgresTx struct {
	tx *sql.Tx
}

// BeginTx implements Database.BeginTx
func (p *PostgresDB) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &PostgresTx{tx: tx}, nil
}

// Commit implements Transaction.Commit
func (t *PostgresTx) Commit() error {
	return t.tx.Commit()
}

// Rollback implements Transaction.Rollback
func (t *PostgresTx) Rollback() error {
	return t.tx.Rollback()
}

// Exec implements Transaction.Exec
func (t *PostgresTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.tx.ExecContext(ctx, query, args...)
	return err
}
