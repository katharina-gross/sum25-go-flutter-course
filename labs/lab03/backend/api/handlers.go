package api

import (
	"encoding/json"
	"lab03-backend/storage"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// APIResponse standard response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HTTPStatusResponse structure for GetHTTPStatus
type HTTPStatusResponse struct {
	StatusCode  int    `json:"status_code"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

// HealthResponse structure for HealthCheck
type HealthResponse struct {
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
	TotalMessages int       `json:"total_messages"`
}

// CreateMessageRequest structure for CreateMessage
type CreateMessageRequest struct {
	Content string `json:"content"`
}

// UpdateMessageRequest structure for UpdateMessage
type UpdateMessageRequest struct {
	Content string `json:"content"`
}

// Handler holds the storage instance
type Handler struct {
	storage *storage.MemoryStorage
}

// NewHandler creates a new handler instance
func NewHandler(storage *storage.MemoryStorage) *Handler {
	return &Handler{storage: storage}
}

// SetupRoutes configures all API routes
func (h *Handler) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/messages", h.GetMessages).Methods("GET")
	apiRouter.HandleFunc("/messages", h.CreateMessage).Methods("POST")
	apiRouter.HandleFunc("/messages/{id}", h.UpdateMessage).Methods("PUT")
	apiRouter.HandleFunc("/messages/{id}", h.DeleteMessage).Methods("DELETE")
	apiRouter.HandleFunc("/status/{code}", h.GetHTTPStatus).Methods("GET")
	apiRouter.HandleFunc("/health", h.HealthCheck).Methods("GET")

	return router
}

// GetMessages handles GET /api/messages
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.storage.GetAll()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to get messages")
		return
	}
	h.writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    messages,
	})
}

// CreateMessage handles POST /api/messages
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req CreateMessageRequest
	if err := h.parseJSON(r, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "content cannot be empty")
		return
	}
	message, err := h.storage.Create(req.Content)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to create message")
		return
	}
	h.writeJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    message,
	})
}

// UpdateMessage handles PUT /api/messages/{id}
func (h *Handler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid message ID")
		return
	}
	var req UpdateMessageRequest
	if err := h.parseJSON(r, &req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "content cannot be empty")
		return
	}
	message, err := h.storage.Update(id, req.Content)
	if err != nil {
		if err == storage.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "message not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to update message")
		return
	}
	h.writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    message,
	})
}

// DeleteMessage handles DELETE /api/messages/{id}
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid message ID")
		return
	}
	if err := h.storage.Delete(id); err != nil {
		if err == storage.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "message not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to delete message")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetHTTPStatus handles GET /api/status/{code}
func (h *Handler) GetHTTPStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, err := strconv.Atoi(vars["code"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid status code")
		return
	}
	if code < 100 || code > 599 {
		h.writeError(w, http.StatusBadRequest, "status code must be between 100-599")
		return
	}
	response := HTTPStatusResponse{
		StatusCode:  code,
		ImageURL:    "https://http.cat/" + strconv.Itoa(code),
		Description: getHTTPStatusDescription(code),
	}
	h.writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// HealthCheck handles GET /api/health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	count, err := h.storage.Count()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to get message count")
		return
	}
	response := HealthResponse{
		Status:        "ok",
		Message:       "API is running",
		Timestamp:     time.Now().UTC(),
		TotalMessages: count,
	}
	h.writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// Helper function to write JSON responses
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

// Helper function to write error responses
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}

// Helper function to parse JSON request body
func (h *Handler) parseJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

// Helper function to get HTTP status description
func getHTTPStatusDescription(code int) string {
	switch code {
	case http.StatusOK:
		return "OK"
	case http.StatusCreated:
		return "Created"
	case http.StatusNoContent:
		return "No Content"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
