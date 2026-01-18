package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Home handler
func Home(w http.ResponseWriter, r *http.Request)  {
	// Only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Welcome to the HTTP server\n")
}

// Health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "ok",
		"service": "httpserver",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Users handler (example REST endpoint)
func Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// In real app: save to database
	user["id"] = 3 // Mock ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//writing the json encoding of user to the output value
	json.NewEncoder(w).Encode(user)
}

// Panic handler (for testing recovery)
func Panic(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	panic("intentional panic for testing recovery middleware")
}

// Slow handler (for testing timeout)
func Slow(w http.ResponseWriter, r *http.Request) {
	// Check if context was cancelled
	select {
	case <-r.Context().Done():
		// Context cancelled (timeout or client disconnect)
		return
	case <-time.After(10 * time.Second):
		// Simulate slow operation
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestTimeout)
		fmt.Fprintf(w, "This took 10 seconds\n")
	}
}