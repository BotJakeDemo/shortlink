# OKRs and NFRs: the labelled format, and the gate every line passes

A scope's Goal, Objectives, Key Results and NFRs are not a document of their
own. They sit inside **one** `## Notes` heading in `scope.md`, as bold-labelled
bullets.

That is not a style preference. The reader this file is written for — Carol —
keeps exactly six `##` sections and drops any other heading it finds:

| Kept section | What it carries |
| --- | --- |
| `## Destination` | the Goal stated from the outside, in one or two sentences |
| `## Notes` | the Goal line, the Objectives, the KRs, the NFRs, and every fact measured during the interview |
| `## Decisions so far` | what has actually been settled — empty at authoring time |
| `## Phases` | the Epic table, with its `Depends on` column |
| `## Not yet specified` | demoted KRs and NFRs, and open questions, each naming what would answer it |
| `## Out of scope` | what was considered and deliberately excluded, each with its reason |

A seventh heading is not a smaller section — it is a deleted one. An `## OKRs`
heading takes every Objective under it down with it, silently, and the scope
then reads as though nobody had ever said when it was finished. Labels inside
`## Notes` survive. That is the entire reason the format looks like this.

## The labels

Write these four verbatim. They are fixed so that a later Carol change can
render OKRs properly by parsing what is already in the file, with no migration
of the scopes written before it:

| Label | How many | Nesting |
| --- | --- | --- |
| `**Goal:**` | exactly one per scope | a top-level bullet, first under `## Notes` |
| `**Objective N:**` | two to four | a top-level bullet |
| `**KR N.M:**` | one or more per Objective | nested one level under Objective `N`, and `M` restarts at 1 under each |
| `**NFR N (<quality attribute>):**` | as many as hold | a top-level bullet, numbered in its own sequence |

`**KR 1.1:**` — not `**KR1.1**`, not `*KR 1.1:*`, not `**KR 1.1**` without the
colon. The bold markers and the colon are part of the label.

The `**Goal:**` line restates `## Destination` in one sentence. The Destination
is the agreed wording and is never edited to match the Goal line; the Goal line
is the compressed version a reader skimming `## Notes` needs, and it is what
every Objective is checked against.

`<quality attribute>` is a real attribute — security, performance, portability,
compatibility, observability, accessibility — in parentheses, lowercase. An NFR
whose attribute is "quality" or "robustness" has not been written yet.

NFRs are **not** nested under an Objective. They must hold while every Objective
is met, so they carry their own flat numbering.

### The shape, in full

```markdown
## Notes

- **Goal:** developers install and update Zurdo through channels they already
  trust, and every provider call it makes is accounted for.
- **Objective 1:** every provider call Zurdo makes is accounted for.
  - **KR 1.1:** 0 iterations in either corpus record a token count below the
    provider's own total — read by `zurdo report --format json`
  - **KR 1.2:** every in-run side call appears in `prd.json`'s `side_calls` —
    read by `cargo test --test tokens`
- **Objective 2:** an install's origin can be checked by somebody who does not
  trust the publisher.
  - **KR 2.1:** every release carries a signature over `checksums.txt`, and an
    unsigned release is refused — read by `[shell: cargo test --test self_update]`
- **NFR 1 (portability):** the install script verifies a release on stock macOS
  and Debian with no package installs — read by `[shell: sh -n dist/install.sh]`
- **NFR 2 (observability):** a refused verification names the file that failed —
  read by `[grep: signature verification failed in dist/install.sh]`
```

Facts measured during the interview go under `## Notes` too, as ordinary
bullets beside the labelled ones, each carrying its number and where it came
from. A measured fact is what stops a phase from being re-argued at every
review.

## The gate: every KR and NFR line ends with its instrument

Each KR and each NFR line ends with an em-dash (U+2014), the words `read by`,
and the instrument in backticks. Nothing follows the closing backtick.

> **No instrument, no key result.**

An **instrument** is a command, a test, or a zurdo hint type — a thing somebody
who does not trust you can run, and get the same answer you got.

| Names no instrument | Names one |
| --- | --- |
| "users are happier" | "`carol status` parses it" |
| "onboarding is smoother" | "`zurdo validate --format json` reports 0 errors on the authored PRD" |
| "the install is more trustworthy" | "`[shell: cargo test --test self_update]`" |
| "token counts are more accurate" | "`zurdo report --format json` records a per-call total equal to the provider's own" |
| "the docs are better" | "`[grep: --sign-off in README.md]`" |

The left column is not a worse way of saying the right column. It is a
different kind of sentence: nothing can disagree with it, so nothing can
confirm it either. "Users are happier" has no reading that comes back false, and
a plan whose finish line cannot come back false has no finish line.

A threshold is a number or a state, not a direction. "Fewer undercounts" is a
direction; "0 undercounting iterations" is a threshold.

### Why a hint type counts as an instrument

`[shell:]`, `[grep:]`, `[file-exists:]` and the rest are instruments, and
writing one is the cheapest line in this skill. A KR or NFR written as a
threshold plus the hint that reads it **is already an acceptance criterion**.
The PRD author inherits it and copies it across; without it they invent a
criterion under pressure, at the moment they are least able to argue about what
"done" meant.

Ask of every KR, in order: *what is the number?* *what reads it?* *what does it
read today?* The third question earns its keep — a KR whose instrument already
reports the target value is not work. It is a fact, and it belongs under
`## Notes` as one.

## The demotion: a claim with no instrument goes to fog

A KR or NFR that cannot name an instrument is **not refused and not left
standing**. Move it to `## Not yet specified`, naming what would measure it:

```markdown
## Not yet specified

- Whether a new user reaches a first passing run without reading the spec.
  Demoted from KR 3.2: no instrument exists. It would need a timed
  first-session corpus, which nobody has taken.
- Copilot's token usage: whether `outputTokens` is per message, and where
  input tokens could come from. Demoted from KR 1.3: no Copilot run exists in
  either corpus to read.
```

Each demoted line carries three things: the claim, the label it was demoted
from, and **the instrument that does not exist yet** — the counter nobody
built, the corpus nobody has taken, the run that has never happened. That third
part is what makes it a live question with an address, answerable at the next
phase review, instead of a regret.

This is the same rule `zurdo-design-author` applies to a claim it cannot
measure: a design record does not drop such a claim and does not assert it
either; it records it as unmeasured and says what would measure it. One
discipline, one tier apart.

Two ways of failing the gate without noticing:

- **Leaving it in the body unmarked.** A KR in `## Notes` with no instrument
  reads to every later reader as though it had been checked. That is the only
  outcome the gate exists to prevent.
- **Weakening the Objective until it sounds measurable.** An Objective is
  allowed to be qualitative — that is what an Objective is for. Only its Key
  Results are not. Trading away the outcome to buy a tidy number loses the
  initiative its reason for existing and keeps none of the rigour.
