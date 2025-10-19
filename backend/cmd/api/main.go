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
	"github.com/yourusername/sourcetracer/internal/db"
	"github.com/yourusername/sourcetracer/internal/search"
)

const version = "v0.7.0-alpha"

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

	// Initialize search aggregator
	aggregator := search.NewAggregator()

	// Add search clients
	aggregator.AddClient("semantic_scholar", search.NewSemanticScholarClient(""))
	aggregator.AddClient("arxiv", search.NewArxivClient())

	// Add Google if API key provided
	if googleKey := os.Getenv("GOOGLE_API_KEY"); googleKey != "" {
		googleCX := os.Getenv("GOOGLE_CX")
		if googleCX != "" {
			aggregator.AddClient("google", search.NewGoogleClient(googleKey, googleCX))
			log.Printf("✅ Google Custom Search enabled")
		}
	}

	log.Printf("🔍 Search providers: semantic_scholar, arxiv")

	// Initialize analyzer with aggregator
	c := classifier.NewClassifier()
	a := analyzer.NewAnalyzer(c, aggregator)

	// Create API handler
	handler := api.NewHandler(a)

	// Initialize database if DATABASE_URL is provided
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		pgDB, err := db.NewPostgresDB(dbURL)
		if err != nil {
			log.Printf("⚠️  Failed to connect to PostgreSQL: %v", err)
			log.Printf("📝 Using in-memory storage")
		} else {
			repo := db.NewRepository(pgDB)
			handler.WithRepository(repo)
			log.Printf("✅ PostgreSQL connected")
		}
	} else {
		log.Printf("📝 DATABASE_URL not set, using mock database")
	}

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
