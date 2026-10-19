---
schema: 1
match:
  hint_types:
    - http
source:
  kind: IntentReview
  prd_slug: prd-01-link-store-and-http-api
  task_id: task-03
  trail_path: docs/shortlink/prds/prd-01-link-store-and-http-api.trail.md
created_at: 2026-10-19T00:00:00Z
---

zurdo 1.30.1 does not unescape `\"` inside an `[http: … contains "<needle>"]` clause: `contains "\"slug\":\"zurdo\""` searches for the literal text `\"slug\":\"zurdo\"`, backslashes included, so it can never match a JSON body. To assert quoted JSON, use `[shell: curl -sf <url> | grep -qF '"slug":"zurdo"']`; keep `[http:]` needles free of double quotes (e.g. `contains "https://github.com/ElOrlis/zurdo"`). The `contains "\"token\""` example in the zurdo-prd-author criteria reference does not hold on 1.30.1.
