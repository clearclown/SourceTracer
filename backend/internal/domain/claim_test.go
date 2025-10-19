package domain

import (
	"encoding/json"
	"testing"
)

func TestClaimType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		ct       ClaimType
		expected string
	}{
		{"Opinion", Opinion, `"opinion"`},
		{"Fact", Fact, `"fact"`},
		{"Mixed", Mixed, `"mixed"`},
		{"Unclear", Unclear, `"unclear"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.ct)
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestClaimType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ClaimType
	}{
		{"Opinion", `"opinion"`, Opinion},
		{"Fact", `"fact"`, Fact},
		{"Mixed", `"mixed"`, Mixed},
		{"Unclear", `"unclear"`, Unclear},
		{"Unknown", `"unknown"`, Unclear},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct ClaimType
			err := json.Unmarshal([]byte(tt.input), &ct)
			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}
			if ct != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, ct)
			}
		})
	}
}

func TestClaim_JSON_RoundTrip(t *testing.T) {
	claim := &Claim{
		ID:         "test-id",
		Text:       "Test claim",
		Type:       Opinion,
		Confidence: 0.9,
		Position:   Position{Start: 0, End: 10},
		Evidences:  []Evidence{},
	}

	// Marshal
	data, err := json.Marshal(claim)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Check that type is a string
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	typeValue, ok := result["type"].(string)
	if !ok {
		t.Fatalf("Type is not a string: %T", result["type"])
	}
	if typeValue != "opinion" {
		t.Errorf("Expected type 'opinion', got '%s'", typeValue)
	}

	// Unmarshal back
	var unmarshaled Claim
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.Type != Opinion {
		t.Errorf("Expected Opinion, got %v", unmarshaled.Type)
	}
}
