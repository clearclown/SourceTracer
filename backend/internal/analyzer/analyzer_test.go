package analyzer

import (
	"context"
	"testing"

	"github.com/yourusername/sourcetracer/internal/classifier"
	"github.com/yourusername/sourcetracer/internal/domain"
)

// 🔴 Red: Analyzer creation test
func TestAnalyzer_Creation(t *testing.T) {
	c := classifier.NewClassifier()
	analyzer := NewAnalyzer(c, nil)

	if analyzer == nil {
		t.Fatal("Expected non-nil analyzer")
	}
}

// 🔴 Red: Analyze text and extract claims
func TestAnalyzer_AnalyzeText(t *testing.T) {
	c := classifier.NewClassifier()
	analyzer := NewAnalyzer(c, nil)

	ctx := context.Background()
	text := "Climate change is accelerating. Python is the best language."

	result, err := analyzer.Analyze(ctx, text, AnalyzeOptions{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// 2つのclaimが抽出されるべき
	if len(result.Claims) < 1 {
		t.Errorf("Expected at least 1 claim, got %d", len(result.Claims))
	}
}

// 🔴 Red: Extract claims from text
func TestAnalyzer_ExtractClaims(t *testing.T) {
	c := classifier.NewClassifier()
	analyzer := NewAnalyzer(c, nil)

	text := "Climate change is real. It is accelerating."

	claims := analyzer.extractClaims(text)

	if len(claims) != 2 {
		t.Fatalf("Expected 2 claims, got %d", len(claims))
	}

	if claims[0].Text != "Climate change is real." {
		t.Errorf("Expected first claim 'Climate change is real.', got %q", claims[0].Text)
	}

	if claims[1].Text != "It is accelerating." {
		t.Errorf("Expected second claim 'It is accelerating.', got %q", claims[1].Text)
	}
}

// 🔴 Red: Classify claims
func TestAnalyzer_ClassifyClaims(t *testing.T) {
	c := classifier.NewClassifier()
	analyzer := NewAnalyzer(c, nil)

	claim1 := domain.NewClaim("Python is the best language", 0, 27)
	claim2 := domain.NewClaim("The Earth orbits the Sun", 0, 24)

	claims := []*domain.Claim{claim1, claim2}

	ctx := context.Background()
	err := analyzer.classifyClaims(ctx, claims)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// claim1 should be Opinion
	if claims[0].Type != domain.Opinion {
		t.Errorf("Expected Opinion for claim1, got %v", claims[0].Type)
	}

	// claim2 should be Fact
	if claims[1].Type != domain.Fact {
		t.Errorf("Expected Fact for claim2, got %v", claims[1].Type)
	}

	// Both should have confidence > 0
	if claims[0].Confidence <= 0 {
		t.Error("Expected positive confidence for claim1")
	}
	if claims[1].Confidence <= 0 {
		t.Error("Expected positive confidence for claim2")
	}
}

// 🔴 Red: Full analysis with options
func TestAnalyzer_AnalyzeWithOptions(t *testing.T) {
	c := classifier.NewClassifier()
	analyzer := NewAnalyzer(c, nil)

	ctx := context.Background()
	text := "Python is great. Climate change is happening."

	opts := AnalyzeOptions{
		IncludeEvidences: false,
		MaxClaims:        10,
	}

	result, err := analyzer.Analyze(ctx, text, opts)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Claims) == 0 {
		t.Error("Expected some claims")
	}

	// Evidenceは検索しないので空のはず
	for _, claim := range result.Claims {
		if len(claim.Evidences) > 0 {
			t.Error("Expected no evidences when IncludeEvidences=false")
		}
	}
}
