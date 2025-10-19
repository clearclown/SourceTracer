package domain

import (
	"testing"
)

// 🔴 Red: 失敗するテストを書く
func TestClaim_Creation(t *testing.T) {
	// Arrange
	text := "Climate change is accelerating"
	startPos := 0
	endPos := 31

	// Act
	claim := NewClaim(text, startPos, endPos)

	// Assert
	if claim.Text != text {
		t.Errorf("Expected text %q, got %q", text, claim.Text)
	}
	if claim.Position.Start != startPos {
		t.Errorf("Expected start position %d, got %d", startPos, claim.Position.Start)
	}
	if claim.Position.End != endPos {
		t.Errorf("Expected end position %d, got %d", endPos, claim.Position.End)
	}
	if claim.ID == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestClaim_DefaultValues(t *testing.T) {
	claim := NewClaim("Test claim", 0, 10)

	// デフォルトでは型が未分類
	if claim.Type != Unclear {
		t.Errorf("Expected type Unclear, got %v", claim.Type)
	}

	// デフォルトでは信頼度0
	if claim.Confidence != 0.0 {
		t.Errorf("Expected confidence 0.0, got %f", claim.Confidence)
	}

	// エビデンスリストは空
	if len(claim.Evidences) != 0 {
		t.Errorf("Expected empty evidences, got %d items", len(claim.Evidences))
	}
}

func TestClaimType_String(t *testing.T) {
	tests := []struct {
		claimType ClaimType
		expected  string
	}{
		{Opinion, "opinion"},
		{Fact, "fact"},
		{Mixed, "mixed"},
		{Unclear, "unclear"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.claimType.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, tt.claimType.String())
			}
		})
	}
}
