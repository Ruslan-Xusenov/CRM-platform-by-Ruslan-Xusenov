-- ==============================================================
-- Migration 002: CRM Schema — Leads, Contacts, Companies, Deals
-- ==============================================================

-- ─── Pipelines ───────────────────────────────────────────────
CREATE TABLE pipelines (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        TEXT NOT NULL,
    is_default  BOOLEAN DEFAULT false,
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE pipelines ENABLE ROW LEVEL SECURITY;

CREATE TABLE pipeline_stages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    color       TEXT DEFAULT '#3B82F6',
    sort_order  INT DEFAULT 0,
    is_won      BOOLEAN DEFAULT false,
    is_lost     BOOLEAN DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ─── Companies ───────────────────────────────────────────────
CREATE TABLE companies (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    name          TEXT NOT NULL,
    industry      TEXT,
    website       TEXT,
    phone         TEXT,
    email         TEXT,
    address       TEXT,
    city          TEXT,
    country       TEXT,
    employee_count INT,
    revenue       DECIMAL(15,2),
    custom_fields JSONB DEFAULT '{}',
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
CREATE INDEX idx_companies_tenant ON companies(tenant_id, deleted_at);
CREATE INDEX idx_companies_custom ON companies USING GIN (custom_fields jsonb_path_ops);

-- ─── Contacts ────────────────────────────────────────────────
CREATE TABLE contacts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    company_id    UUID REFERENCES companies(id),
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    email         TEXT,
    phone         TEXT,
    mobile        TEXT,
    position      TEXT,
    source        TEXT,
    avatar_url    TEXT,
    custom_fields JSONB DEFAULT '{}',
    assigned_to   UUID REFERENCES users(id),
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

ALTER TABLE contacts ENABLE ROW LEVEL SECURITY;
CREATE INDEX idx_contacts_tenant ON contacts(tenant_id, deleted_at);
CREATE INDEX idx_contacts_phone ON contacts(phone);
CREATE INDEX idx_contacts_custom ON contacts USING GIN (custom_fields jsonb_path_ops);

-- ─── Leads ───────────────────────────────────────────────────
CREATE TABLE leads (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    title         TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'new',
    source        TEXT,
    budget        DECIMAL(15,2),
    currency      TEXT DEFAULT 'UZS',
    contact_name  TEXT,
    contact_phone TEXT,
    contact_email TEXT,
    company_name  TEXT,
    description   TEXT,
    custom_fields JSONB DEFAULT '{}',
    assigned_to   UUID REFERENCES users(id),
    created_by    UUID REFERENCES users(id),
    converted_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

ALTER TABLE leads ENABLE ROW LEVEL SECURITY;
CREATE INDEX idx_leads_tenant_status ON leads(tenant_id, status, created_at DESC);
CREATE INDEX idx_leads_assigned ON leads(assigned_to) WHERE deleted_at IS NULL;
CREATE INDEX idx_leads_custom ON leads USING GIN (custom_fields jsonb_path_ops);

-- ─── Deals ───────────────────────────────────────────────────
CREATE TABLE deals (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    title         TEXT NOT NULL,
    pipeline_id   UUID REFERENCES pipelines(id),
    stage_id      UUID REFERENCES pipeline_stages(id),
    contact_id    UUID REFERENCES contacts(id),
    company_id    UUID REFERENCES companies(id),
    amount        DECIMAL(15,2),
    currency      TEXT DEFAULT 'UZS',
    probability   INT DEFAULT 0,
    expected_close_date DATE,
    closed_at     TIMESTAMPTZ,
    won           BOOLEAN,
    custom_fields JSONB DEFAULT '{}',
    assigned_to   UUID REFERENCES users(id),
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

ALTER TABLE deals ENABLE ROW LEVEL SECURITY;
CREATE INDEX idx_deals_tenant_pipeline ON deals(tenant_id, pipeline_id, stage_id);
CREATE INDEX idx_deals_custom ON deals USING GIN (custom_fields jsonb_path_ops);

-- ─── Activities ──────────────────────────────────────────────
CREATE TABLE activities (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    type          TEXT NOT NULL,
    subject       TEXT NOT NULL,
    description   TEXT,
    entity_type   TEXT,
    entity_id     UUID,
    due_date      TIMESTAMPTZ,
    completed     BOOLEAN DEFAULT false,
    completed_at  TIMESTAMPTZ,
    assigned_to   UUID REFERENCES users(id),
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activities_tenant ON activities(tenant_id, entity_type, entity_id);

-- ─── Notes ───────────────────────────────────────────────────
CREATE TABLE notes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    entity_type TEXT NOT NULL,
    entity_id   UUID NOT NULL,
    content     TEXT NOT NULL,
    created_by  UUID REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);

CREATE INDEX idx_notes_entity ON notes(tenant_id, entity_type, entity_id);

-- ─── Insert default pipeline ─────────────────────────────────
-- (This will be done per-tenant during registration)
