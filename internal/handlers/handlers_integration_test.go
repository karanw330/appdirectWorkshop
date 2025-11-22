// +build integration

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"appdirect-workshop/internal/firestore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests require a real Firestore connection
// These should be run separately with: go test -tags=integration

func setupIntegrationTest(t *testing.T) (*Handlers, func()) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		t.Skip("FIREBASE_PROJECT_ID not set, skipping integration test")
	}

	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if serviceAccountPath == "" {
		serviceAccountPath = "./serviceAccountKey.json"
	}

	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, projectID, serviceAccountPath)
	require.NoError(t, err)

	handler := NewHandlers(fsClient, "test_collection")

	cleanup := func() {
		fsClient.Close()
	}

	return handler, cleanup
}

func TestIntegrationGetAttendees(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/attendees", nil)
	w := httptest.NewRecorder()

	handler.GetAttendees(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var attendees []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &attendees)
	assert.NoError(t, err)
	assert.IsType(t, []map[string]interface{}{}, attendees)
}

func TestIntegrationRegisterAttendee(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	attendee := map[string]interface{}{
		"name":        "Test User",
		"email":       "test@example.com",
		"designation": "Software Engineer",
	}

	jsonBody, _ := json.Marshal(attendee)
	req := httptest.NewRequest("POST", "/api/attendees", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RegisterAttendee(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", response["name"])
	assert.NotEmpty(t, response["id"])
}

func TestIntegrationGetAttendeeCount(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/attendees/count", nil)
	w := httptest.NewRecorder()

	handler.GetAttendeeCount(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "count")
	assert.GreaterOrEqual(t, response["count"], 0)
}
