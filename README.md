# SAMS — System Architecture Modeling & Simulation

A developer-centric platform that **models, validates, and simulates** system
architectures instead of just drawing diagrams. Drag components onto a canvas,
connect them, describe your workload, and SAMS checks the design against a
library of architecture rules — flagging security holes, scalability bottlenecks,
reliability gaps, and performance problems in real time, with a health score.

---

## Features

- **Visual canvas** — drag-and-drop components (frontend, backend, database,
  cache, queue, load balancer, API gateway, CDN, object storage, monitoring…)
  and wire them together with typed connections (sync/async, HTTP/HTTPS/gRPC/AMQP).
- **Live validation** — the design is re-checked automatically (debounced) as you
  edit; issues highlight the offending nodes and edges directly on the canvas.
- **Context-aware rules** — the same topology is judged differently based on your
  declared workload: traffic level, read/write ratio, user base, team size, and
  product stage.
- **Architecture health score** — a single 0–100 number, weighted by rule severity
  and scaled by traffic, so you can see a design improve as you fix issues.
- **Save / Load** — diagrams auto-persist to the browser (survive a refresh) and can
  be exported/imported as JSON files for sharing or version control.
- **Declarative rule library** — rules live as YAML files under `rules/`, grouped by
  category, and are loaded at startup.

---

## Architecture

```
┌─────────────────────────┐        POST /api/v1/validate        ┌──────────────────────────┐
│  Frontend (React + Vite) │  ───────────────────────────────▶  │   Backend (Go + chi)      │
│                          │                                     │                           │
│  • React Flow canvas     │  ◀───────────────────────────────  │  • Rule engine            │
│  • Zustand store (SoT)    │        validation results          │  • YAML rule loader        │
│  • Auto-persist + export  │                                     │  • Scoring                │
└─────────────────────────┘                                     └────────────┬─────────────┘
                                                                             │ loads at startup
                                                                             ▼
                                                                     rules/*.yaml (29 rules)
```

- **Backend** (`/backend`) — Go service using the [chi](https://github.com/go-chi/chi)
  router. It loads rules from YAML at startup, builds a graph from the posted
  nodes/edges, runs every rule against it, and returns errors, warnings, passed
  rules, and a score. Middleware provides CORS, structured request logging, a
  token-bucket rate limiter, request IDs, and panic recovery.
- **Frontend** (`/frontend`) — React 19 + TypeScript + Vite, with
  [React Flow](https://reactflow.dev) for the canvas and
  [Zustand](https://github.com/pmndrs/zustand) as the single source of truth for
  the diagram (controlled mode). Styling is Tailwind CSS v4.
- **Rules** (`/rules`) — one YAML file per rule, organized into `security/`,
  `performance/`, `scalability/`, and `reliability/`.

---

## Getting started

### Option A — Docker Compose (everything at once)

```bash
docker compose up --build
```

- Frontend → http://localhost:3000
- Backend  → http://localhost:8080

### Option B — Local development

**Backend:**

```bash
cd backend
go run ./cmd/server
# listens on :8080, loads rules from ../rules by default
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
# Vite dev server on :3000, proxies /api → http://localhost:8080
```

Open http://localhost:3000 and start dragging components onto the canvas.

---

## Configuration (backend)

The backend is configured via environment variables — see
[`.env.example`](.env.example) for a template (`cp .env.example .env`). Every
variable has a sensible default, so no configuration is required to run locally.

| Variable          | Default                  | Description                                   |
| ----------------- | ------------------------ | --------------------------------------------- |
| `PORT`            | `8080`                   | Port to listen on                             |
| `RULES_PATH`      | `../rules`               | Rule YAML dir (relative to `backend/`)        |
| `ALLOWED_ORIGINS` | `http://localhost:3000`  | Comma-separated CORS allow-list               |
| `ENV`             | `development`            | Environment label (used in logs)              |

The frontend needs no build-time API URL: it calls a relative `/api` path, which
Vite proxies in development and nginx proxies in the Docker image.

---

## API

Base path: `/api/v1`

| Method | Endpoint      | Description                                              |
| ------ | ------------- | ------------------------------------------------------- |
| `POST` | `/validate`   | Validate a diagram; returns errors, warnings, score     |
| `GET`  | `/components` | The draggable component palette                         |
| `GET`  | `/rules`      | All loaded rules and their metadata                     |
| `GET`  | `/health`     | Liveness + number of rules loaded                       |


---

## Rule system

Each rule is a YAML file (e.g. `rules/performance/rule_001_missing_cache.yaml`):

```yaml
id: rule_001
category: performance
title: "Missing Cache Layer"
description: "Your backend queries the database directly under high traffic..."
suggestion: "Add a Redis or Memcached node between your backend and database..."
baseWeight: 10
conditions:
  contextChecks:
    trafficLevel: "high"
  requiredNodes: [backend, database]
  forbiddenNodes: [cache]
```

Rules are grouped into four categories:

| Category       | Examples                                                              |
| -------------- | -------------------------------------------------------------------- |
| **security**   | frontend→DB direct connection, no auth, SQL injection risk, no HTTPS |
| **performance**| missing cache, N+1 queries, missing indexes, no CDN                   |
| **scalability**| no queue, stateful servers, premature microservices, no autoscaling  |
| **reliability**| single point of failure, no replication, no backups, no monitoring   |

**Scoring:** each fired rule subtracts its `baseWeight` from 100 (warnings count
half), scaled by a traffic multiplier, and clamped to `[0, 100]`.

---

## Project structure

```
SAMS/
├── backend/
│   ├── cmd/server/        # main entrypoint
│   ├── config/            # env-based configuration
│   ├── internal/
│   │   ├── engine/        # graph builder, rule validator, scorer, YAML loader
│   │   ├── handlers/      # HTTP handlers (validate, components, rules, health)
│   │   ├── middleware/    # CORS, logging, rate limiting
│   │   └── models/        # domain types (node, edge, graph, rule, result)
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/     # canvas, sidebar, layout, validation panels
│   │   ├── store/          # Zustand store (single source of truth + persistence)
│   │   ├── hooks/          # debounced auto-validation
│   │   ├── services/       # API client
│   │   └── types/          # shared TypeScript types
│   └── Dockerfile
├── rules/                  # YAML rule library (security/performance/scalability/reliability)
├── docker-compose.yml
└── .github/workflows/ci.yml
```

CI runs on every push and pull request to `main` — see
[`.github/workflows/ci.yml`](.github/workflows/ci.yml).

---