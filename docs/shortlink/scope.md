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

## Phases

| Phase | Title | PRD | Status | Depends on | Design |
|-------|-------|-----|--------|------------|--------|
| phase-01 | Link store and HTTP API | docs/shortlink/prds/prd-01-link-store-and-http-api.md | running | | design/shortlink.md#Phase 1 — link store and HTTP API |
| phase-02 | CLI client | | ready | phase-01 | design/shortlink.md#Phase 2 — CLI client |
| phase-03 | Click stats | | ready | phase-02 | design/shortlink.md#Phase 3 — click stats |

## Not yet specified

- Whether a later phase adds a `[manual]` criterion to show `zurdo review --sign-off` (design record open question 1).
- Link expiry and deletion (`DELETE /api/links/{slug}`).

## Out of scope

- Authentication and users — they add no new kind of criterion to the demo.
- A database — a JSON file keeps every run hermetic.
- A web UI — `[http:]` checks status codes and bodies; HTML adds nothing to show.
