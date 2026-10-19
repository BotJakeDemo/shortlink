> drafted by `zurdo scaffold` from `docs/shortlink/design/shortlink.md` §Phase 1 — link store and HTTP API

# Trail: prd-01-link-store-and-http-api

## Decisions

- **Split into three tasks** along the verification surface: the store (unit tests), the routes (httptest), and the command plus seed data (live dev server). Each task's evidence depends only on what the earlier tasks produce.
- **task-02 depends on task-01** because `server.New` opens the store. **task-03 depends on task-02** because `serve` calls `server.New`. Both edges are real data dependencies.
- **Store API pinned by stub** (`Open`/`Create`/`Get`/`List`, three sentinel errors). The pre-authored tests can then name symbols and still compile at the PRD commit.
- **`server.New(dataPath string) (http.Handler, error)`** takes a path, not a `*store.Store`. The server stub then doesn't need to import `internal/store`, so the stubs compile without depending on the module path in `go.mod`. The cost is that `server` owns opening the store.
- **Data-file format is a JSON array of `{slug, url}`.** The server test seeds it directly, and `testdata/seed.json` uses the same shape. `Link` stays exactly two fields: `TestStoreReload` uses `reflect.DeepEqual` across a JSON round trip, which an added timestamp would break.
- **Compact JSON is required.** The live criterion matches `"slug":"zurdo"` byte-exactly, so the task-02 Description forbids indentation.
- **Location is asserted only on 302s.** A 201 that also sets `Location` is legitimate REST and should not fail.
- **Two hints on one criterion for 404-on-nope.** `GET /api/links/nope -> 404` alone could pass against any server that 404s everything. AND'ing it with the carol 200 removes that tautology.
- **The unreadable-file shell check also greps for the path.** Exit 1 alone could pass pre-run if the current `main` exits 1 on an unknown subcommand. The `timeout 10` turns a server that keeps listening into a failure instead of a hang.
- **The full-suite criterion runs with `-tags=preauthored`.** The untagged `gofmt`/`vet`/`test` trio already passes on the PRD commit (tautology). Including the tagged files makes it red until both pre-authored suites pass.
- **Frozen paths:** task-02 freezes task-01's test file, and task-03 freezes task-02's test file and `scripts/devserver.sh`. None of these overlaps a same-task evidence path. Per-task baselines (v1.13.1+) keep the earlier tasks' writes from counting as violations.
- **Effort keys** `medium`/`medium`/`low` assume the shipped `low|medium|high` defaults. Check them against `[effort_map.<provider>]` in `.zurdo/config.toml` before `zurdo run`.

## Open assumptions to check before committing

- `docs/shortlink/design/shortlink.md` §4 names the packages `internal/store` and `internal/server`. If it names them differently, rename the stubs, the tests and the hints together.
- The §4 slug pattern accepts the lowercase 4–9 character slugs the tests use (`zurdo`, `carol`, `bravo`, `plainhttp`, `taken`, `newone`, `kept`, `badurl`, `ftpfile`, `relative`, `emptyurl`) and rejects `has space` and `a/b`.
- `scripts/devserver.sh` runs `shortlink serve` with `--data` pointing at `testdata/seed.json` (or a copy of it) on `127.0.0.1:8787`.
- `zurdo analyze --static-only docs/shortlink/prds/prd-01-link-store-and-http-api.md` has not been run in this authoring pass. Run it and clear any `[error]` findings.

## Pre-authored tests

- **task-01** — `internal/store/preauthored_task_01_test.go`, 4 tests. Pre-run: `go test -tags=preauthored -run '^TestStoreCreate$' -count=1 ./internal/store/` → **failed on assertion**
- **task-02** — `internal/server/preauthored_task_02_test.go`, 3 tests. Pre-run: `go test -tags=preauthored -run '^TestServerRoutes$' -count=1 ./internal/server/` → **failed on assertion**
- **task-03** — **declined.** No pre-authored test in this draft.
