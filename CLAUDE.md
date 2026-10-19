# shortlink

Go 1.22, standard library only. One binary, `cmd/shortlink`; packages live under
`internal/` as laid out in `docs/shortlink/design/shortlink.md` §4.

- Run `gofmt -l .`, `go vet ./...` and `go test ./...` before finishing; all must be clean.
- Tests are table-driven; HTTP handlers are tested with `net/http/httptest`.
- `[http:]` acceptance criteria hit `http://127.0.0.1:8787`, served by
  `scripts/devserver.sh` from the newest tree. Don't start another server on that port.
- Seed data for the dev server lives in `testdata/seed.json`.
- `docs/shortlink/scope.md` is the source of truth for the initiative; GitHub issues are projections of it.
