package search

// contains checks if string s contains substr
func contains(s, substr string) bool {
	if len(s) < len(substr) || s == "" || substr == "" {
		return s == substr
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
