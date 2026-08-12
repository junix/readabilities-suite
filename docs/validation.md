# Validation evidence

Date: 2026-08-11

## Selected Oracle: Defuddle

The local `/Users/junix/defuddle` checkout was probed through its public CLI:
`node dist/cli.js parse HTML --json --markdown` (Defuddle 0.19.2). The suite
passes one static HTML snapshot to that CLI and to Rust; neither participant
downloads a second copy. The adapter only normalizes Defuddle's JSON envelope
and converts its explicit “No content could be extracted” diagnostic into the
shared `no_content` status. All other nonzero exits remain failures.

`record-golden --force` freezes Defuddle's status, Markdown, metadata, fixture
SHA-256, and content SHA-256 for every reviewed self-authored fixture. Normal
runs validate those digests. Explicit recording intentionally loads fixture
definitions without the old Golden body so an approved fixture change can
replace a stale Golden; it is the only path that bypasses the stale-digest
check.

## Oracle contract and safety boundary

The Golden asserts Defuddle's extraction status and at least 90% normalized
token recall by Rust. Lines marked as fixture noise are excluded before this
recall calculation, so Rust may safely remove declared noise. Required text,
forbidden text, structure, Markdown, ordering, metadata, and Rust site-config
facts remain independent checks.

Defuddle is an extraction Oracle, not a security authority. On the captured
MDN and W3C live snapshots it retained iframe markup, which the suite reports
as two Oracle security observations. Rust's final sanitizer remains a hard
gate and returned zero security violations. Oracle residue is visible but does
not make a safe Rust response fail merely for omitting unsafe markup.

## Defuddle-driven Rust refinement

Defuddle considers the Hacker News front page and the self-authored `READ-017`
reduction readable link-index content. Rust had converted this shape to typed
`no_content` using a global link-density heuristic. The heuristic was removed:
Rust now accepts the cleaned, non-empty extracted candidate, preserves list
content, and still applies the mandatory sanitizer. `READ-017` now requires
successful extraction and has a Rust regression test for its first and last
entries.

## Results

The refreshed offline corpus has 17 Defuddle Golden artifacts. Rust passed all
17 direct cases, all 17 status/Golden comparisons, and the zero-security gate.
Defuddle's own direct-score failures remain reported because the suite's richer
structure and security assertions are deliberately independent of Oracle
recall; they do not weaken the Rust release gate.

The same-snapshot live run passed all 14 Defuddle/Rust observations:

| Page | Rust status | Defuddle-token coverage |
|---|---|---:|
| English Wikipedia | success | 99.6% |
| Chinese Wikipedia | success | 99.8% |
| MDN article reference | success | 94.0% |
| NASA news release | success | 100.0% |
| Rust Blog | success | 100.0% |
| Hacker News front page | success | 100.0% |
| W3C ARIA accordion | success | 94.6% |

Live output is diagnostic only: public pages can drift, and no third-party
page body is committed. A future mismatch must be reduced to a self-authored
fixture, reviewed, and re-frozen before it becomes a release gate.
