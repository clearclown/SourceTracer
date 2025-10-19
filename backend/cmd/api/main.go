package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/sourcetracer/internal/analyzer"
	"github.com/yourusername/sourcetracer/internal/api"
	"github.com/yourusername/sourcetracer/internal/classifier"
	"github.com/yourusername/sourcetracer/internal/config"
)

const version = "v0.3.0-alpha"

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("🚀 Starting SourceTracer API Server %s", version)
	log.Printf("📝 Environment: %s", cfg.AppEnv)
	log.Printf("🔌 Port: %d", cfg.AppPort)
	log.Printf("🤖 Default LLM: %s", cfg.DefaultLLMProvider)

	// Initialize components
	c := classifier.NewClassifier()
	a := analyzer.NewAnalyzer(c, nil) // TODO: Add search clients

	// Create API handler
	handler := api.NewHandler(a)
	router := handler.SetupRouter()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.AppPort)
		log.Printf("✅ Server listening on %s", addr)
		log.Printf("📖 API documentation: http://localhost:%d/api/v1/health", cfg.AppPort)

		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("\n🛑 Shutting down server...")
	log.Println("✅ Server gracefully stopped")
}
