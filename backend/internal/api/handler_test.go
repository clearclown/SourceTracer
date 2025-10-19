package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/sourcetracer/internal/analyzer"
	"github.com/yourusername/sourcetracer/internal/classifier"
)

// 🔴 Red: Test health check endpoint
func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

// 🔴 Red: Test analyze endpoint
func TestAnalyzeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	requestBody := map[string]interface{}{
		"text": "Python is the best language.",
		"options": map[string]interface{}{
			"include_evidences": false,
			"max_claims":        10,
		},
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}

	data := response["data"].(map[string]interface{})
	claims := data["claims"].([]interface{})

	if len(claims) == 0 {
		t.Error("Expected at least 1 claim")
	}
}

// 🔴 Red: Test analyze endpoint with invalid input
func TestAnalyzeHandler_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	requestBody := map[string]interface{}{
		"text": "", // Empty text
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["success"].(bool) {
		t.Error("Expected success to be false")
	}
}

// 🔴 Red: Test CORS middleware
func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	req := httptest.NewRequest("OPTIONS", "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected CORS headers to be set")
	}
}

// 🔴 Red: Test history endpoint
func TestHistoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/v1/history", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}

	data := response["data"].(map[string]interface{})
	if _, ok := data["analyses"]; !ok {
		t.Error("Expected 'analyses' field in response")
	}

	if _, ok := data["limit"]; !ok {
		t.Error("Expected 'limit' field in response")
	}

	if _, ok := data["offset"]; !ok {
		t.Error("Expected 'offset' field in response")
	}
}

// 🔴 Red: Test history endpoint with pagination
func TestHistoryHandler_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/v1/history?limit=5&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data := response["data"].(map[string]interface{})

	if limit := data["limit"].(float64); limit != 5 {
		t.Errorf("Expected limit 5, got %v", limit)
	}

	if offset := data["offset"].(float64); offset != 10 {
		t.Errorf("Expected offset 10, got %v", offset)
	}
}

// 🔴 Red: Test get analysis by ID endpoint
func TestGetAnalysisHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	// Test with non-existent ID (should return 404)
	req := httptest.NewRequest("GET", "/api/v1/history/non-existent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["success"].(bool) {
		t.Error("Expected success to be false")
	}

	errorData := response["error"].(map[string]interface{})
	if errorData["code"] != "NOT_FOUND" {
		t.Errorf("Expected error code 'NOT_FOUND', got %v", errorData["code"])
	}
}

// 🔴 Red: Test get analysis with empty ID
func TestGetAnalysisHandler_EmptyID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupTestRouter()

	// Test with empty ID
	req := httptest.NewRequest("GET", "/api/v1/history/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should redirect to history list endpoint or return 404
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status 200, 301, or 404, got %d", w.Code)
	}
}

// setupTestRouter creates a test router
func setupTestRouter() *gin.Engine {
	c := classifier.NewClassifier()
	a := analyzer.NewAnalyzer(c, nil)

	handler := NewHandler(a)
	return handler.SetupRouter()
}
