---
implementation: unshipped
probe: "/api/links"
scope: internal/
---

# Design record: shortlink

**Status**: decided 2026-10-19 with the maintainer, when the initiative was scoped.
**Phases**: [phase-01, phase-02, phase-03](../scope.md).

## 1. Motivation

The repository exists to show Zurdo and Carol driving an initiative end to end:
a phase is accepted on GitHub, drafted, reviewed, run and reported with humans
deciding only at the gates. It needs a product small enough that each phase runs
in minutes, and that exercises both kinds of criteria a reader cares about:
`[http:]` checks against a live server and `[shell:]` checks against a CLI.

## 2. What already ships

| Surface | Where | What it does today |
|---|---|---|
| `shortlink` binary | `cmd/shortlink/main.go` | `version` and `help`; any other command exits 2 |
| Version | `internal/version/version.go` | `version.Version`, set by `-ldflags` |
| Dev server | `scripts/devserver.sh` | Rebuilds and restarts `shortlink serve --addr 127.0.0.1:8787 --data <copy of testdata/seed.json>` on every `.go` change, following the newest `.carol/worktrees/<phase>` |
| Tests | `go test ./...` | Two tests on the CLI dispatcher, one on the version |

Nothing serves HTTP and nothing stores a link.

## 3. Goals and non-goals

Goals:

- One binary. `shortlink serve` is the API; every other subcommand is a client of it.
- Standard library only (`net/http`, `encoding/json`, `flag`), Go 1.22 routing
  patterns (`mux.HandleFunc("GET /{slug}", …)`).
- Every phase's exit criteria are machine-checkable: `[http:]` against the dev
  server on `127.0.0.1:8787`, `[shell:]` for `go test` and the CLI.

Non-goals:

- Authentication, users, and link expiry. They add no new kind of criterion. Reopen
  if the demo needs a `[manual]` sign-off story.
- A database. A JSON file is enough for a demo, and keeps the run hermetic.
- A web UI. `[http:]` checks status codes and bodies; HTML adds nothing to show.

## 4. The decision

Packages, by layer:

- `internal/store` — `Store`, an in-memory map of slug → `Link` guarded by a
  mutex and persisted to one JSON file (write to a temp file, then rename).
  `Link` is `{slug, url}` in phase 1. Fields arrive with the phase that
  observes them: phase 3 adds `clicks` and `last_clicked_at`. There is no
  `created_at`, because no phase reads or checks it (amended 2026-10-19 with
  the maintainer, on the phase 1 PRD review).
- `internal/server` — an `http.Handler` over a `Store`. The API lives under
  `/api/links`; a bare `/{slug}` redirects.
- `internal/client` — a small HTTP client for the API, used by the CLI.
- `cmd/shortlink` — subcommand dispatch only.

Conventions that every phase keeps:

- The server address is `--addr` (default `127.0.0.1:8787`); the client reads
  `--server` or `SHORTLINK_SERVER` (default `http://127.0.0.1:8787`).
- JSON responses carry `Content-Type: application/json`. Errors are
  `{"error": "<message>"}` with a 4xx status.
- `testdata/seed.json` holds the seed links the dev server loads on every
  restart. It contains at least `zurdo` → `https://github.com/ElOrlis/zurdo` and
  `carol` → `https://github.com/ElOrlis/carol`, so `[http:]` criteria have known
  data to read.
- Slugs match `^[a-z0-9-]{3,32}$`. A generated slug is 6 lowercase letters/digits.

## 5. Affected areas

`cmd/shortlink/`, `internal/store/`, `internal/server/`, `internal/client/`,
`testdata/seed.json`, `README.md` (usage section).

## Phase 1 — link store and HTTP API

**Goal.** `shortlink serve` runs an HTTP API that creates, lists, reads and
resolves short links, persisted to a JSON file.

### Surface

- `shortlink serve [--addr 127.0.0.1:8787] [--data links.json]`. A missing data
  file starts an empty store; an unreadable one exits 1 with a message.
- `GET /healthz` → `200`, body `ok`.
- `POST /api/links` with `{"url": "...", "slug": "optional"}` → `201` and the
  link as JSON. A URL that is not absolute `http`/`https` → `400`. A slug taken
  already → `409`. A slug that breaks the pattern → `400`.
- `GET /api/links` → `200`, a JSON array of every link, sorted by slug.
- `GET /api/links/{slug}` → `200` and the link, or `404`.
- `GET /{slug}` → `302` with `Location` set to the link's URL, or `404`.
- `testdata/seed.json` with the two seed links in §4.

### Exit criteria (observables)

- `go test ./...` passes, with table-driven `httptest` tests for every route
  and status above, and a store test that a saved file reloads to the same links.
- With the dev server running: `GET http://127.0.0.1:8787/healthz` → 200
  containing `ok`; `GET http://127.0.0.1:8787/api/links` → 200 containing
  `"slug":"zurdo"`; `GET http://127.0.0.1:8787/api/links/carol` → 200 containing
  `https://github.com/ElOrlis/carol`; `GET http://127.0.0.1:8787/api/links/nope`
  → 404.
- `curl -s -o /dev/null -w '%{http_code} %{redirect_url}' http://127.0.0.1:8787/zurdo`
  prints `302 https://github.com/ElOrlis/zurdo` (a `[shell:]` check, so a
  redirect is not followed).

### Depends on

Nothing.

## Phase 2 — CLI client

**Goal.** The other `shortlink` subcommands drive the API, so a person never
needs `curl`.

### Surface

- `shortlink add <url> [--slug <slug>]` prints the short URL
  (`http://127.0.0.1:8787/<slug>`) and exits 0. A rejected request prints the
  server's error to stderr and exits 1.
- `shortlink list [--json]` prints one `<slug>\t<url>` line per link, sorted, or
  the raw JSON array with `--json`.
- `shortlink resolve <slug>` prints the target URL, or exits 1 with
  `shortlink: no link "<slug>"` on stderr.
- Every client command takes `--server`, falling back to `SHORTLINK_SERVER` and
  then `http://127.0.0.1:8787`. An unreachable server exits 3 with a message
  naming the address.
- `shortlink help` lists every subcommand.

### Exit criteria (observables)

- `go test ./...` passes, with client tests against an `httptest.Server` for each
  command, including the unreachable-server exit 3.
- With the dev server running: `go run ./cmd/shortlink list` prints a line
  starting `zurdo` followed by a tab; `go run ./cmd/shortlink resolve carol` prints
  `https://github.com/ElOrlis/carol`; `go run ./cmd/shortlink resolve nope` exits 1.
- `go run ./cmd/shortlink add https://go.dev --slug godev` exits 0, and then
  `GET http://127.0.0.1:8787/api/links/godev` → 200 containing `https://go.dev`.
- `SHORTLINK_SERVER=http://127.0.0.1:1 go run ./cmd/shortlink list` exits 3.

### Depends on

Phase 1: the API it calls.

## Phase 3 — click stats

**Goal.** Every redirect is counted, and the count is readable from both the API
and the CLI.

### Surface

- `Link` gains `clicks` (0 for a link loaded without it) and `last_clicked_at`
  (null until the first redirect), both persisted with the link (§4).
- `GET /{slug}` increments the link's `clicks` before redirecting; the count is
  persisted with the link. Concurrent redirects never lose a count.
- `GET /api/links/{slug}/stats` → `200` and `{"slug": "...", "clicks": N,
  "last_clicked_at": "<RFC 3339 or null>"}`, or `404`.
- `GET /api/stats` → `200` and `{"links": N, "clicks": N}` over every link.
- `shortlink stats [slug]` prints `<slug>\t<clicks>` for one link, or for every
  link sorted by clicks (most first) with no argument.

### Exit criteria (observables)

- `go test ./...` passes, including a test that 100 concurrent redirects to one
  slug leave `clicks` at exactly 100 (run with `-race`).
- With the dev server freshly started: `GET http://127.0.0.1:8787/api/links/zurdo/stats`
  → 200 containing `"clicks":0`; after one `curl -s -o /dev/null http://127.0.0.1:8787/zurdo`,
  `go run ./cmd/shortlink stats zurdo` prints `zurdo	1`.
- `GET http://127.0.0.1:8787/api/stats` → 200 containing `"links":`.
- `GET http://127.0.0.1:8787/api/links/nope/stats` → 404.

### Depends on

Phase 1 for the redirect, phase 2 for the CLI.

## 6. Considered and rejected

- **A separate client binary.** Two binaries double the release and install
  surface for no new kind of check; one binary with `serve` keeps a phase's
  `[shell:]` and `[http:]` criteria on the same build.
- **Starting the server inside each criterion** (`[shell: go run … serve & …]`).
  It makes every `[http:]` criterion a shell script and hides what is being
  shown. The dev server keeps `[http:]` criteria declarative.
- **SQLite.** It needs cgo or a third-party driver, which breaks the
  standard-library-only goal for a store of a handful of links.

## 7. Open questions

1. Whether a phase should add a `[manual]` criterion so the demo shows
   `zurdo review --sign-off` (answered by the maintainer at a phase review).
2. Whether redirects should be `302` or `301`. `302` is chosen so browsers do not
   cache past the click counter; reopen if a phase needs permanent links.
