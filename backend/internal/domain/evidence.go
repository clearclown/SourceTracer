package domain

import (
	"time"
)

// SourceType represents the type of source
type SourceType string

const (
	SourceTypePeerReviewedJournal SourceType = "peer_reviewed_journal"
	SourceTypePreprint            SourceType = "preprint"
	SourceTypeNews                SourceType = "news"
	SourceTypeBlog                SourceType = "blog"
	SourceTypeGovernment          SourceType = "government"
	SourceTypeUnknown             SourceType = "unknown"
)

// Evidence represents supporting or contradicting evidence for a claim
type Evidence struct {
	ID            string            `json:"id"`
	Source        string            `json:"source"`         // e.g., "Semantic Scholar"
	Title         string            `json:"title"`
	Authors       []string          `json:"authors,omitempty"`
	URL           string            `json:"url"`
	Snippet       string            `json:"snippet"`
	Credibility   float64           `json:"credibility"`   // 0.0-1.0
	Relevance     float64           `json:"relevance"`     // 0.0-1.0
	PublishedDate *time.Time        `json:"published_date,omitempty"`
	SourceType    SourceType        `json:"source_type"`
	CitationCount int               `json:"citation_count,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}
