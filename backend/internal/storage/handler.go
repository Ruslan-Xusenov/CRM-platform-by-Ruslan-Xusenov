package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Filename string `json:"filename"`
		Bucket   string `json:"bucket"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Bucket == "" { req.Bucket = "crm-files" }
	objectName := fmt.Sprintf("%s/%s_%s", time.Now().Format("2006/01/02"), uuid.NewString()[:8], req.Filename)

	url, err := h.svc.GetPresignedUploadURL(req.Bucket, objectName, 15*time.Minute)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"upload_url": url, "object_name": objectName})
}

func (h *Handler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	objectName := chi.URLParam(r, "id")
	bucket := r.URL.Query().Get("bucket")
	if bucket == "" { bucket = "crm-files" }

	url, err := h.svc.GetPresignedDownloadURL(bucket, objectName, 1*time.Hour)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"download_url": url})
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	objectName := chi.URLParam(r, "id")
	bucket := r.URL.Query().Get("bucket")
	if bucket == "" { bucket = "crm-files" }

	if err := h.svc.DeleteObject(bucket, objectName); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}
