package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminLogin(t *testing.T) {
	os.Setenv("ADMIN_PASSWORD", "testpass")
	defer os.Unsetenv("ADMIN_PASSWORD")

	handler := &Handlers{
		adminPassword: "testpass",
	}

	tests := []struct {
		name           string
		password       string
		expectedStatus int
	}{
		{
			name:           "Valid password",
			password:       "testpass",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid password",
			password:       "wrongpass",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty password",
			password:       "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]string{"password": tt.password}
			jsonBody, _ := json.Marshal(body)
			req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.AdminLogin(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "Login successful", response["message"])
			} else {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "Invalid password", response["error"])
			}
		})
	}
}

func TestAdminLoginInvalidJSON(t *testing.T) {
	handler := &Handlers{
		adminPassword: "testpass",
	}

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AdminLogin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestNewHandlers(t *testing.T) {
	os.Setenv("ADMIN_PASSWORD", "envpassword")
	defer os.Unsetenv("ADMIN_PASSWORD")

	// This test requires a real firestore client, so we'll just test the password logic
	// In a real scenario, you'd inject a mock
	handler := NewHandlers(nil, "test_collection")
	assert.Equal(t, "envpassword", handler.adminPassword)
	assert.Equal(t, "test_collection", handler.subcollectionID)
}

func TestNewHandlersDefaultPassword(t *testing.T) {
	os.Unsetenv("ADMIN_PASSWORD")
	
	handler := NewHandlers(nil, "test_collection")
	assert.Equal(t, "admin123", handler.adminPassword)
}
