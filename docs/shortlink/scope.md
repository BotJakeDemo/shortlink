# Scope: Shortlink

## Destination

Shortlink is a small Go URL shortener, built phase by phase by Zurdo and driven by Carol: `shortlink serve` is an HTTP API that stores and redirects links, and `shortlink add/list/resolve/stats` is a CLI client for it.

Done when every phase is `done`, each having gone from `carol:accepted` on its proposal to `carol:reported` on its epic, with humans acting only at the gates.

## Notes

- **Domain:** a URL shortener: links, slugs, redirects and click counts. One binary; `serve` is the API and every other subcommand is a client.
- **Design record:** [design/shortlink.md](design/shortlink.md) — one section per phase, which `carol specify` hands to `zurdo scaffold`.
- **Standing preferences:** Go 1.22 standard library only; `gofmt` and `go vet` clean; table-driven tests with `httptest`. `[http:]` criteria target the dev server on `127.0.0.1:8787` (`scripts/devserver.sh`), which must be running while a phase runs.

## Decisions so far

- Initiative scoped (maintainer, 2026-10-19). Three phases, each with its own section in [the design record](design/shortlink.md): the HTTP API first, the CLI client on top of it, then click stats across both. Every phase carries a `Design` cell so Carol can specify it without a hand-authored PRD.
- `[http:]` criteria run against a live dev server (2026-10-19). zurdo's PRD grammar has no setup step, so `scripts/devserver.sh` rebuilds and restarts `shortlink serve` on every `.go` change and follows the newest `.carol/worktrees/<phase>`, with fresh `testdata/seed.json` data on each restart.
- PRD-01 reviewed at gate H2 (maintainer, 2026-10-19, [trail](prds/prd-01-link-store-and-http-api.trail.md)). `carol tick` specified phase-01 into #5 once `carol:accepted` was on #2. The H2 session asked one question, whether `Link` keeps the design's `{slug, url, created_at, clicks}`, and `carol:take-recommended` answered it: `Link` is `{slug, url}` in phase 1, `clicks` and `last_clicked_at` arrive in phase 3, and `created_at` is dropped (design §4 amended). The review added `{"error"}` 4xx bodies, the slug rules and cross-task frozen test files.
- Phase 1 review (maintainer, 2026-10-19). Phase-01 is `done`. The run (`prd-01-link-store-and-http-api-d68d`, 2m 39s, $0.29) passed task-01 and task-02 at one attempt each. task-03 landed but was recorded `failed`: its `[http: … contains "\"slug\":\"zurdo\""]` needle can't pass, because zurdo 1.30.1 doesn't unescape `\"` (`lessons/lesson-9aedb6c4.md`). H3 found all three tasks landed. [The follow-up](prds/prd-01-link-store-and-http-api-followup.md) re-proved the list check with `curl | grep -qF` and passed at pre-flight. `TestRunServeBadData` (unattributed) is kept. #11 merged as `595c554`. Two tooling gaps: the first H3 session left its follow-up uncommitted, so Carol's follow-up run found no file and halted #11 (CLAUDE.md now tells review sessions to commit follow-ups). And `carol run` commits leftovers with `:(exclude).zurdo`, so run state is committed to the run branch by hand before H3.

## Phases

| Phase | Title | PRD | Status | Depends on | Design |
|-------|-------|-----|--------|------------|--------|
| phase-01 | Link store and HTTP API | docs/shortlink/prds/prd-01-link-store-and-http-api.md | done | | design/shortlink.md#Phase 1 — link store and HTTP API |
| phase-02 | CLI client | | ready | phase-01 | design/shortlink.md#Phase 2 — CLI client |
| phase-03 | Click stats | | ready | phase-02 | design/shortlink.md#Phase 3 — click stats |

## Not yet specified

- Whether a later phase adds a `[manual]` criterion to show `zurdo review --sign-off` (design record open question 1).
- Link expiry and deletion (`DELETE /api/links/{slug}`).

## Out of scope

- Authentication and users — they add no new kind of criterion to the demo.
- A database — a JSON file keeps every run hermetic.
- A web UI — `[http:]` checks status codes and bodies; HTML adds nothing to show.
