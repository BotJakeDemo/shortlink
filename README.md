# shortlink

A small Go URL shortener, built phase by phase by [Zurdo](https://github.com/ElOrlis/zurdo)
and driven on GitHub by [Carol](https://github.com/ElOrlis/carol).

- `shortlink serve` — the HTTP API (create, list, resolve and redirect links).
- `shortlink add | list | resolve | stats` — a CLI client for it.

Nothing here is written by hand past the skeleton: every feature lands through
an initiative phase. See [the scope](docs/shortlink/scope.md) and
[the design record](docs/shortlink/design/shortlink.md), or the **Shortlink**
Project board.

## How a phase moves

1. Carol **proposes** each open phase as an issue (written as the GitHub App).
2. A maintainer adds `carol:accepted` to a proposal — gate **H1**.
3. Carol **specifies** it: `zurdo scaffold` drafts the PRD from the phase's
   design-record section and opens a PRD pull request. A review session tightens
   it, and the maintainer approves — gate **H2**.
4. Carol **runs** it: `zurdo run` drives an agent through every task until each
   acceptance criterion passes, then opens the run pull request.
5. Carol **reports**: task issues are synced, an intent review posts its verdict,
   and the maintainer approves the merge — gate **H3**.

`carol watch docs/shortlink/scope.md --apply` takes every step between the gates.

## Running the demo locally

```bash
scripts/devserver.sh                               # keeps `shortlink serve` up on :8787 for [http:] criteria
carol watch docs/shortlink/scope.md --interval 1m --apply
carol status shortlink                             # where every phase stands
```

## Development

```bash
go test ./...
go run ./cmd/shortlink help
```
