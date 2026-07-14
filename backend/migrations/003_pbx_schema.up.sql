-- ==============================================================
-- Migration 003: PBX Schema — Extensions, CDR, Routing
-- ==============================================================

-- ─── Extensions ──────────────────────────────────────────────
CREATE TABLE extensions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(id),
    user_id          UUID REFERENCES users(id),
    extension_number TEXT NOT NULL,
    display_name     TEXT,
    password         TEXT NOT NULL,
    transport        TEXT DEFAULT 'wss',
    context          TEXT DEFAULT 'crm-stasis',
    enabled          BOOLEAN DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ,
    UNIQUE(tenant_id, extension_number)
);

CREATE INDEX idx_ext_tenant ON extensions(tenant_id);
CREATE INDEX idx_ext_user ON extensions(user_id);

-- ─── Call Detail Records (partitioned by month) ──────────────
CREATE TABLE call_logs (
    id               UUID NOT NULL DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL,
    caller           TEXT NOT NULL,
    callee           TEXT NOT NULL,
    direction        TEXT NOT NULL,
    status           TEXT NOT NULL,
    started_at       TIMESTAMPTZ NOT NULL,
    answered_at      TIMESTAMPTZ,
    ended_at         TIMESTAMPTZ,
    duration_seconds INT DEFAULT 0,
    recording_url    TEXT,
    linked_entity_type TEXT,
    linked_entity_id   UUID,
    asterisk_unique_id TEXT,
    metadata         JSONB DEFAULT '{}',
    PRIMARY KEY (id, started_at)
) PARTITION BY RANGE (started_at);

-- Monthly partitions
CREATE TABLE call_logs_2026_07 PARTITION OF call_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE call_logs_2026_08 PARTITION OF call_logs
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE call_logs_2026_09 PARTITION OF call_logs
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE call_logs_2026_10 PARTITION OF call_logs
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE call_logs_2026_11 PARTITION OF call_logs
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE call_logs_2026_12 PARTITION OF call_logs
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

CREATE INDEX idx_cdr_tenant_date ON call_logs(tenant_id, started_at DESC);
CREATE INDEX idx_cdr_caller ON call_logs(caller, started_at DESC);
CREATE INDEX idx_cdr_callee ON call_logs(callee, started_at DESC);

-- ─── Routing Rules ───────────────────────────────────────────
CREATE TABLE routing_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        TEXT NOT NULL,
    pattern     TEXT NOT NULL,
    action      TEXT NOT NULL,
    target      TEXT NOT NULL,
    priority    INT DEFAULT 0,
    enabled     BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ─── Ring Groups ─────────────────────────────────────────────
CREATE TABLE ring_groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        TEXT NOT NULL,
    strategy    TEXT DEFAULT 'ringall',
    timeout     INT DEFAULT 30,
    members     JSONB DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
