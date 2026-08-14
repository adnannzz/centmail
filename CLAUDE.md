# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

CentMail is a self-hosted newsletter/mailing-list manager: a Go backend (single binary) + a Vue 2/Buefy admin SPA, backed by PostgreSQL. It is a rebranded fork of [listmonk](https://listmonk.app) — the Go module path is `github.com/adnannzz/centmail`, but most of the architecture, patterns, and even some file/variable names still reflect the upstream project. The `origin` git remote points at the real upstream (`knadh/listmonk`) for reference; `centmail` is the actual push target (`adnannzz/centmail`).

## Commands

### Local dev stack (preferred way to run this project)
The dev stack runs everything in Docker with the app directory bind-mounted, so most changes don't need a rebuild:
```bash
make init-dev-docker   # first-time only: builds images, installs the DB
make dev-docker        # starts postgres, mailhog, the Vite frontend, and the Go backend
make rm-dev-docker      # tears everything down, including the DB volume
```
Visit `http://localhost:8080` (frontend dev server, proxies API calls to the backend on `:9000`).

- **Frontend hot-reloads** automatically (Vite + yarn watch) — no action needed after editing `frontend/src/**`.
- **Backend does NOT hot-reload.** After any change to a `.go` file, `i18n/*.json`, `schema.sql`, or `permissions.json`, restart the backend container (e.g. `docker restart dev-backend-1`) for it to take effect.
- Nightly/dev builds don't record their migration version in the DB, so every backend restart re-runs any pending `internal/migrations/vX.Y.Z.go` migrations — they must stay idempotent (`CREATE TABLE IF NOT EXISTS`, etc.).

### Building
```bash
make build          # compiles the Go backend to ./centmail
make build-frontend  # builds the Vue admin SPA into frontend/dist
make dist            # full build: backend + frontend + bundles all static assets into the binary via stuffbin
```

### Backend
```bash
go build ./...    # compile check
go vet ./...
gofmt -l .         # list files needing formatting; gofmt -w . to fix
make test          # go test ./... (no *_test.go files exist in this repo yet)
```
There is no `.golangci.yml` / lint config in this repo — `go vet` + `gofmt` is what CI (`build-sanity.yml`) effectively relies on via `make dist`.

### Frontend
```bash
cd frontend
yarn dev            # standalone Vite dev server (prefer `make dev-docker` instead, which wires up the backend too)
yarn lint           # eslint --ext .js,.vue src
yarn build          # production build
```
Cypress e2e tests live in `frontend/cypress/` but aren't wired into `package.json` scripts or CI — run via `npx cypress open`/`run` directly if needed. Note `frontend/cypress.config.js` still references the pre-rebrand binary name (`./listmonk`) and `LISTMONK_ADMIN_*` env vars; this was deliberately left alone during the rebrand since it's test-only tooling.

## Architecture

### Backend layering
`cmd/` (HTTP handlers, route registration, echo wiring) → `internal/core/` (business logic, one file per resource, e.g. `lists.go`, `forms.go`) → `models/` (structs with `db`/`json` tags) → `queries/*.sql` (raw SQL, goyesql-tagged).

- All `queries/*.sql` files are globbed and parsed by `cmd/init.go`'s `readQueries`, then matched to fields on the `models.Queries` struct (`internal/` doesn't touch SQL directly) by the `query:"name"` struct tag and prepared via `goyesqlx.ScanToStruct`. Adding a new query is: write it in a `queries/*.sql` file with a `-- name: your-query` header, add a matching `*sqlx.Stmt` field to `models.Queries`, then call it from `internal/core/`.
- Route permission gating uses the `pm(handler, "perm:string")` middleware wrapper in `cmd/handlers.go`. Permission strings are declared in two places that must stay in sync: `permissions.json` (the master list, loaded at startup and exposed to the frontend for the roles UI) and Go constants in `internal/auth/models.go` (`Perm*`). `auth.SuperAdminRoleID` (role ID 1) bypasses all permission checks.
- DB schema: `schema.sql` is the full schema applied on fresh installs. Upgrades run through `internal/migrations/vX.Y.Z.go` files registered in the `migList` slice in `cmd/upgrade.go` (ordered by semver, run sequentially from the DB's last recorded version). Nightly builds don't record the version, so migrations re-run every boot — write them idempotently.
- Public-facing (non-admin) pages — the subscription form, unsubscribe page, archive, etc. — are server-rendered Go `html/template`s under `static/public/templates/*.html`, completely separate from the Vue admin SPA. `cmd/public.go` renders them via `tplRenderer`/`subFormTpl`-style structs.

### Frontend
- Vue 2 + Buefy. Settings-style pages share one reactive `form`/`data` object passed as a prop into sibling tab components rather than each tab fetching its own copy (see `Settings.vue` + `views/settings/*.vue`).
- Vuex state (`store/index.js`) is auto-initialized for every key in the `models` object in `constants.js` — adding a new API resource usually just means adding one entry there plus `loading`/`store` options on the corresponding call in `api/index.js`.
- `api/index.js`'s http layer **auto-camelCases GET response JSON keys** by default (e.g. `list_ids` → `listIds`) for consistency with the Vue/AirBnB lint spec, but this only applies to responses — outgoing POST/PUT bodies must be written with the backend's actual snake_case JSON field names. A handful of endpoints (`/api/config`, `/api/settings`, campaign headers) opt out via `camelCase: false` or a path-based `camelCase` test function.
- Icons: this project ships a **hand-curated Fontello icon font** (`frontend/src/assets/icons/fontello.css`), not the full Material Design Icons set — only ~44 glyphs are actually compiled in (see `frontend/fontello/config.json`'s `selected: true` entries). Using an MDI icon name that isn't in that curated set silently renders nothing (no error). To add more icons, see the instructions in `frontend/README.md`.
- i18n strings live in `i18n/*.json` at the repo root (not under `frontend/`) and are served by the backend via `/api/lang/:lang`, fetched once at app init (`main.js`) — editing them requires a backend restart in the dev stack to take effect, not just a frontend reload.

### Rebrand-specific notes
Real upstream infrastructure this fork doesn't own was deliberately left pointing at listmonk during the rebrand rather than being renamed to something fictional: the `docs/` directory (Hugo marketing site, mkdocs, Postman/swagger collections), `.goreleaser*.yml` and `.github/workflows/nightly.yml`'s Docker image tags (`ghcr.io/knadh/listmonk`), and outbound docs links in the admin UI (`https://listmonk.app/docs/...`). Don't "fix" these to say centmail.app/etc. without setting up the corresponding real infrastructure first.
