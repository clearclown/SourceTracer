package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/sourcetracer/internal/analyzer"
)

// Handler handles HTTP requests
type Handler struct {
	analyzer *analyzer.Analyzer
}

// NewHandler creates a new Handler
func NewHandler(a *analyzer.Analyzer) *Handler {
	return &Handler{
		analyzer: a,
	}
}

// SetupRouter sets up the Gin router
func (h *Handler) SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(CORSMiddleware())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", h.HealthHandler)
		v1.POST("/analyze", h.AnalyzeHandler)
	}

	return router
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string                 `json:"status"`
	Version string                 `json:"version"`
	Time    string                 `json:"time"`
	Services map[string]string     `json:"services,omitempty"`
}

// HealthHandler handles health check
func (h *Handler) HealthHandler(c *gin.Context) {
	response := HealthResponse{
		Status:  "ok",
		Version: "v0.3.0-alpha",
		Time:    time.Now().Format(time.RFC3339),
		Services: map[string]string{
			"analyzer": "ok",
		},
	}

	c.JSON(http.StatusOK, response)
}

// AnalyzeRequest represents the analyze request
type AnalyzeRequest struct {
	Text    string         `json:"text" binding:"required"`
	Options AnalyzeOptions `json:"options"`
}

// AnalyzeOptions represents analysis options
type AnalyzeOptions struct {
	IncludeEvidences bool     `json:"include_evidences"`
	MaxClaims        int      `json:"max_claims"`
	MaxEvidences     int      `json:"max_evidences"`
	SearchEngines    []string `json:"search_engines"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success  bool                   `json:"success"`
	Data     interface{}            `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Error   ErrorDetail            `json:"error"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// AnalyzeHandler handles text analysis
func (h *Handler) AnalyzeHandler(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ErrorDetail{
				Code:    "INVALID_INPUT",
				Message: "Invalid request body",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
		return
	}

	// Validate text
	if len(req.Text) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ErrorDetail{
				Code:    "INVALID_INPUT",
				Message: "Text input is required",
				Details: map[string]interface{}{
					"field":  "text",
					"reason": "cannot be empty",
				},
			},
			Metadata: map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
		return
	}

	// Convert to analyzer options
	analyzerOpts := analyzer.AnalyzeOptions{
		IncludeEvidences: req.Options.IncludeEvidences,
		MaxClaims:        req.Options.MaxClaims,
		MaxEvidences:     req.Options.MaxEvidences,
		SearchEngines:    req.Options.SearchEngines,
	}

	// Analyze
	result, err := h.analyzer.Analyze(c.Request.Context(), req.Text, analyzerOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to analyze text",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			},
			Metadata: map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    result,
		Metadata: map[string]interface{}{
			"timestamp":        time.Now().Format(time.RFC3339),
			"processing_time_ms": result.ProcessingTimeMS,
		},
	})
}
