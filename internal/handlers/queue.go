package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/dayzer0-hq/be-project-apple-launch-queue-template/internal/queue"
)

type Handler struct{ Store queue.QueueStore }

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/queue/join", h.Join)
	r.Get("/queue/{queueID}/position", h.GetPosition)
	r.Post("/queue/{queueID}/claim", h.Claim)
	r.Post("/queue/advance", h.Advance)
	return r
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad request", 400); return }
	// TODO: validate req.ProductID and req.UserID non-empty
	entry, err := h.Store.Join(r.Context(), req.ProductID, req.UserID)
	if err != nil { http.Error(w, err.Error(), 500); return }
	writeJSON(w, 201, entry)
}

func (h *Handler) GetPosition(w http.ResponseWriter, r *http.Request) {
	queueID := chi.URLParam(r, "queueID")
	entry, err := h.Store.Position(r.Context(), queueID)
	if errors.Is(err, queue.ErrNotFound) { http.Error(w, "not found", 404); return }
	if err != nil { http.Error(w, err.Error(), 500); return }
	writeJSON(w, 200, entry)
}

func (h *Handler) Claim(w http.ResponseWriter, r *http.Request) {
	queueID := chi.URLParam(r, "queueID")
	token, expiresAt, err := h.Store.Claim(r.Context(), queueID)
	if errors.Is(err, queue.ErrNotReady) { writeJSON(w, 409, map[string]string{"error": "not_ready"}); return }
	if errors.Is(err, queue.ErrNotFound) { http.Error(w, "not found", 404); return }
	if err != nil { http.Error(w, err.Error(), 500); return }
	writeJSON(w, 200, map[string]any{"purchase_token": token, "expires_at": expiresAt})
}

func (h *Handler) Advance(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	batchSize, _ := strconv.Atoi(r.URL.Query().Get("batch_size"))
	if batchSize <= 0 { batchSize = 50 }
	n, err := h.Store.Advance(r.Context(), productID, batchSize)
	if err != nil { http.Error(w, err.Error(), 500); return }
	writeJSON(w, 200, map[string]int{"advanced": n})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
