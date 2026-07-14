package crm

import (
	"context"
	"log/slog"

	"github.com/crm-platform/backend/internal/shared/middleware"
	"github.com/crm-platform/backend/pkg/broker"
	"github.com/crm-platform/backend/pkg/ws"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	mq   *broker.RabbitMQ
	hub  *ws.Hub
}

func NewService(repo *Repository, mq *broker.RabbitMQ, hub *ws.Hub) *Service {
	return &Service{repo: repo, mq: mq, hub: hub}
}

// ─── Leads ───────────────────────────────────────────────────

func (s *Service) ListLeads(ctx context.Context, p ListParams) (*ListResponse, error) {
	return s.repo.ListLeads(ctx, p)
}

func (s *Service) GetLead(ctx context.Context, id uuid.UUID) (*Lead, error) {
	return s.repo.GetLead(ctx, id)
}

func (s *Service) CreateLead(ctx context.Context, l *Lead) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	l.TenantID = uuid.MustParse(tid)
	creatorID := uuid.MustParse(uid)
	l.CreatedBy = &creatorID
	if l.Status == "" { l.Status = "new" }
	if l.Currency == "" { l.Currency = "UZS" }

	if err := s.repo.CreateLead(ctx, l); err != nil {
		return err
	}

	// Audit log
	s.repo.CreateAuditLog(ctx, tid, "lead", l.ID, "create", creatorID, l)

	// Publish event
	s.mq.Publish(ctx, "crm.events", "lead.created", broker.Event{Type: "lead.created", Payload: l})

	// Broadcast via WebSocket
	s.hub.BroadcastToTenant(tid, ws.Message{Type: "lead.created", Payload: l})

	slog.Info("Lead created", "id", l.ID, "title", l.Title, "tenant", tid)
	return nil
}

func (s *Service) UpdateLead(ctx context.Context, l *Lead) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	l.TenantID = uuid.MustParse(tid)

	if err := s.repo.UpdateLead(ctx, l); err != nil {
		return err
	}

	s.repo.CreateAuditLog(ctx, tid, "lead", l.ID, "update", uuid.MustParse(uid), l)
	s.mq.Publish(ctx, "crm.events", "lead.updated", broker.Event{Type: "lead.updated", Payload: l})
	s.hub.BroadcastToTenant(tid, ws.Message{Type: "lead.updated", Payload: l})
	return nil
}

func (s *Service) DeleteLead(ctx context.Context, id uuid.UUID) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)

	if err := s.repo.DeleteLead(ctx, id, tid); err != nil {
		return err
	}

	s.repo.CreateAuditLog(ctx, tid, "lead", id, "delete", uuid.MustParse(uid), nil)
	s.hub.BroadcastToTenant(tid, ws.Message{Type: "lead.deleted", Payload: map[string]string{"id": id.String()}})
	return nil
}

func (s *Service) ConvertLead(ctx context.Context, leadID uuid.UUID) (*Deal, error) {
	lead, err := s.repo.GetLead(ctx, leadID)
	if err != nil { return nil, err }

	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	creatorID := uuid.MustParse(uid)

	// Create contact from lead
	contact := &Contact{
		TenantID:  uuid.MustParse(tid),
		FirstName: lead.ContactName,
		LastName:  "",
		Email:     lead.ContactEmail,
		Phone:     lead.ContactPhone,
		Source:    lead.Source,
		CreatedBy: &creatorID,
	}
	s.repo.CreateContact(ctx, contact)

	// Create deal from lead
	deal := &Deal{
		TenantID:  uuid.MustParse(tid),
		Title:     lead.Title,
		Amount:    lead.Budget,
		Currency:  lead.Currency,
		ContactID: &contact.ID,
		CreatedBy: &creatorID,
	}
	s.repo.CreateDeal(ctx, deal)

	// Mark lead as converted (soft update)
	lead.Status = "converted"
	s.repo.UpdateLead(ctx, lead)

	s.repo.CreateAuditLog(ctx, tid, "lead", leadID, "convert", creatorID, map[string]interface{}{"contact_id": contact.ID, "deal_id": deal.ID})
	return deal, nil
}

// ─── Contacts ────────────────────────────────────────────────

func (s *Service) ListContacts(ctx context.Context, p ListParams) (*ListResponse, error) {
	return s.repo.ListContacts(ctx, p)
}

func (s *Service) GetContact(ctx context.Context, id uuid.UUID) (*Contact, error) {
	return s.repo.GetContact(ctx, id)
}

func (s *Service) CreateContact(ctx context.Context, c *Contact) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	c.TenantID = uuid.MustParse(tid)
	creatorID := uuid.MustParse(uid)
	c.CreatedBy = &creatorID
	if err := s.repo.CreateContact(ctx, c); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "contact", c.ID, "create", creatorID, c)
	return nil
}

func (s *Service) UpdateContact(ctx context.Context, c *Contact) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	c.TenantID = uuid.MustParse(tid)
	if err := s.repo.UpdateContact(ctx, c); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "contact", c.ID, "update", uuid.MustParse(uid), c)
	return nil
}

func (s *Service) DeleteContact(ctx context.Context, id uuid.UUID) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	if err := s.repo.DeleteContact(ctx, id, tid); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "contact", id, "delete", uuid.MustParse(uid), nil)
	return nil
}

// ─── Companies ───────────────────────────────────────────────

func (s *Service) ListCompanies(ctx context.Context, p ListParams) (*ListResponse, error) {
	return s.repo.ListCompanies(ctx, p)
}
func (s *Service) GetCompany(ctx context.Context, id uuid.UUID) (*Company, error) {
	return s.repo.GetCompany(ctx, id)
}
func (s *Service) CreateCompany(ctx context.Context, c *Company) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	c.TenantID = uuid.MustParse(tid)
	creatorID := uuid.MustParse(uid)
	c.CreatedBy = &creatorID
	if err := s.repo.CreateCompany(ctx, c); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "company", c.ID, "create", creatorID, c)
	return nil
}
func (s *Service) UpdateCompany(ctx context.Context, c *Company) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	c.TenantID = uuid.MustParse(tid)
	if err := s.repo.UpdateCompany(ctx, c); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "company", c.ID, "update", uuid.MustParse(uid), c)
	return nil
}
func (s *Service) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	if err := s.repo.DeleteCompany(ctx, id, tid); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "company", id, "delete", uuid.MustParse(uid), nil)
	return nil
}

// ─── Deals ───────────────────────────────────────────────────

func (s *Service) ListDeals(ctx context.Context, p ListParams) (*ListResponse, error) {
	return s.repo.ListDeals(ctx, p)
}
func (s *Service) GetDeal(ctx context.Context, id uuid.UUID) (*Deal, error) {
	return s.repo.GetDeal(ctx, id)
}
func (s *Service) CreateDeal(ctx context.Context, d *Deal) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	d.TenantID = uuid.MustParse(tid)
	creatorID := uuid.MustParse(uid)
	d.CreatedBy = &creatorID
	if d.Currency == "" { d.Currency = "UZS" }
	if err := s.repo.CreateDeal(ctx, d); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "deal", d.ID, "create", creatorID, d)
	s.hub.BroadcastToTenant(tid, ws.Message{Type: "deal.created", Payload: d})
	return nil
}
func (s *Service) UpdateDeal(ctx context.Context, d *Deal) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	d.TenantID = uuid.MustParse(tid)
	if err := s.repo.UpdateDeal(ctx, d); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "deal", d.ID, "update", uuid.MustParse(uid), d)
	s.hub.BroadcastToTenant(tid, ws.Message{Type: "deal.updated", Payload: d})
	return nil
}
func (s *Service) DeleteDeal(ctx context.Context, id uuid.UUID) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	if err := s.repo.DeleteDeal(ctx, id, tid); err != nil { return err }
	s.repo.CreateAuditLog(ctx, tid, "deal", id, "delete", uuid.MustParse(uid), nil)
	return nil
}

// ─── Pipelines ───────────────────────────────────────────────

func (s *Service) ListPipelines(ctx context.Context) ([]Pipeline, error) {
	return s.repo.ListPipelines(ctx)
}
func (s *Service) CreatePipeline(ctx context.Context, p *Pipeline) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	p.TenantID = uuid.MustParse(tid)
	return s.repo.CreatePipeline(ctx, p)
}
func (s *Service) UpdatePipeline(ctx context.Context, p *Pipeline) error {
	return nil // TODO: implement pipeline update
}

// ─── Activities ──────────────────────────────────────────────

func (s *Service) ListActivities(ctx context.Context, p ListParams) (*ListResponse, error) {
	return s.repo.ListActivities(ctx, p)
}
func (s *Service) CreateActivity(ctx context.Context, a *Activity) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	uid, _ := ctx.Value(middleware.UserIDKey).(string)
	a.TenantID = uuid.MustParse(tid)
	creatorID := uuid.MustParse(uid)
	a.CreatedBy = &creatorID
	return s.repo.CreateActivity(ctx, a)
}

// ─── Audit Logs ──────────────────────────────────────────────

func (s *Service) ListAuditLogs(ctx context.Context, p ListParams) (*ListResponse, error) {
	return s.repo.ListAuditLogs(ctx, p)
}
