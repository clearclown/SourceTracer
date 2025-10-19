package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ClaimType represents the classification of a claim
type ClaimType int

const (
	Unclear ClaimType = iota // デフォルト: 未分類
	Opinion                  // 意見
	Fact                     // 事実
	Mixed                    // 意見と事実の混在
)

// String returns the string representation of ClaimType
func (ct ClaimType) String() string {
	switch ct {
	case Opinion:
		return "opinion"
	case Fact:
		return "fact"
	case Mixed:
		return "mixed"
	case Unclear:
		return "unclear"
	default:
		return "unclear"
	}
}

// Position represents the position of a claim in the original text
type Position struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Claim represents a single claim extracted from text
type Claim struct {
	ID         string     `json:"id"`
	Text       string     `json:"text"`
	Type       ClaimType  `json:"type"`
	Confidence float64    `json:"confidence"` // 0.0-1.0
	Position   Position   `json:"position"`
	Evidences  []Evidence `json:"evidences"`
	Summary    string     `json:"summary,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// NewClaim creates a new Claim with default values
func NewClaim(text string, start, end int) *Claim {
	return &Claim{
		ID:         generateClaimID(),
		Text:       text,
		Type:       Unclear,
		Confidence: 0.0,
		Position: Position{
			Start: start,
			End:   end,
		},
		Evidences: []Evidence{},
		CreatedAt: time.Now(),
	}
}

// generateClaimID generates a unique ID for a claim
func generateClaimID() string {
	return fmt.Sprintf("clm_%s", uuid.New().String()[:8])
}
