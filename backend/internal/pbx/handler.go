package pbx

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crm-platform/backend/internal/shared/middleware"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) ListActiveCalls(w http.ResponseWriter, r *http.Request) {
	calls := h.svc.GetActiveCalls()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calls)
}

func (h *Handler) ListCallHistory(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 50 }
	records, err := h.svc.ListCallHistory(r.Context(), page, pageSize)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *Handler) OriginateCall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		To string `json:"to"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	tid, _ := r.Context().Value(middleware.TenantIDKey).(string)
	uid, _ := r.Context().Value(middleware.UserIDKey).(string)
	if err := h.svc.OriginateCall(r.Context(), uid, req.To, tid); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "call originated"})
}

func (h *Handler) TransferCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "transfer not yet implemented"})
}

func (h *Handler) HoldCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "hold not yet implemented"})
}

func (h *Handler) HangupCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "hangup not yet implemented"})
}

func (h *Handler) ListExtensions(w http.ResponseWriter, r *http.Request) {
	exts, err := h.svc.ListExtensions(r.Context())
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exts)
}

func (h *Handler) CreateExtension(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "create extension not yet implemented"})
}

func (h *Handler) UpdateExtension(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "update extension not yet implemented"})
}

func (h *Handler) DeleteExtension(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "delete extension not yet implemented"})
}

func (h *Handler) ListRoutingRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]struct{}{})
}

func (h *Handler) UpdateRoutingRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "update routing not yet implemented"})
}
