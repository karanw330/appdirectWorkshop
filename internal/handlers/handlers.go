package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"appdirect-workshop/internal/firestore"

	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

type Handlers struct {
	fsClient        *firestore.Client
	subcollectionID string
	adminPassword   string
}

func NewHandlers(fsClient *firestore.Client, subcollectionID string) *Handlers {
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123" // Default for development only
	}

	return &Handlers{
		fsClient:        fsClient,
		subcollectionID: subcollectionID,
		adminPassword:   adminPassword,
	}
}

// Response helpers
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Attendee handlers
func (h *Handlers) GetAttendees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collection := h.fsClient.GetCollection(ctx, "attendees")

	var attendees []map[string]interface{}
	iter := collection.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		attendees = append(attendees, data)
	}

	respondJSON(w, http.StatusOK, attendees)
}

func (h *Handlers) RegisterAttendee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var attendee map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&attendee); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Add timestamp
	attendee["createdAt"] = time.Now()

	collection := h.fsClient.GetCollection(ctx, "attendees")
	docRef, _, err := collection.Add(ctx, attendee)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	attendee["id"] = docRef.ID
	respondJSON(w, http.StatusCreated, attendee)
}

func (h *Handlers) GetAttendeeCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collection := h.fsClient.GetCollection(ctx, "attendees")

	docs, err := collection.Documents(ctx).GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"count": len(docs)})
}

// Speaker handlers
func (h *Handlers) GetSpeakers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collection := h.fsClient.GetCollection(ctx, "speakers")

	var speakers []map[string]interface{}
	iter := collection.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		speakers = append(speakers, data)
	}

	respondJSON(w, http.StatusOK, speakers)
}

func (h *Handlers) CreateSpeaker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var speaker map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&speaker); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	collection := h.fsClient.GetCollection(ctx, "speakers")
	docRef, _, err := collection.Add(ctx, speaker)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	speaker["id"] = docRef.ID
	respondJSON(w, http.StatusCreated, speaker)
}

func (h *Handlers) UpdateSpeaker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var speaker map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&speaker); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	docRef := h.fsClient.GetCollection(ctx, "speakers").Doc(id)
	_, err := docRef.Set(ctx, speaker)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	speaker["id"] = id
	respondJSON(w, http.StatusOK, speaker)
}

func (h *Handlers) DeleteSpeaker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	docRef := h.fsClient.GetCollection(ctx, "speakers").Doc(id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Speaker deleted"})
}

// Session handlers
func (h *Handlers) GetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	collection := h.fsClient.GetCollection(ctx, "sessions")

	var sessions []map[string]interface{}
	iter := collection.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		sessions = append(sessions, data)
	}

	respondJSON(w, http.StatusOK, sessions)
}

func (h *Handlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var session map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	collection := h.fsClient.GetCollection(ctx, "sessions")
	docRef, _, err := collection.Add(ctx, session)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session["id"] = docRef.ID
	respondJSON(w, http.StatusCreated, session)
}

func (h *Handlers) UpdateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var session map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	docRef := h.fsClient.GetCollection(ctx, "sessions").Doc(id)
	_, err := docRef.Set(ctx, session)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session["id"] = id
	respondJSON(w, http.StatusOK, session)
}

func (h *Handlers) DeleteSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	docRef := h.fsClient.GetCollection(ctx, "sessions").Doc(id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Session deleted"})
}

// Admin handler
func (h *Handlers) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Password != h.adminPassword {
		respondError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

