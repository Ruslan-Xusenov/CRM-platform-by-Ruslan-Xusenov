package crm

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"github.com/crm-platform/backend/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) setTenant(ctx context.Context) error {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	if tid == "" { return fmt.Errorf("no tenant") }
	_, err := r.db.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tid))
	return err
}

// ─── Leads ───────────────────────────────────────────────────

func (r *Repository) ListLeads(ctx context.Context, p ListParams) (*ListResponse, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	offset := (p.Page - 1) * p.PageSize
	if p.PageSize == 0 { p.PageSize = 20 }
	if p.Page == 0 { p.Page = 1 }

	var total int64
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM leads WHERE tenant_id=$1 AND deleted_at IS NULL", tid).Scan(&total)

	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, title, status, source, budget, currency, contact_name, contact_phone, 
		 contact_email, company_name, description, custom_fields, assigned_to, created_by, created_at, updated_at
		 FROM leads WHERE tenant_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tid, p.PageSize, offset)
	if err != nil { return nil, err }
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var l Lead
		var cf []byte
		err := rows.Scan(&l.ID, &l.TenantID, &l.Title, &l.Status, &l.Source, &l.Budget, &l.Currency,
			&l.ContactName, &l.ContactPhone, &l.ContactEmail, &l.CompanyName, &l.Description,
			&cf, &l.AssignedTo, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt)
		if err != nil { return nil, err }
		json.Unmarshal(cf, &l.CustomFields)
		leads = append(leads, l)
	}

	return &ListResponse{Data: leads, Total: total, Page: p.Page, PageSize: p.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PageSize)))}, nil
}

func (r *Repository) GetLead(ctx context.Context, id uuid.UUID) (*Lead, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	var l Lead; var cf []byte
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, title, status, source, budget, currency, contact_name, contact_phone,
		 contact_email, company_name, description, custom_fields, assigned_to, created_by, created_at, updated_at
		 FROM leads WHERE id=$1 AND tenant_id=$2 AND deleted_at IS NULL`, id, tid).
		Scan(&l.ID, &l.TenantID, &l.Title, &l.Status, &l.Source, &l.Budget, &l.Currency,
			&l.ContactName, &l.ContactPhone, &l.ContactEmail, &l.CompanyName, &l.Description,
			&cf, &l.AssignedTo, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt)
	if err != nil { return nil, err }
	json.Unmarshal(cf, &l.CustomFields)
	return &l, nil
}

func (r *Repository) CreateLead(ctx context.Context, l *Lead) error {
	cf, _ := json.Marshal(l.CustomFields)
	return r.db.QueryRow(ctx,
		`INSERT INTO leads (tenant_id, title, status, source, budget, currency, contact_name, contact_phone,
		 contact_email, company_name, description, custom_fields, assigned_to, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id, created_at`,
		l.TenantID, l.Title, l.Status, l.Source, l.Budget, l.Currency, l.ContactName, l.ContactPhone,
		l.ContactEmail, l.CompanyName, l.Description, cf, l.AssignedTo, l.CreatedBy).Scan(&l.ID, &l.CreatedAt)
}

func (r *Repository) UpdateLead(ctx context.Context, l *Lead) error {
	cf, _ := json.Marshal(l.CustomFields)
	_, err := r.db.Exec(ctx,
		`UPDATE leads SET title=$1, status=$2, source=$3, budget=$4, contact_name=$5, contact_phone=$6,
		 contact_email=$7, company_name=$8, description=$9, custom_fields=$10, assigned_to=$11, updated_at=NOW()
		 WHERE id=$12 AND tenant_id=$13 AND deleted_at IS NULL`,
		l.Title, l.Status, l.Source, l.Budget, l.ContactName, l.ContactPhone,
		l.ContactEmail, l.CompanyName, l.Description, cf, l.AssignedTo, l.ID, l.TenantID)
	return err
}

func (r *Repository) DeleteLead(ctx context.Context, id uuid.UUID, tenantID string) error {
	_, err := r.db.Exec(ctx, `UPDATE leads SET deleted_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

// ─── Contacts ────────────────────────────────────────────────

func (r *Repository) ListContacts(ctx context.Context, p ListParams) (*ListResponse, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	offset := (p.Page - 1) * p.PageSize
	if p.PageSize == 0 { p.PageSize = 20 }
	var total int64
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM contacts WHERE tenant_id=$1 AND deleted_at IS NULL", tid).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, company_id, first_name, last_name, email, phone, mobile, position, source,
		 custom_fields, assigned_to, created_by, created_at, updated_at
		 FROM contacts WHERE tenant_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tid, p.PageSize, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var contacts []Contact
	for rows.Next() {
		var c Contact; var cf []byte
		rows.Scan(&c.ID, &c.TenantID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone,
			&c.Mobile, &c.Position, &c.Source, &cf, &c.AssignedTo, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
		json.Unmarshal(cf, &c.CustomFields)
		contacts = append(contacts, c)
	}
	return &ListResponse{Data: contacts, Total: total, Page: p.Page, PageSize: p.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PageSize)))}, nil
}

func (r *Repository) CreateContact(ctx context.Context, c *Contact) error {
	cf, _ := json.Marshal(c.CustomFields)
	return r.db.QueryRow(ctx,
		`INSERT INTO contacts (tenant_id, company_id, first_name, last_name, email, phone, mobile, position, source, custom_fields, assigned_to, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id, created_at`,
		c.TenantID, c.CompanyID, c.FirstName, c.LastName, c.Email, c.Phone, c.Mobile, c.Position, c.Source, cf, c.AssignedTo, c.CreatedBy).Scan(&c.ID, &c.CreatedAt)
}

func (r *Repository) GetContact(ctx context.Context, id uuid.UUID) (*Contact, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	var c Contact; var cf []byte
	err := r.db.QueryRow(ctx,
		`SELECT id, tenant_id, company_id, first_name, last_name, email, phone, mobile, position, source, custom_fields, assigned_to, created_by, created_at, updated_at
		 FROM contacts WHERE id=$1 AND tenant_id=$2 AND deleted_at IS NULL`, id, tid).
		Scan(&c.ID, &c.TenantID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Mobile, &c.Position, &c.Source, &cf, &c.AssignedTo, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil { return nil, err }
	json.Unmarshal(cf, &c.CustomFields)
	return &c, nil
}

func (r *Repository) UpdateContact(ctx context.Context, c *Contact) error {
	cf, _ := json.Marshal(c.CustomFields)
	_, err := r.db.Exec(ctx,
		`UPDATE contacts SET first_name=$1, last_name=$2, email=$3, phone=$4, mobile=$5, position=$6, company_id=$7, custom_fields=$8, assigned_to=$9, updated_at=NOW()
		 WHERE id=$10 AND tenant_id=$11 AND deleted_at IS NULL`,
		c.FirstName, c.LastName, c.Email, c.Phone, c.Mobile, c.Position, c.CompanyID, cf, c.AssignedTo, c.ID, c.TenantID)
	return err
}

func (r *Repository) DeleteContact(ctx context.Context, id uuid.UUID, tid string) error {
	_, err := r.db.Exec(ctx, `UPDATE contacts SET deleted_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tid)
	return err
}

// ─── Companies (similar pattern) ─────────────────────────────

func (r *Repository) ListCompanies(ctx context.Context, p ListParams) (*ListResponse, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	if p.PageSize == 0 { p.PageSize = 20 }
	offset := (p.Page - 1) * p.PageSize
	var total int64
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM companies WHERE tenant_id=$1 AND deleted_at IS NULL", tid).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, name, industry, website, phone, email, address, custom_fields, created_at
		 FROM companies WHERE tenant_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tid, p.PageSize, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var companies []Company
	for rows.Next() {
		var c Company; var cf []byte
		rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Industry, &c.Website, &c.Phone, &c.Email, &c.Address, &cf, &c.CreatedAt)
		json.Unmarshal(cf, &c.CustomFields)
		companies = append(companies, c)
	}
	return &ListResponse{Data: companies, Total: total, Page: p.Page, PageSize: p.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PageSize)))}, nil
}

func (r *Repository) CreateCompany(ctx context.Context, c *Company) error {
	cf, _ := json.Marshal(c.CustomFields)
	return r.db.QueryRow(ctx, `INSERT INTO companies (tenant_id, name, industry, website, phone, email, address, custom_fields, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_at`,
		c.TenantID, c.Name, c.Industry, c.Website, c.Phone, c.Email, c.Address, cf, c.CreatedBy).Scan(&c.ID, &c.CreatedAt)
}

func (r *Repository) GetCompany(ctx context.Context, id uuid.UUID) (*Company, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	var c Company; var cf []byte
	err := r.db.QueryRow(ctx, `SELECT id, tenant_id, name, industry, website, phone, email, address, custom_fields, created_at FROM companies WHERE id=$1 AND tenant_id=$2 AND deleted_at IS NULL`, id, tid).
		Scan(&c.ID, &c.TenantID, &c.Name, &c.Industry, &c.Website, &c.Phone, &c.Email, &c.Address, &cf, &c.CreatedAt)
	if err != nil { return nil, err }
	json.Unmarshal(cf, &c.CustomFields)
	return &c, nil
}

func (r *Repository) UpdateCompany(ctx context.Context, c *Company) error {
	cf, _ := json.Marshal(c.CustomFields)
	_, err := r.db.Exec(ctx, `UPDATE companies SET name=$1, industry=$2, website=$3, phone=$4, email=$5, address=$6, custom_fields=$7, updated_at=NOW() WHERE id=$8 AND tenant_id=$9 AND deleted_at IS NULL`,
		c.Name, c.Industry, c.Website, c.Phone, c.Email, c.Address, cf, c.ID, c.TenantID)
	return err
}

func (r *Repository) DeleteCompany(ctx context.Context, id uuid.UUID, tid string) error {
	_, err := r.db.Exec(ctx, `UPDATE companies SET deleted_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tid)
	return err
}

// ─── Deals ───────────────────────────────────────────────────

func (r *Repository) ListDeals(ctx context.Context, p ListParams) (*ListResponse, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	if p.PageSize == 0 { p.PageSize = 20 }
	offset := (p.Page - 1) * p.PageSize
	var total int64
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM deals WHERE tenant_id=$1 AND deleted_at IS NULL", tid).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, title, pipeline_id, stage_id, contact_id, company_id, amount, currency, probability, custom_fields, assigned_to, created_at
		 FROM deals WHERE tenant_id=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tid, p.PageSize, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var deals []Deal
	for rows.Next() {
		var d Deal; var cf []byte
		rows.Scan(&d.ID, &d.TenantID, &d.Title, &d.PipelineID, &d.StageID, &d.ContactID, &d.CompanyID, &d.Amount, &d.Currency, &d.Probability, &cf, &d.AssignedTo, &d.CreatedAt)
		json.Unmarshal(cf, &d.CustomFields)
		deals = append(deals, d)
	}
	return &ListResponse{Data: deals, Total: total, Page: p.Page, PageSize: p.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PageSize)))}, nil
}

func (r *Repository) CreateDeal(ctx context.Context, d *Deal) error {
	cf, _ := json.Marshal(d.CustomFields)
	return r.db.QueryRow(ctx, `INSERT INTO deals (tenant_id, title, pipeline_id, stage_id, contact_id, company_id, amount, currency, probability, custom_fields, assigned_to, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id, created_at`,
		d.TenantID, d.Title, d.PipelineID, d.StageID, d.ContactID, d.CompanyID, d.Amount, d.Currency, d.Probability, cf, d.AssignedTo, d.CreatedBy).Scan(&d.ID, &d.CreatedAt)
}

func (r *Repository) GetDeal(ctx context.Context, id uuid.UUID) (*Deal, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	var d Deal; var cf []byte
	err := r.db.QueryRow(ctx, `SELECT id, tenant_id, title, pipeline_id, stage_id, contact_id, company_id, amount, currency, probability, custom_fields, assigned_to, created_at FROM deals WHERE id=$1 AND tenant_id=$2 AND deleted_at IS NULL`, id, tid).
		Scan(&d.ID, &d.TenantID, &d.Title, &d.PipelineID, &d.StageID, &d.ContactID, &d.CompanyID, &d.Amount, &d.Currency, &d.Probability, &cf, &d.AssignedTo, &d.CreatedAt)
	if err != nil { return nil, err }
	json.Unmarshal(cf, &d.CustomFields)
	return &d, nil
}

func (r *Repository) UpdateDeal(ctx context.Context, d *Deal) error {
	cf, _ := json.Marshal(d.CustomFields)
	_, err := r.db.Exec(ctx, `UPDATE deals SET title=$1, pipeline_id=$2, stage_id=$3, contact_id=$4, company_id=$5, amount=$6, probability=$7, custom_fields=$8, assigned_to=$9, updated_at=NOW() WHERE id=$10 AND tenant_id=$11 AND deleted_at IS NULL`,
		d.Title, d.PipelineID, d.StageID, d.ContactID, d.CompanyID, d.Amount, d.Probability, cf, d.AssignedTo, d.ID, d.TenantID)
	return err
}

func (r *Repository) DeleteDeal(ctx context.Context, id uuid.UUID, tid string) error {
	_, err := r.db.Exec(ctx, `UPDATE deals SET deleted_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tid)
	return err
}

// ─── Pipelines ───────────────────────────────────────────────

func (r *Repository) ListPipelines(ctx context.Context) ([]Pipeline, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, is_default, created_at FROM pipelines WHERE tenant_id=$1 ORDER BY sort_order`, tid)
	if err != nil { return nil, err }
	defer rows.Close()
	var pipelines []Pipeline
	for rows.Next() {
		var p Pipeline
		rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.IsDefault, &p.CreatedAt)
		pipelines = append(pipelines, p)
	}
	return pipelines, nil
}

func (r *Repository) CreatePipeline(ctx context.Context, p *Pipeline) error {
	return r.db.QueryRow(ctx, `INSERT INTO pipelines (tenant_id, name, is_default) VALUES ($1,$2,$3) RETURNING id, created_at`,
		p.TenantID, p.Name, p.IsDefault).Scan(&p.ID, &p.CreatedAt)
}

// ─── Activities ──────────────────────────────────────────────

func (r *Repository) ListActivities(ctx context.Context, p ListParams) (*ListResponse, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	if p.PageSize == 0 { p.PageSize = 20 }
	offset := (p.Page - 1) * p.PageSize
	var total int64
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM activities WHERE tenant_id=$1", tid).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id, tenant_id, type, subject, description, entity_type, entity_id, due_date, completed, assigned_to, created_by, created_at
		 FROM activities WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tid, p.PageSize, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var activities []Activity
	for rows.Next() {
		var a Activity
		rows.Scan(&a.ID, &a.TenantID, &a.Type, &a.Subject, &a.Description, &a.EntityType, &a.EntityID, &a.DueDate, &a.Completed, &a.AssignedTo, &a.CreatedBy, &a.CreatedAt)
		activities = append(activities, a)
	}
	return &ListResponse{Data: activities, Total: total, Page: p.Page, PageSize: p.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PageSize)))}, nil
}

func (r *Repository) CreateActivity(ctx context.Context, a *Activity) error {
	return r.db.QueryRow(ctx, `INSERT INTO activities (tenant_id, type, subject, description, entity_type, entity_id, due_date, assigned_to, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_at`,
		a.TenantID, a.Type, a.Subject, a.Description, a.EntityType, a.EntityID, a.DueDate, a.AssignedTo, a.CreatedBy).Scan(&a.ID, &a.CreatedAt)
}

// ─── Audit Logs ──────────────────────────────────────────────

func (r *Repository) CreateAuditLog(ctx context.Context, tenantID, entityType string, entityID uuid.UUID, action string, changedBy uuid.UUID, changes interface{}) error {
	ch, _ := json.Marshal(changes)
	_, err := r.db.Exec(ctx, `INSERT INTO audit_logs (tenant_id, entity_type, entity_id, action, changed_by, changes) VALUES ($1,$2,$3,$4,$5,$6)`,
		tenantID, entityType, entityID, action, changedBy, ch)
	return err
}

func (r *Repository) ListAuditLogs(ctx context.Context, p ListParams) (*ListResponse, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	if p.PageSize == 0 { p.PageSize = 50 }
	offset := (p.Page - 1) * p.PageSize
	var total int64
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM audit_logs WHERE tenant_id=$1", tid).Scan(&total)
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, entity_type, entity_id, action, changed_by, changes, created_at FROM audit_logs WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tid, p.PageSize, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	type AuditLog struct {
		ID         int64     `json:"id"`
		TenantID   string    `json:"tenant_id"`
		EntityType string    `json:"entity_type"`
		EntityID   uuid.UUID `json:"entity_id"`
		Action     string    `json:"action"`
		ChangedBy  uuid.UUID `json:"changed_by"`
		Changes    json.RawMessage `json:"changes"`
		CreatedAt  string    `json:"created_at"`
	}
	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		rows.Scan(&l.ID, &l.TenantID, &l.EntityType, &l.EntityID, &l.Action, &l.ChangedBy, &l.Changes, &l.CreatedAt)
		logs = append(logs, l)
	}
	return &ListResponse{Data: logs, Total: total, Page: p.Page, PageSize: p.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PageSize)))}, nil
}
