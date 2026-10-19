# PRD: Phase 1 — link store and HTTP API

<!-- zurdo:scaffold unreviewed — run the zurdo-prd-author skill before zurdo run -->

`shortlink serve` runs an HTTP API that creates, lists, reads and resolves
short links, persisted to a JSON file. Source: `docs/shortlink/scope.md`,
Phase 1. Three tasks: the store (`internal/store`), the routes
(`internal/server`), and the `serve` command with seed data, proven live
against the dev server on `127.0.0.1:8787`.

## Task: task-01 — JSON-file link store
**Effort**: medium
**Depends-on**: []
**Frozen**: internal/server/preauthored_task_02_test.go

### Requirements

- req-create: `Create` adds a link only when the URL is absolute `http`/`https` and the slug matches `^[a-z0-9-]{3,32}$` and is not taken. Otherwise it returns `ErrInvalidURL`, `ErrInvalidSlug` or `ErrSlugTaken` and stores nothing. An empty slug gets a generated one of 6 lowercase letters/digits.
- req-persist: `Open` on a missing file starts an empty store. On a directory or malformed JSON it returns an error. Every successful `Create` saves the file, a reopened file lists the same links, and `List` sorts by slug.

### Description

Implement the link store in `internal/store/store.go`. The file already holds
the signature stubs (`Link`, `ErrInvalidURL`, `ErrInvalidSlug`, `ErrSlugTaken`,
`Store`, `Open`, `Create`, `Get`, `List`). Each stub currently panics with
`todo: task-01`. Replace those bodies and add whatever fields `Store` needs.
Keep the exported signatures as they are.

- The data file is a JSON array of link objects, e.g.
  `[{"slug":"zurdo","url":"https://github.com/ElOrlis/zurdo"}]`. `Link` stays
  exactly `{slug, url}` for this phase.
- A slug must match `^[a-z0-9-]{3,32}$` (`docs/shortlink/design/shortlink.md` §4).
  A generated slug is 6 lowercase letters/digits and must not collide with an
  existing one.
- Writes save the whole file. Write to a temp file in the same directory and
  rename it over the original. Guard the store with a mutex, because the HTTP
  server calls it concurrently.
- Standard library only (Go 1.22).

The tests in `internal/store/preauthored_task_01_test.go` are evidence, not
scratch. Do not delete them, weaken their assertions, or change what they
assert. The only edit this task may make to that file is removing the
`//go:build preauthored` line (and the blank line after it) once the behavior
is implemented.

### Acceptance Criteria

- [ ] Create validates URLs and slugs against the pattern and its length bounds, rejects taken slugs, generates 6-character slugs and stores nothing on rejection [shell: go test -tags=preauthored -run '^TestStoreCreate$' -count=1 ./internal/store/] [proves: req-create]
- [ ] a missing file opens empty, an unreadable or malformed file errors, and a saved file reloads to the same sorted links [shell: go test -tags=preauthored -run '^TestStore(OpenMissing|OpenUnreadable|Reload)$' -count=1 ./internal/store/] [proves: req-persist]
- [ ] the pre-authored store tests are no longer build-tagged [no-grep: //go:build preauthored in internal/store/preauthored_task_01_test.go] [proves: req-create, req-persist]

## Task: task-02 — HTTP routes over the store
**Effort**: medium
**Depends-on**: [task-01]
**Frozen**: internal/store/preauthored_task_01_test.go

### Requirements

- req-routes: `GET /healthz` returns 200 with `ok`. `POST /api/links` returns 201 and the link, or 400 for a bad URL, slug or body, or 409 for a taken slug. `GET /api/links/{slug}` returns 200 and the link, or 404. `GET /{slug}` returns 302 with `Location` set to the URL, or 404.
- req-json: every `/api/links` response carries `Content-Type: application/json`, and every 4xx from it has the body `{"error": "<message>"}` with a non-empty message.
- req-list: `GET /api/links` returns 200 and a JSON array (`[]` when empty, never `null`) of every link, sorted by slug.
- req-persist-api: a link created through the API is still there after the data file is reopened.

### Description

Implement `server.New(dataPath string) (http.Handler, error)` in
`internal/server/server.go`. The stub currently panics with `todo: task-02`.
`New` opens the `internal/store` store at `dataPath` and returns its error
unchanged if the file is unreadable. It returns an `http.ServeMux` that uses
Go 1.22 method+pattern routes:

- `GET /healthz` → `200`, body `ok`.
- `POST /api/links` with body `{"url": "...", "slug": "optional"}` → `201` and the link as JSON. Map `store.ErrInvalidURL` and `store.ErrInvalidSlug` to `400`, `store.ErrSlugTaken` to `409`, and a body that is not JSON to `400`.
- `GET /api/links` → `200`, a JSON array of every link sorted by slug.
- `GET /api/links/{slug}` → `200` and the link, or `404`.
- `GET /{slug}` → `302` with `Location` set to the link's URL, or `404`.

Write compact JSON with `encoding/json` (`{"slug":"zurdo","url":"…"}`, no
indentation) and set `Content-Type: application/json` on every `/api/links`
response, errors included. Every 4xx under `/api/links` (400, 404, 409) has
the body `{"error": "<message>"}`, per the §4 convention phase 2's client
relies on. The `GET /{slug}` 404 may stay plain text.

The tests in `internal/server/preauthored_task_02_test.go` are evidence, not
scratch. Do not delete them, weaken their assertions, or change what they
assert. The only edit this task may make to that file is removing the
`//go:build preauthored` line (and the blank line after it) once the behavior
is implemented.

### Acceptance Criteria

- [ ] every route returns its specified status, body and Location under httptest, with JSON content type and `{"error"}` bodies on API errors [shell: go test -tags=preauthored -run '^TestServerRoutes$' -count=1 ./internal/server/] [proves: req-routes, req-json]
- [ ] the list is a sorted JSON array, empty when the data file is missing, and API-created links survive a reopen [shell: go test -tags=preauthored -run '^TestServer(List|CreatePersists)$' -count=1 ./internal/server/] [proves: req-list, req-persist-api]
- [ ] the pre-authored server tests are no longer build-tagged [no-grep: //go:build preauthored in internal/server/preauthored_task_02_test.go] [proves: req-routes, req-json, req-list, req-persist-api]

## Task: task-03 — serve command, seed data and live checks
**Effort**: low
**Depends-on**: [task-02]
**Frozen**: scripts/devserver.sh, internal/store/preauthored_task_01_test.go, internal/server/preauthored_task_02_test.go

### Requirements

- req-serve: `shortlink serve [--addr 127.0.0.1:8787] [--data links.json]` serves `server.New(data)` on `addr`. When the data file can't be read, it exits 1 with a message on stderr that names the path.
- req-seed: `testdata/seed.json` holds the two seed links (`zurdo` → `https://github.com/ElOrlis/zurdo`, `carol` → `https://github.com/ElOrlis/carol`), and the dev server serves them.
- req-redirect: the live dev server redirects `/zurdo` with a 302 to the zurdo repository.
- req-clean: `gofmt -l .`, `go vet ./...` and `go test ./...` are clean, with the pre-authored tests included.

### Description

Add a `serve` subcommand to `cmd/shortlink`. Use the `flag` package for its
flags, with `--addr` defaulting to `127.0.0.1:8787` and `--data` defaulting to
`links.json`. It calls `server.New(data)`. If that fails, it prints the error,
including the data path, to stderr and exits 1. Otherwise it runs
`http.ListenAndServe(addr, handler)`. Keep whatever `cmd/shortlink` already
does (e.g. version output), and keep `cmd/shortlink/main_test.go` passing.

Create `testdata/seed.json` in the store's file format (a JSON array of
`{"slug","url"}` objects) holding the two seed links listed under req-seed.

The live checks run against the dev server that `scripts/devserver.sh` serves
from the newest tree on `http://127.0.0.1:8787`. Do not start another server
on that port, and do not modify `scripts/devserver.sh`.

### Acceptance Criteria

- [ ] serve with an unreadable data file exits 1 and names the path [shell: bin="$(mktemp -d)/shortlink" && go build -o "$bin" ./cmd/shortlink && { out="$(timeout 10 "$bin" serve --addr 127.0.0.1:0 --data testdata 2>&1)"; code=$?; test "$code" -eq 1 && printf '%s' "$out" | grep -q testdata; }] [proves: req-serve]
- [ ] the live dev server is healthy [http: GET http://127.0.0.1:8787/healthz -> 200 contains "ok"] [proves: req-serve]
- [ ] the live list carries the zurdo seed link [http: GET http://127.0.0.1:8787/api/links -> 200 contains "\"slug\":\"zurdo\""] [proves: req-seed]
- [ ] the live server reads the carol seed link and 404s an unknown slug [http: GET http://127.0.0.1:8787/api/links/carol -> 200 contains "https://github.com/ElOrlis/carol"] [http: GET http://127.0.0.1:8787/api/links/nope -> 404] [proves: req-seed]
- [ ] the live server redirects /zurdo without following it [shell: curl -s -o /dev/null -w '%{http_code} %{redirect_url}\n' http://127.0.0.1:8787/zurdo | grep -qx '302 https://github.com/ElOrlis/zurdo'] [proves: req-redirect]
- [ ] the tree is gofmt-clean, vets clean and the full suite passes with pre-authored tests included [shell: test -z "$(gofmt -l .)" && go vet ./... && go vet -tags=preauthored ./... && go test -tags=preauthored -count=1 ./...] [proves: req-clean]
