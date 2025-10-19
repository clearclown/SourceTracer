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

// setupTestRouter creates a test router
func setupTestRouter() *gin.Engine {
	c := classifier.NewClassifier()
	a := analyzer.NewAnalyzer(c, nil)

	handler := NewHandler(a)
	return handler.SetupRouter()
}
