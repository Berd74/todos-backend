package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

var database *spanner.Client

func init() {
	// Load .env file
	err := godotenv.Load() // This will look for a .env file in the current directory
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	projectId := os.Getenv("PROJECTID")
	if projectId == "" {
		log.Fatal("PROJECTID must be set.")
	}

	instanceName := os.Getenv("INSTANCE")
	if instanceName == "" {
		log.Fatal("INSTANCE must be set.")
	}

	databaseName := os.Getenv("DATABASE")
	if databaseName == "" {
		log.Fatal("DATABASE must be set.")
	}

	// Creates a Spanner client.
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceName, databaseName), option.WithCredentialsFile("firebase-adminsdk.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	database = client
}

func GetDatabase() *spanner.Client {
	return database
}
