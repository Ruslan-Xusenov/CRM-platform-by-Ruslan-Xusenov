package crm

import (
	"time"
	"github.com/google/uuid"
)

type Lead struct {
	ID           uuid.UUID              `json:"id"`
	TenantID     uuid.UUID              `json:"tenant_id"`
	Title        string                 `json:"title"`
	Status       string                 `json:"status"`
	Source       string                 `json:"source,omitempty"`
	Budget       *float64               `json:"budget,omitempty"`
	Currency     string                 `json:"currency"`
	ContactName  string                 `json:"contact_name,omitempty"`
	ContactPhone string                 `json:"contact_phone,omitempty"`
	ContactEmail string                 `json:"contact_email,omitempty"`
	CompanyName  string                 `json:"company_name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields"`
	AssignedTo   *uuid.UUID             `json:"assigned_to,omitempty"`
	CreatedBy    *uuid.UUID             `json:"created_by,omitempty"`
	ConvertedAt  *time.Time             `json:"converted_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
	DeletedAt    *time.Time             `json:"-"`
}

type Contact struct {
	ID           uuid.UUID              `json:"id"`
	TenantID     uuid.UUID              `json:"tenant_id"`
	CompanyID    *uuid.UUID             `json:"company_id,omitempty"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	Email        string                 `json:"email,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	Mobile       string                 `json:"mobile,omitempty"`
	Position     string                 `json:"position,omitempty"`
	Source       string                 `json:"source,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields"`
	AssignedTo   *uuid.UUID             `json:"assigned_to,omitempty"`
	CreatedBy    *uuid.UUID             `json:"created_by,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
	DeletedAt    *time.Time             `json:"-"`
}

type Company struct {
	ID            uuid.UUID              `json:"id"`
	TenantID      uuid.UUID              `json:"tenant_id"`
	Name          string                 `json:"name"`
	Industry      string                 `json:"industry,omitempty"`
	Website       string                 `json:"website,omitempty"`
	Phone         string                 `json:"phone,omitempty"`
	Email         string                 `json:"email,omitempty"`
	Address       string                 `json:"address,omitempty"`
	EmployeeCount *int                   `json:"employee_count,omitempty"`
	Revenue       *float64               `json:"revenue,omitempty"`
	CustomFields  map[string]interface{} `json:"custom_fields"`
	CreatedBy     *uuid.UUID             `json:"created_by,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at,omitempty"`
	DeletedAt     *time.Time             `json:"-"`
}

type Deal struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenant_id"`
	Title             string                 `json:"title"`
	PipelineID        *uuid.UUID             `json:"pipeline_id,omitempty"`
	StageID           *uuid.UUID             `json:"stage_id,omitempty"`
	ContactID         *uuid.UUID             `json:"contact_id,omitempty"`
	CompanyID         *uuid.UUID             `json:"company_id,omitempty"`
	Amount            *float64               `json:"amount,omitempty"`
	Currency          string                 `json:"currency"`
	Probability       int                    `json:"probability"`
	ExpectedCloseDate *time.Time             `json:"expected_close_date,omitempty"`
	Won               *bool                  `json:"won,omitempty"`
	CustomFields      map[string]interface{} `json:"custom_fields"`
	AssignedTo        *uuid.UUID             `json:"assigned_to,omitempty"`
	CreatedBy         *uuid.UUID             `json:"created_by,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	DeletedAt         *time.Time             `json:"-"`
}

type Pipeline struct {
	ID        uuid.UUID       `json:"id"`
	TenantID  uuid.UUID       `json:"tenant_id"`
	Name      string          `json:"name"`
	IsDefault bool            `json:"is_default"`
	Stages    []PipelineStage `json:"stages,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type PipelineStage struct {
	ID         uuid.UUID `json:"id"`
	PipelineID uuid.UUID `json:"pipeline_id"`
	Name       string    `json:"name"`
	Color      string    `json:"color"`
	SortOrder  int       `json:"sort_order"`
	IsWon      bool      `json:"is_won"`
	IsLost     bool      `json:"is_lost"`
}

type Activity struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Type        string     `json:"type"`
	Subject     string     `json:"subject"`
	Description string     `json:"description,omitempty"`
	EntityType  string     `json:"entity_type,omitempty"`
	EntityID    *uuid.UUID `json:"entity_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Completed   bool       `json:"completed"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ListParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
	Search   string `json:"search"`
	Status   string `json:"status"`
}

type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
