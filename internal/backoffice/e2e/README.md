# Backoffice end-to-end check

Optional browser check for the backoffice UI. It is **not** part of
`go test ./...` — the Go suite covers the services, this covers the pages.

It drives a real browser through the whole flow: login, dashboard, the three
CRUDs, the role-permission manager, theme toggle, the auth guard and a bad
login. It fails on any console error or unexpected failed request, so it
catches the class of bug that only shows up once the JS actually runs.

## Running it

Uses `playwright-core` against the Chrome you already have installed — no
browser download.

```bash
# 1. start a server with a seeded database
DB_PATH=/tmp/e2e.db go run ./cmd/migrate
DB_PATH=/tmp/e2e.db go run ./cmd/seed
DB_PATH=/tmp/e2e.db PORT=8099 go run ./cmd/api

# 2. in another terminal
cd internal/backoffice/e2e
npm install
npm test
```

It writes screenshots to `internal/backoffice/e2e/shots/` (gitignored) —
useful for eyeballing the result of a styling change.

## Configuration

| Env var | Default |
|---|---|
| `BASE_URL` | `http://localhost:8099` |
| `ADMIN_EMAIL` | `admin@admin.com` |
| `ADMIN_PASSWORD` | `admin123` |
| `CHROME_PATH` | macOS Google Chrome location |

The test creates one application (`pw-<timestamp>`) and leaves it behind, so
point it at a throwaway database rather than anything you care about.
