package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abdul-hamid-achik/chessdrill/internal/config"
	"github.com/abdul-hamid-achik/chessdrill/internal/mongo"
)

func main() {
	up := flag.Bool("up", false, "Run migrations (create indexes)")
	flag.Parse()

	if !*up {
		fmt.Println("ChessDrill Migration Tool")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  migrate -up    Run migrations (create indexes)")
		fmt.Println("")
		fmt.Println("Environment variables:")
		fmt.Println("  MONGODB_URI       MongoDB connection URI (default: mongodb://localhost:27017)")
		fmt.Println("  MONGODB_DATABASE  Database name (default: chessdrill)")
		os.Exit(0)
	}

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Printf("Connecting to MongoDB at %s...", cfg.MongoDBURI)

	mongoClient, err := mongo.NewClient(ctx, cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Close(context.Background()); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	log.Printf("Running migrations on database '%s'...", cfg.MongoDBDatabase)

	if err := mongoClient.CreateIndexes(ctx); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}
