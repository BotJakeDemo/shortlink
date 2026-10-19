---
name: zurdo-goal-author
description: "Authoring an initiative's scope from a Goal — the tier above design records: Objectives and Key Results and non-functional requirements that each name a threshold and the instrument that reads it, Epics cut into the phase rows of a `scope.md` that Carol reads, and the handoff of each Epic to a design record or a PRD."
allowed-tools: [Read, Write, Edit, Bash]
disable-model-invocation: true
---

## When to invoke

Invoke this skill when a body of work is still a **Goal** — an outcome somebody wants, with
no agreed shape. It produces one `scope.md`: the Destination stated from the outside, the
Objectives and Key Results that say when it is reached, the non-functional requirements it
must hold while reaching it, and the Epics it is cut into. It is the tier above
`zurdo-design-author`: Goal → OKRs → NFRs → Epics → a design record or a PRD per Epic →
`zurdo run`.

Reach for it when the honest answer to *"what are the phases?"* is a shrug, when two people
would describe "done" differently, or when the work is large enough that the first PRD cannot
be chosen without deciding what the last one is for.

Do not invoke it to adjudicate **how** something will be built on evidence — that is
`zurdo-design-author`, and step 5 hands off to it. Do not invoke it to turn a settled slice
into tasks and criteria — that is `zurdo-prd-author`. And do not invoke it on an initiative
that already has a `scope.md`: an existing scope is amended at a phase review, never
rewritten from a fresh interview.

## Where the method comes from

The method below is not invented here. It is extracted from this repository's own initiative
scopes under `docs/planning/*/scope.md`, which are the artifacts this skill writes:

- **A Destination stated from the outside.** `release-readiness` opens on what a developer
  can do — install through channels they already trust, update from a binary whose origin can
  be checked — not on the modules that would change.
- **Measured facts under `## Notes`, with the number attached.** That initiative records
  "8 of 273 Claude iterations undercount" against the file and symbol that undercounts. That
  single number is what justified cutting a token-accounting phase at all, and it is why
  nobody re-argued the phase afterwards.
- **Phases ordered by `Depends on`, not by enthusiasm.** The table's edges are what let the
  review reorder six phases in one line without re-cutting any of them.
- **Fog kept apart from decisions.** `## Not yet specified` carries the open questions
  verbatim — including ones found mid-run — so `## Decisions so far` stays a record of what
  was actually settled.

Those four properties are the whole method. Everything below is the interview that produces
them.

## The gate

`zurdo-prd-author` gates on `zurdo analyze --static-only`, and `zurdo-design-author` gates on
a measured number per claim. A scope's gate is **measurability**:

> Every Key Result and every NFR names a threshold and the instrument that reads it — a
> command, a test, or a zurdo hint type.
>
> **No instrument, no key result.**

An instrument is a thing that can be run by somebody who does not trust you. `cargo test
--test docs_truth_corpus` is an instrument. `zurdo validate --format json` is an instrument.
A zurdo hint type is an instrument too, which is the point of the rule: an NFR written as a
threshold plus the hint that reads it is already a criterion, and the PRD that inherits it
copies it rather than inventing it.

"Faster", "more reliable", "better onboarding" are not instruments and not thresholds.

### When a Key Result cannot name its instrument

It is not refused, and it is not left standing. **Demote it to `## Not yet specified`, naming
what would measure it** — the counter that does not exist yet, the corpus nobody has, the
run that has never happened. A demoted KR is a live question with an address, and the next
phase review can answer it. A KR left in the body unmarked is a claim that reads as though
it had been checked.

Do not resolve the fog by weakening the Objective until it sounds measurable. An Objective is
allowed to be qualitative; only its Key Results are not.

## The five steps

Each step has a gate. Do not open the next one until it clears.

### 1. Goal → `## Destination`

Pin the Goal in **one or two sentences**, stated from the outside: what somebody can do, see,
or trust once this is done. Name no module, no command that does not exist, and no approach.

The test is substitution. If the sentence survives replacing every implementation noun in it
with a different plausible one, it is a Destination. If it does not, it is a design decision
wearing a Goal's clothes, and it will quietly rule out the phases nobody has thought of yet.

This paragraph becomes the file's `## Destination`, verbatim.

**Gate:** the Goal names an outcome and no implementation, and the human has agreed to its
wording — not to its gist.

### 2. Objectives and Key Results

Each **Objective** is one qualitative outcome the Destination decomposes into. Two to four is
the useful range; an initiative with eight Objectives is two initiatives.

Each **Key Result** under it names a threshold and the instrument that reads it. Thresholds
are numbers or states, not directions: "every provider call's tokens appear in `prd.json`,
read by `zurdo report --format json`" is a Key Result; "better token accounting" is a mood.

Ask of every KR, in order: *what is the number?* *what reads it?* *what does it read today?*
The third question is the one that earns its keep — a KR whose instrument already reports the
target value is done, and belongs in Notes as a fact rather than in the plan as work.

They sit under `## Notes` as bold-labelled bullets, because that is what Carol renders today:

```markdown
- **Objective 1:** every provider call Zurdo makes is accounted for.
  - **KR 1.1:** 0 iterations in either corpus record a token count below the
    provider's own total — read by `zurdo report --format json`
  - **KR 1.2:** every in-run side call appears in `prd.json`'s `side_calls` —
    read by `cargo test --test tokens`
```

**Gate:** every KR carries a threshold and an instrument, or it has been demoted to
`## Not yet specified` naming what would measure it.

### 3. Non-functional requirements

Each NFR names a **quality attribute** — security, performance, portability, compatibility,
observability — a threshold, and the instrument. Same gate, same demotion, same bullet shape:

```markdown
- **NFR 1 (portability):** the install script runs on stock macOS and Debian with no
  package installs — read by `[shell: sh -n shared/dist/install.sh]`
- **NFR 2 (security):** no release is installed whose `checksums.txt` signature fails —
  read by `[shell: cargo test --test self_update]`
```

Writing the instrument as a zurdo hint is the cheapest thing in this skill. An NFR written
this way becomes an acceptance criterion later with no translation step, and the PRD author
inherits a line they would otherwise have had to invent under pressure.

NFRs are the part of a scope that is silently dropped. Ask explicitly what must stay true
*while* the Objectives are met — what must not get slower, less portable, or less verifiable
— and write the answer down even when it is "nothing new".

**Gate:** every NFR names an attribute, a threshold, and an instrument, or it is demoted.

### 4. Epics → the `## Phases` table

Each **Epic** is a slice that ships independently and **moves at least one Key Result**. An
Epic that moves no KR is either missing a KR or is not in this initiative; say which.

Each Epic becomes one row of `## Phases`: a `phase-NN` id, a Title, an **empty PRD cell**,
status `planned`, and its `Depends on` edges. The table always carries the `Depends on`
column, even when every cell in it is empty — the column is what later reviews reorder the
work with.

```markdown
| Phase | Title | PRD | Status | Depends on |
|-------|-------|-----|--------|------------|
| phase-01 | Every provider call's tokens, counted once and in full | | planned | |
| phase-02 | A verified install, and a binary that updates itself | | planned | |
| phase-03 | Signed releases, verified by every install channel | | planned | phase-02 |
```

Say the edges out loud and challenge each one: a `Depends on` that means "I would rather do
that first" is not a dependency, and it costs the initiative the right to reorder later.
Number the rows in dependency order where you can, so the first phase to run is `phase-01`.

**The vocabulary, stated once.** The interview says **Epic**; the file says **phase**; and a
phase becomes a **GitHub epic** only when its PRD is published. Three names, one thing, at
three stages — use the file's word once the file exists, and do not introduce "milestone".

**Gate:** every Epic ships independently, names the KR it moves, and has a row with its
`Depends on` edges.

### 5. Write the file, then hand off

Write **exactly these six `##` sections, in this order**:

| Section | What it carries |
| --- | --- |
| `## Destination` | step 1's paragraph, unedited |
| `## Notes` | the Objectives, KRs and NFRs, plus any fact measured during the interview, with its number and the file or command it came from |
| `## Decisions so far` | empty at authoring time — the first phase review fills it |
| `## Phases` | step 4's table |
| `## Not yet specified` | every demoted KR and NFR, and every open question, each naming what would answer it |
| `## Out of scope` | what was considered and deliberately excluded, each with its reason |

**Write no other `##` heading.** The reader this file is written for keeps these six and
drops anything else, so a seventh section is not a smaller section — it is a deleted one.
Write no `###` heading inside a kept section either: whether Carol keeps one
has not been verified, which is why the OKRs and NFRs are bold-labelled
bullets.

Then **refer the human onward, Epic by Epic** — the input rule
`docs/design/engineering-skills-verdict.md`'s Amendment 2026-10-10 fixed:

- to **`zurdo-design-author`** when the Epic must adjudicate approaches on evidence before
  any PRD can exist — two plausible designs, an unknown incumbent, a number nobody has taken.
  From that design record, the human goes on to **`zurdo scaffold`** (a draft the
  `zurdo-prd-author` skill then revises) or straight to **`zurdo-prd-author`**. `zurdo scaffold`
  never takes a Phases row — its input is the design record (or one `--section` of it), not
  this file;
- to **`zurdo-prd-author`** otherwise, when the approach is already settled in
  `## Decisions so far` or is obvious enough that the first honest question is "what are the
  tasks?".

All three are named, not called. A referral hands the next step to the person driving the
session, who decides when the scope is settled enough to be that command or skill's input
(`docs/design/skill-invocation-tiers.md` §5.4). Name them by name; they are peer skills and a
runtime command, not directories to read.

**Gate:** the file carries the six sections and no others, and every Epic has been referred
onward by name.

## Where it writes

The target is `<initiatives>/<initiative>/scope.md`, where `<initiatives>` is resolved like
this:

1. Read `.carol/config.toml` at the repository root. If it exists, take the `initiatives` key
   under its `[paths]` table.
2. If `.carol/config.toml` does not exist, use `docs/`.

`<initiative>` is a short kebab-case name for the Goal, not for its first phase:
`release-readiness`, `docs-truth`, `index-effectiveness`.

**If the target `scope.md` already exists, stop and tell the human.** Do not merge, do not
append, and do not write a second file beside it. An existing scope carries recorded
decisions and phase statuses that this interview has not seen, and amending an initiative is
the phase review's job — a different conversation, with the run state in front of it. Say
which file exists and what the review would do instead.

## References

- [references/okrs-and-nfrs.md](references/okrs-and-nfrs.md) — the labelled format in full: the fixed `**Goal:**` / `**Objective N:**` / `**KR N.M:**` / `**NFR N (<attribute>):**` labels, why they live inside `## Notes` rather than a heading of their own, the gate that ends every KR and NFR line with `read by <instrument>`, what counts as an instrument and what only sounds like one, and the demotion of an unmeasurable claim to `## Not yet specified`
- [references/example-scope.md](references/example-scope.md) — one complete `scope.md`, the file step 5 writes: a Destination stated from the outside, four Objectives with instrumented Key Results, four NFRs, seven phase rows with their `Depends on` edges, and demoted claims in fog. Read it before the interview and again before writing
- `zurdo-design-author` and `zurdo-prd-author` — where step 5 refers each Epic onward; call the Skill tool with the name (these are peer skills, named rather than addressed by path, and the two `references/` files above are this skill's own, which is why they carry a path)
