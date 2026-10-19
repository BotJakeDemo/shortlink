# PRD: Phase 1 follow-up — link store and HTTP API

Follow-up to [the phase 1 PRD](prd-01-link-store-and-http-api.md) from its
intent review on PR #11 (run `prd-01-link-store-and-http-api-d68d`, commit
`a5763e0`). All three tasks landed as intended. task-03 is recorded `failed`
only because one criterion could not pass: zurdo 1.30.1 does not unescape `\"`
inside an `[http:] contains` needle, so
`contains "\"slug\":\"zurdo\""` searched for literal backslashes, though the
live body did contain `"slug":"zurdo"` (`lessons/lesson-9aedb6c4.md`).

This follow-up re-proves req-seed's list check with a hint that can pass. The
behavior was built in the phase 1 run, so the task is expected to pass at
pre-flight with no attempts; that is the intended shape here, not a vacuous pass.

## Task: task-01 — the live list carries the zurdo seed link
**Effort**: low
**Depends-on**: []
**Frozen**: scripts/devserver.sh, testdata/seed.json, internal/store/preauthored_task_01_test.go, internal/server/preauthored_task_02_test.go

### Requirements

- req-seed-list: `GET http://127.0.0.1:8787/api/links` on the dev server returns 200 and a compact JSON array that contains `"slug":"zurdo"`.

### Description

No code change is expected. The dev server (`scripts/devserver.sh`) serves the
newest tree on `127.0.0.1:8787` with a fresh copy of `testdata/seed.json`. Do
not start another server on that port, and do not modify the frozen files. If
the check fails, fix `internal/server` so `/api/links` writes compact JSON.

### Acceptance Criteria

- [ ] the live list carries the zurdo seed link [shell: curl -sf http://127.0.0.1:8787/api/links | grep -qF '"slug":"zurdo"'] [proves: req-seed-list]
- [ ] the live list answers 200 [http: GET http://127.0.0.1:8787/api/links -> 200 contains "https://github.com/ElOrlis/zurdo"] [proves: req-seed-list]
