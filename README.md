# Omnichannel CRM & WebRTC PBX Platform

Enterprise-grade CRM platform with integrated WebRTC softphone, built with Go, Next.js, and Asterisk PBX.

## 🏗 Architecture

```
┌─────────────┐     ┌──────────────┐     ┌───────────────┐
│   Next.js    │────▶│   Nginx      │────▶│   Go Backend  │
│   Frontend   │     │ (HTTPS/WSS)  │     │   (Chi/REST)  │
│  SIP.js      │     └──────────────┘     └───────┬───────┘
│  WebRTC      │                                  │
└─────────────┘                    ┌──────────────┼──────────────┐
                                   │              │              │
                              ┌────▼───┐    ┌─────▼────┐   ┌────▼───┐
                              │ Postgres│    │  Redis   │   │RabbitMQ│
                              │  (RLS)  │    │(Cache/PS)│   │(Events)│
                              └────────┘    └──────────┘   └────────┘
                                                                │
                              ┌──────────┐    ┌──────────┐     │
                              │ Asterisk │    │  MinIO   │◀────┘
                              │ 20 (ARI) │    │(S3 Files)│
                              └──────────┘    └──────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local development)
- Node.js 20+ (for local development)

### Start All Services

```bash
# Copy environment file
cp .env.example .env

# Start everything
make dev
```

### Access Points

| Service    | URL                          |
|-----------|------------------------------|
| Frontend  | http://localhost:3000         |
| Backend   | http://localhost:8080         |
| API Docs  | http://localhost:8080/health  |
| RabbitMQ  | http://localhost:15672        |
| MinIO     | http://localhost:9001         |
| Grafana   | http://localhost:3001         |
| Prometheus| http://localhost:9090         |

## 📂 Project Structure

```
├── backend/           # Go modular monolith (DDD)
│   ├── cmd/server/    # Application entrypoint
│   ├── internal/      # Domain modules (auth, crm, pbx, storage)
│   ├── pkg/           # Shared packages (database, cache, broker, ws)
│   └── migrations/    # PostgreSQL migrations
├── frontend/          # Next.js 14 (App Router)
├── asterisk/          # Asterisk PBX configuration
│   ├── conf/          # PJSIP, extensions, ARI, RTP configs
│   └── scripts/       # Recording upload scripts
├── nginx/             # Reverse proxy config
├── monitoring/        # Prometheus, Grafana, Loki
└── docker-compose.yml # Full stack orchestration
```

## 🔑 API Endpoints

### Auth
- `POST /api/v1/auth/register` — Create tenant + admin user
- `POST /api/v1/auth/login` — Login, returns JWT tokens
- `POST /api/v1/auth/refresh` — Refresh access token
- `POST /api/v1/auth/logout` — Revoke refresh token

### CRM
- `GET/POST /api/v1/leads` — List/Create leads
- `GET/PUT/DELETE /api/v1/leads/:id` — Lead CRUD
- `POST /api/v1/leads/:id/convert` — Convert to contact+deal
- `GET/POST /api/v1/contacts` — Contacts
- `GET/POST /api/v1/companies` — Companies
- `GET/POST /api/v1/deals` — Deals
- `GET/POST /api/v1/pipelines` — Pipeline management
- `GET /api/v1/audit-logs` — Audit trail

### PBX
- `POST /api/v1/calls/originate` — Click-to-call
- `GET /api/v1/calls/active` — Active calls
- `GET /api/v1/calls/history` — CDR records
- `GET/POST /api/v1/pbx/extensions` — Extensions
- `GET/PUT /api/v1/pbx/routing` — Routing rules

### WebSocket
- `GET /ws` — Real-time events (calls, lead updates, notifications)

## 🔒 Security

- **JWT** with access/refresh token rotation
- **RBAC** roles: admin, manager, operator, viewer
- **Row Level Security** (PostgreSQL) for multi-tenancy
- **SRTP + WSS** for encrypted telephony
- **Network isolation** via Docker networks
- **Nginx** as sole public entry point

## 📊 Monitoring

- **Prometheus** — Metrics collection
- **Grafana** — Dashboards
- **Loki + Promtail** — Centralized logging
- **TraceID** — Request chain tracing

## 📝 License

Proprietary — All rights reserved.