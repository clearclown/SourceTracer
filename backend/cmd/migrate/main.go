package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Get command
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down|version|force]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	defer m.Close() //nolint:errcheck

	// Execute command
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✅ Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✅ Migrations rolled back successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current version: %d (dirty: %v)\n", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate force [version]")
		}
		version := os.Args[2]
		var v int
		fmt.Sscanf(version, "%d", &v) //nolint:errcheck
		if err := m.Force(v); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Printf("✅ Forced version to %d\n", v)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
