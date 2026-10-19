# Scope: Release readiness — self-update, a verified install, true token counts, and Goal-to-PRD authoring

## Destination

Zurdo reaches developers through install channels they already trust —
Homebrew and a verified install script — and updates itself from a public
binary whose origin can be checked; reports the real token usage of every
provider call it makes; lets an agent harness drive `zurdo review` without a
terminal; takes a project from Goal → OKRs → NFRs → Epics to a scaffolded
first-draft PRD; and lets developers sign in from the terminal, free, to put
their name on what they sign off.

## Notes

- **Goal:** a developer installs Zurdo from a channel they already trust,
  updates it from a binary whose origin they can check, and gets true token
  numbers out of every run — and a project reaches its first PRD without its
  author inventing the format on the way.
- **Objective 1:** every provider call Zurdo makes is accounted for.
  - **KR 1.1:** 0 Claude iterations record a token total below the provider's
    own (baseline: 8 of 273 undercount today, by about 26.96M input and 172k
    output tokens, because only the last `"type":"result"` event is read and a
    background-task stream carries several) — read by `zurdo report --format json`
  - **KR 1.2:** every in-run side call (probe, diagnosis, post-mortem, lesson
    extraction) appears in the run record's totals; today none of them do —
    read by `cargo test --test tokens`
  - **KR 1.3:** a transient attempt's tokens survive its retry, and each retry
    keeps its own transcript; today the retry overwrites the attempt's
    transcript and the usage cannot be recovered (5 transient retries in the
    second corpus) — read by `cargo test --test tokens`
  - **KR 1.4:** every iteration carries a cost estimate (baseline: null on 438
    of 491, because the pricing table has a row for one model only) — read by
    `zurdo report --format json`
- **Objective 2:** an install's origin can be checked by somebody who does not
  trust the publisher.
  - **KR 2.1:** `zurdo self-update` moves a released binary to the latest
    version, pins an older one with `--to`, and refuses a release whose
    signature fails — read by `cargo test --test self_update`
  - **KR 2.2:** `curl … | sh` installs a verified release on a machine with no
    Rust toolchain, with completions and man pages — read by
    `cargo test --test install_sh`
  - **KR 2.3:** a binary installed under a Homebrew prefix refuses to replace
    itself and names `brew upgrade zurdo` instead — read by
    `[shell: cargo test --test self_update brew]`
- **Objective 3:** an agent harness drives a review to sign-off without a
  terminal, and a sign-off carries a name somebody verified.
  - **KR 3.1:** `zurdo review --no-prompt` prints every pending `[manual]`
    criterion with its evidence as one JSON object and exits 0 even when
    nothing is pending — read by `cargo test --test review_headless`
  - **KR 3.2:** a single criterion can be signed off by address, no command
    approves all of them at once, and the log records which channel signed —
    read by `[shell: cargo test --test review_headless]`
  - **KR 3.3:** a signed-in developer's verified name replaces a self-asserted
    one on a sign-off, and `zurdo auth status` reports who is signed in — read
    by `cargo test --test auth`
- **Objective 4:** a Goal reaches a first-draft PRD without its author
  inventing the format on the way.
  - **KR 4.1:** a scope authored from a Goal parses with every phase row and
    every dependency edge intact — read by
    `[shell: carol status <initiative> --format json]`
  - **KR 4.2:** a first-draft PRD scaffolded from a design document validates
    with 0 errors — read by `zurdo validate --format json`
- **NFR 1 (portability):** `install.sh` verifies a release with stock
  `openssl`, macOS LibreSSL included, and installs no package to do it — read
  by
  `[shell: openssl dgst -sha256 -verify release-signing.pub.pem -signature checksums.txt.sig checksums.txt]`
- **NFR 2 (security):** no release is installed whose `checksums.txt` fails
  verification, through any channel and with no override — read by
  `cargo test --test self_update`
- **NFR 3 (compatibility):** the run loop works signed out, offline, and in
  CI; signing in is optional and adds identity only — read by
  `cargo test --test auth`
- **NFR 4 (compatibility):** every document this initiative authors keeps the
  six `##` sections its reader renders and adds no seventh — read by
  `[shell: carol status <initiative> --format json]`
- **Measured 2026-10-07** over this repository's and a second project's run
  state: the Claude undercount above (8 of 273 iterations), one example
  iteration recording 99,479 in / 33 out against a true 4,494,906 / 20,067,
  the null cost estimates (438 of 491), the analyzer, judge, heal and
  lesson-extraction calls that never reach a run total, and 3 of 218 Codex
  iterations whose failed turns carry no usage at all — the provider's limit,
  not a parse bug.
- **Self-update has code to reuse.** The Pyxis model installer already
  downloads from the public release repository and checks sha256 with the
  `ureq` and `sha2` dependencies. Release tarballs are named
  `zurdo-v<version>-<target>.tar.gz`, and `checksums.txt` covers them.
- **The reader, read 2026-10-07** (`carol` v0.1.4): it reads `scope.md`,
  ticket files, PRDs and handoffs only. It keeps the six `scope.md` sections
  this file uses and drops any other `##` heading. It has no notion of goals,
  OKRs, NFRs or epic documents, so rendering them properly is a change in its
  own repository, not here.
- homebrew-core rejects formulae without an open-source license, so Zurdo,
  under its EULA, stays in its own tap.
- Opened 2026-10-07 as five urgent must-haves, ahead of a parked
  index-effectiveness initiative whose first comparison reads attempt-1
  outcomes from the run record — and the token fields beside them are wrong
  today, which is why Objective 1 runs first.

## Decisions so far

## Phases

| Phase | Title | PRD | Status | Depends on |
|-------|-------|-----|--------|------------|
| phase-01 | Every provider call's tokens, counted once and in full | | planned | |
| phase-02 | `zurdo self-update` and a verified install script | | planned | |
| phase-03 | Signed releases, verified by every install channel | | planned | phase-02 |
| phase-04 | Goal → OKRs → NFRs → Epics, as a skill Carol can read | | planned | phase-07 |
| phase-05 | `zurdo scaffold` — a first-draft PRD from a design document | | planned | phase-04 |
| phase-06 | Terminal sign-in | | planned | phase-05 |
| phase-07 | Headless `zurdo review` for agent harnesses | | planned | phase-03 |

## Not yet specified

- Whether Copilot's `outputTokens` is per message, and where its input tokens
  could come from. Demoted from KR 1.5: no instrument exists, because no
  Copilot run exists in either corpus to read.
- Whether a developer who has never used Zurdo reaches a first passing run
  without reading the spec. Demoted from KR 4.3: no instrument exists. It
  would need a timed first-session corpus, which nobody has taken.
- Whether a first install can be something other than trust-on-first-use, and
  what the procedure is if the signing key is compromised.
- Whether the scaffold command takes one phase row as input as well as a
  design document.
- Which OAuth flow and which identity source terminal sign-in uses: the
  planned identity provider offers no native device-authorization flow, and
  GitHub is not a standard OIDC provider for it.
- Whether a run whose transient-retry budget is exhausted should record its
  invocations; today it writes no iteration, so their tokens are lost.

## Out of scope

- Retiring the Homebrew tap — Homebrew stays a supported channel because
  developers already trust it.
- Changes to the reader itself — they live in its own repository; this
  initiative only produces documents it already reads.
- Intel macOS and Windows binaries — the destination is a verified install of
  the targets already built.
- Other package managers (apt, `cargo install`, npm) — the source is closed,
  and the destination names one install path.
