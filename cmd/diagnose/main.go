package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"google.golang.org/api/iterator"

	"appdirect-workshop/internal/firestore"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	log.Println("--- Diagnostic Start ---")
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	} else {
		log.Println("✓ Loaded .env file")
	}

	// 2. Check Env Vars
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")

	fmt.Printf("FIREBASE_PROJECT_ID: %s\n", projectID)
	fmt.Printf("FIREBASE_SERVICE_ACCOUNT_PATH: %s\n", serviceAccountPath)

	// 3. Resolve Service Account Path
	if serviceAccountPath == "" {
		serviceAccountPath = serviceAccountPath
		fmt.Printf("Defaulting serviceAccountPath to: %s\n", serviceAccountPath)
	}

	absPath, _ := filepath.Abs(serviceAccountPath)
	fmt.Printf("Absolute path: %s\n", absPath)

	if _, err := os.Stat(serviceAccountPath); err != nil {
		fmt.Printf("❌ FILE NOT FOUND: %s\n", serviceAccountPath)
		// Try to find if there is a similar file
		files, _ := filepath.Glob("service*.json")
		if len(files) > 0 {
			fmt.Printf("Did you mean one of these? %v\n", files)
		}
	} else {
		fmt.Printf("✓ File exists: %s\n", serviceAccountPath)
	}

	// 4. Connect to Firestore
	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	ctx := context.Background()
	log.Println("Attempting to initialize Firestore client...")
	fsClient, err := firestore.NewClient(ctx, projectID, databaseID, serviceAccountPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Firestore: %v", err)
	}
	defer fsClient.Close()
	log.Println("✓ Firestore client initialized")

	// 5. Test Query
	log.Println("Attempting to list 1 document from 'sessions' collection...")
	iter := fsClient.GetCollection(ctx, "sessions").Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		log.Println("✓ Connection successful! Collection is empty or exhausted.")
	} else if err != nil {
		log.Fatalf("❌ Query Failed: %v", err)
	} else {
		log.Printf("✓ Connection successful! Found document ID: %s", doc.Ref.ID)
	}
	log.Println("--- Diagnostic End ---")
}
