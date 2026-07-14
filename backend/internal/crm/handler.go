package crm

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func parseListParams(r *http.Request) ListParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	return ListParams{
		Page: page, PageSize: pageSize,
		SortBy: r.URL.Query().Get("sort_by"), SortDir: r.URL.Query().Get("sort_dir"),
		Search: r.URL.Query().Get("search"), Status: r.URL.Query().Get("status"),
	}
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, code, map[string]string{"error": msg})
}

// ─── Leads ───────────────────────────────────────────────────

func (h *Handler) ListLeads(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListLeads(r.Context(), parseListParams(r))
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, resp)
}

func (h *Handler) CreateLead(w http.ResponseWriter, r *http.Request) {
	var l Lead
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil { respondError(w, 400, "Invalid body"); return }
	if err := h.svc.CreateLead(r.Context(), &l); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 201, l)
}

func (h *Handler) GetLead(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil { respondError(w, 400, "Invalid ID"); return }
	lead, err := h.svc.GetLead(r.Context(), id)
	if err != nil { respondError(w, 404, "Lead not found"); return }
	respondJSON(w, 200, lead)
}

func (h *Handler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil { respondError(w, 400, "Invalid ID"); return }
	var l Lead
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil { respondError(w, 400, "Invalid body"); return }
	l.ID = id
	if err := h.svc.UpdateLead(r.Context(), &l); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, l)
}

func (h *Handler) DeleteLead(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteLead(r.Context(), id); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, map[string]string{"message": "deleted"})
}

func (h *Handler) ConvertLead(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	deal, err := h.svc.ConvertLead(r.Context(), id)
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, deal)
}

// ─── Contacts ────────────────────────────────────────────────

func (h *Handler) ListContacts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListContacts(r.Context(), parseListParams(r))
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, resp)
}

func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil { respondError(w, 400, "Invalid body"); return }
	if err := h.svc.CreateContact(r.Context(), &c); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 201, c)
}

func (h *Handler) GetContact(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	c, err := h.svc.GetContact(r.Context(), id)
	if err != nil { respondError(w, 404, "Not found"); return }
	respondJSON(w, 200, c)
}

func (h *Handler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var c Contact
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = id
	if err := h.svc.UpdateContact(r.Context(), &c); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, c)
}

func (h *Handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.svc.DeleteContact(r.Context(), id)
	respondJSON(w, 200, map[string]string{"message": "deleted"})
}

// ─── Companies ───────────────────────────────────────────────

func (h *Handler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListCompanies(r.Context(), parseListParams(r))
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, resp)
}

func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var c Company
	json.NewDecoder(r.Body).Decode(&c)
	if err := h.svc.CreateCompany(r.Context(), &c); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 201, c)
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	c, err := h.svc.GetCompany(r.Context(), id)
	if err != nil { respondError(w, 404, "Not found"); return }
	respondJSON(w, 200, c)
}

func (h *Handler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var c Company
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = id
	if err := h.svc.UpdateCompany(r.Context(), &c); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, c)
}

func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.svc.DeleteCompany(r.Context(), id)
	respondJSON(w, 200, map[string]string{"message": "deleted"})
}

// ─── Deals ───────────────────────────────────────────────────

func (h *Handler) ListDeals(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListDeals(r.Context(), parseListParams(r))
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, resp)
}

func (h *Handler) CreateDeal(w http.ResponseWriter, r *http.Request) {
	var d Deal
	json.NewDecoder(r.Body).Decode(&d)
	if err := h.svc.CreateDeal(r.Context(), &d); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 201, d)
}

func (h *Handler) GetDeal(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	d, err := h.svc.GetDeal(r.Context(), id)
	if err != nil { respondError(w, 404, "Not found"); return }
	respondJSON(w, 200, d)
}

func (h *Handler) UpdateDeal(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var d Deal
	json.NewDecoder(r.Body).Decode(&d)
	d.ID = id
	if err := h.svc.UpdateDeal(r.Context(), &d); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, d)
}

func (h *Handler) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.svc.DeleteDeal(r.Context(), id)
	respondJSON(w, 200, map[string]string{"message": "deleted"})
}

// ─── Pipelines ───────────────────────────────────────────────

func (h *Handler) ListPipelines(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.ListPipelines(r.Context())
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, p)
}

func (h *Handler) CreatePipeline(w http.ResponseWriter, r *http.Request) {
	var p Pipeline
	json.NewDecoder(r.Body).Decode(&p)
	if err := h.svc.CreatePipeline(r.Context(), &p); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 201, p)
}

func (h *Handler) UpdatePipeline(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]string{"message": "not implemented"})
}

// ─── Activities ──────────────────────────────────────────────

func (h *Handler) ListActivities(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListActivities(r.Context(), parseListParams(r))
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, resp)
}

func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var a Activity
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.svc.CreateActivity(r.Context(), &a); err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 201, a)
}

// ─── Audit Logs ──────────────────────────────────────────────

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.ListAuditLogs(r.Context(), parseListParams(r))
	if err != nil { respondError(w, 500, err.Error()); return }
	respondJSON(w, 200, resp)
}
