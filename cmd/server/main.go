package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"appdirect-workshop/internal/handlers"
	"appdirect-workshop/internal/firestore"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables from .env file
	// Get the current working directory
	wd, _ := os.Getwd()
	
	// Try to find .env file - start from project root (two levels up from cmd/server/)
	envPaths := []string{
		".env",                              // Current directory
		filepath.Join(wd, ".env"),           // Absolute from working dir
		filepath.Join(wd, "..", ".env"),     // One level up
		filepath.Join(wd, "../..", ".env"),  // Two levels up (project root)
		"../.env",                           // Relative one level up
		"../../.env",                        // Relative two levels up
	}
	
	// Try relative to source file location
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		sourceDir := filepath.Dir(filename)
		projectRoot := filepath.Join(sourceDir, "../..")
		envPaths = append(envPaths, filepath.Join(projectRoot, ".env"))
	}
	
	var envLoaded bool
	for _, path := range envPaths {
		// Check if file exists first
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				log.Printf("✓ Loaded .env file from: %s", path)
				envLoaded = true
				break
			} else {
				log.Printf("Failed to load .env from %s: %v", path, err)
			}
		}
	}
	if !envLoaded {
		log.Printf("⚠ No .env file found. Working directory: %s", wd)
		log.Println("Tried paths:", envPaths)
		log.Println("Please ensure .env file exists in the appdirectWorkshop directory")
	}

	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	
	// For Cloud Run, use Application Default Credentials (ADC)
	// Set FIREBASE_SERVICE_ACCOUNT_PATH=ADC or leave empty in Cloud Run
	// For local development, use service account file path
	if serviceAccountPath == "" {
		// Check if running in Cloud Run (has K_SERVICE env var)
		if os.Getenv("K_SERVICE") != "" {
			// Running in Cloud Run, use ADC
			serviceAccountPath = "ADC"
		} else {
			// Local development, try to use service account file
			serviceAccountPath = serviceAccountPath
		}
	}

	subcollectionID := os.Getenv("FIRESTORE_SUBCOLLECTION_ID")
	if subcollectionID == "" {
		subcollectionID = "workshop_attendees"
	}

	// Initialize Firestore client
	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	log.Printf("DEBUG: FIRESTORE_DATABASE_ID is: '%s'", databaseID)
	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, projectID, databaseID, serviceAccountPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer fsClient.Close()

	// Initialize handlers
	h := handlers.NewHandlers(fsClient, subcollectionID)

	// Setup router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// API info endpoint
	api.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "AppDirect Workshop API",
			"version": "1.0.0",
			"endpoints": map[string]interface{}{
				"attendees": map[string]string{
					"GET":  "/api/attendees",
					"POST": "/api/attendees",
					"GET_count": "/api/attendees/count",
				},
				"speakers": map[string]string{
					"GET":    "/api/speakers",
					"POST":   "/api/speakers",
					"PUT":    "/api/speakers/{id}",
					"DELETE": "/api/speakers/{id}",
				},
				"sessions": map[string]string{
					"GET":    "/api/sessions",
					"POST":   "/api/sessions",
					"PUT":    "/api/sessions/{id}",
					"DELETE": "/api/sessions/{id}",
				},
				"admin": map[string]string{
					"POST": "/api/admin/login",
				},
			},
		})
	}).Methods("GET")

	// Attendees
	api.HandleFunc("/attendees", h.GetAttendees).Methods("GET")
	api.HandleFunc("/attendees", h.RegisterAttendee).Methods("POST")
	api.HandleFunc("/attendees/count", h.GetAttendeeCount).Methods("GET")

	// Speakers
	api.HandleFunc("/speakers", h.GetSpeakers).Methods("GET")
	api.HandleFunc("/speakers", h.CreateSpeaker).Methods("POST")
	api.HandleFunc("/speakers/{id}", h.UpdateSpeaker).Methods("PUT")
	api.HandleFunc("/speakers/{id}", h.DeleteSpeaker).Methods("DELETE")

	// Sessions
	api.HandleFunc("/sessions", h.GetSessions).Methods("GET")
	api.HandleFunc("/sessions", h.CreateSession).Methods("POST")
	api.HandleFunc("/sessions/{id}", h.UpdateSession).Methods("PUT")
	api.HandleFunc("/sessions/{id}", h.DeleteSession).Methods("DELETE")

	// Admin
	api.HandleFunc("/admin/login", h.AdminLogin).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Serve static files in production (if static directory exists)
	var handler http.Handler = c.Handler(r)
	if _, err := os.Stat("./static"); err == nil {
		fs := http.FileServer(http.Dir("./static"))
		handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Don't serve static files for API routes or health check
			if len(req.URL.Path) >= 4 && req.URL.Path[:4] == "/api" {
				c.Handler(r).ServeHTTP(w, req)
				return
			}
			if req.URL.Path == "/health" {
				c.Handler(r).ServeHTTP(w, req)
				return
			}
			// For React Router, serve index.html for all non-API routes if file doesn't exist
			if req.URL.Path != "/" {
				if _, err := os.Stat("./static" + req.URL.Path); os.IsNotExist(err) {
					req.URL.Path = "/"
				}
			}
			fs.ServeHTTP(w, req)
		})
	} else {
		handler = c.Handler(r)
	}

	// Start server
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

