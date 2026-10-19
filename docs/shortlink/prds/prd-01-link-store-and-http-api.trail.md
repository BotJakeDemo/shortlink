> reviewed by the `zurdo-prd-author` skill (phase-01-h2, 2026-10-19) from the `zurdo scaffold` draft of `docs/shortlink/design/shortlink.md` §Phase 1 — verdict `✓ READY TO RUN`, the scaffold marker removed once the maintainer answered the `Link`-shape question (below).

# Trail: prd-01-link-store-and-http-api

## Decisions

- **Three tasks along the verification surface** (confirmed): the store (unit tests), the routes (httptest), and the command plus seed data (live dev server). Each task's evidence depends only on what the earlier tasks produce.
- **task-02 depends on task-01** because `server.New` opens the store. **task-03 depends on task-02** because `serve` calls `server.New`. Both edges are real data dependencies (confirmed).
- **Store API pinned by stub** (`Open`/`Create`/`Get`/`List`, three sentinel errors), so the pre-authored tests name symbols and still compile at the PRD commit (confirmed).
- **`server.New(dataPath string) (http.Handler, error)`** takes a path, so the server stub compiles without importing `internal/store`. The cost is that `server` owns opening the store (confirmed).
- **Slug rules are spelled out and asserted.** The draft deferred to "the §4 pattern" and only checked that a generated slug was non-empty, so a 1-character generator or a pattern without length bounds would have passed. The review added the literal `^[a-z0-9-]{3,32}$` and the 6-character generated-slug rule to req-create and the Description, and added test cases: a hyphenated slug (accepted), `ab`, a 33-character slug and `Zurdo` (all rejected), and a `^[a-z0-9]{6}$` match on the generated slug.
- **New req-json: the §4 error and content-type convention.** Design §4 makes `{"error": "<message>"}` with a 4xx status and `Content-Type: application/json` conventions that every phase keeps, and phase 2's client prints the server's error. The draft asserted neither. `TestServerRoutes` now checks the JSON content type on every `/api/` case and a non-empty `error` field on every `/api/` 4xx. The `GET /{slug}` 404 is left unconstrained.
- **Compact JSON is required** (confirmed). The live criterion matches `"slug":"zurdo"` byte-exactly.
- **Location is asserted only on 302s** (confirmed). A 201 that also sets `Location` is legitimate REST.
- **Two hints on one criterion for 404-on-nope** (confirmed). `GET /api/links/nope -> 404` alone passes against a server that 404s everything; AND'ing it with the carol 200 removes the tautology.
- **The unreadable-file shell check greps for the path** (confirmed). Today's `main` exits 2 on `serve`, and its message does not contain `testdata`. `timeout 10` turns a server that keeps listening into a failure instead of a hang.
- **The redirect check is a pipe into `grep -qx`.** `zurdo analyze` flagged the draft's `test "$(curl …)" = "…"` as a no-op. The rewritten check, `curl … -w '%{http_code} %{redirect_url}\n' | grep -qx '302 https://github.com/ElOrlis/zurdo'`, clears the lint, and it exits 1 against the dead port today.
- **The full-suite criterion runs with `-tags=preauthored`** (confirmed). The untagged trio already passes on the PRD commit (tautology).
- **Frozen paths:** task-01 freezes task-02's test file, task-02 freezes task-01's, and task-03 freezes both plus `scripts/devserver.sh`. None of these overlaps a same-task evidence path, and per-task baselines (v1.13.1+; this repo runs zurdo 1.30.1) keep the earlier tasks' writes from counting. Advisory: `scripts/devserver.sh` must never change in any phase. Add it to `[verification] protected_paths` in `.zurdo/config.toml` rather than freezing it per task in later PRDs.
- **Accepted and unverified: the mutex and the temp-file-and-rename write.** No phase-1 criterion can observe a lost write or a torn file cheaply. Phase 3's 100-concurrent-redirects test under `-race` is where concurrency is proven.
- **Effort keys** `medium`/`medium`/`low` (confirmed). `.zurdo/config.toml` maps all three keys for every provider.

## Lessons written

- `lessons/lesson-3634089f.md` (task-03): `zurdo analyze` flags `test "$(…)" = "…"` as a no-op shell hint, and `… | grep -qx '…'` replaces it. This is the only correction from this review that transfers to other PRDs. The other decisions above are scope calls and stay in this trail.

## Resolved question

- **`Link` shape vs design §4** (resolved 2026-10-19, BotJakeDemo took the recommendation on PR #5). The PRD pins `Link` to exactly `{slug, url}`, which the original §4 contradicted with `{slug, url, created_at, clicks}`. Decision: keep `{slug, url}` in phase 1. Design §4 is amended in the same commit so that fields arrive with the phase that observes them: phase 3's Surface now adds `clicks` and `last_clicked_at`, and `created_at` is dropped because no phase reads or checks it. Every phase-1 criterion stands as drafted, and `TestStoreReload` keeps its `reflect.DeepEqual` comparison.

## Pre-authored tests

- **task-01** — `internal/store/preauthored_task_01_test.go`, 4 tests. Negative cases: relative, `ftp` and empty URLs; slugs with a space or a slash, too short, too long or uppercase; a taken slug; a directory or malformed file on `Open`. Pre-run: `go test -tags=preauthored -run '^TestStoreCreate$' -count=1 ./internal/store/` → **failed on assertion** (`todo: task-01` panic). Satisfiability was checked against a throwaway reference store, which passed all 4 tests; a 4-character generator failed `TestStoreCreate`.
- **task-02** — `internal/server/preauthored_task_02_test.go`, 3 tests. Negative cases: unknown slug on both read routes, a non-URL, an `ftp` URL, a taken slug, a bad slug, a malformed body, and plain-text error bodies. Pre-run: `go test -tags=preauthored -run '^TestServerRoutes$' -count=1 ./internal/server/` → **failed on assertion** (`todo: task-02` panic). The throwaway reference server passed all 3 tests; swapping its JSON errors for `http.Error` failed 12 assertions.
- **task-03** — **declined.** Its behavior is the `serve` process and the live dev server, which the `[shell:]` and `[http:]` criteria exercise directly. A Go test would duplicate the store/server tests.
