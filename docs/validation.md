# Validation evidence

Date: 2026-08-11

## Contract gate before adapters

The shared contract was first proved directly on `READ-001` through each native
Markdown surface. All three preserved the required central idea, evidence, and
footnote while rejecting related-content, newsletter, and footer facts. Rust
additionally retained the Rust code language, resolved relative links/images,
preserved the image-to-prose boundary and superscript, and passed its final
sanitizer.

The empty-input path remained observably different rather than being normalized
away: Defuddle returned a nonzero internal parse error, Python returned its
native null/no-content projection, and Rust returned a nonzero typed
`invalid_input`/`validate` error. The Python driver preserves this as
`no_content`; it does not convert it into success.

## Frozen-corpus comparison

Command:

```sh
./readabilities-suite run --profile offline \
  --report reports/offline.json \
  --markdown-report reports/offline.md
```

Rust passed 16/16 and met the bounded best-evidence gate:

| Participant | Cases | Required | Noise | Structure | Markdown | Order | Metadata | Security |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| Defuddle | 4/16 | 43/43 | 34/38 | 23/33 | 67/70 | 14/14 | 19/19 | 1 |
| readabilities-py | 3/16 | 43/43 | 32/38 | 22/33 | 60/70 | 14/14 | 19/19 | 2 |
| readabilities-rs | 16/16 | 43/43 | 38/38 | 33/33 | 70/70 | 14/14 | 19/19 | 0 |

The added cases prove actual Markdown semantics (code language, absolute URLs,
image spacing, math, footnotes), Medium/Wikipedia/MDN built-ins, an external
configuration, and a useful no-config fallback. Defuddle and Python retain
declared site noise in some cases; Python also loses code-language and math
semantics. In the active-content case, Defuddle retains an iframe and Python
retains script/dangerous-Markdown-URI evidence; Rust retains none.

## Deliberate mutation

A temporary fake Rust executable returned only:

```html
<article><h1>Broken extractor</h1><p>This output intentionally omits every required fact.</p></article>
```

Running only `READ-001` exited 1 and reported:

```text
READ-001/readabilities-rs:
failed_facts=[REQUIRED-CENTRAL-IDEA REQUIRED-EVIDENCE REQUIRED-STRUCTURE
REQUIRED-FOOTNOTE code table image footnote title author]
actual_preview="Broken extractor This output intentionally omits every required fact."
```

The fake executable was external to the repository. The real participant was
not changed; the complete baseline was rerun afterward.

## Live common-snapshot smoke

`just live-common` captures eight public page types once per run: Example Domain (short),
English Wikipedia (long), Chinese Wikipedia, MDN technical documentation, a
NASA news release, the Rust Blog, Hacker News (list/forum), and the W3C ARIA
accordion example. See `reports/live-2026-08-11.json` for URL, digest, output
size, bounded preview, status, security evidence, and timings.

Final smoke result:

| Participant | Native success | Suite pass | Security violations |
|---|---:|---:|---:|
| Defuddle | 8/8 | 6/8 | 2 iframe residues |
| readabilities-py | 8/8 | 8/8 | 0 |
| readabilities-rs | 8/8 | 8/8 | 0 |

The first run found a false positive in Rust's post-sanitizer invariant: literal
prose beginning `JavaScript:` was mistaken for a `javascript:` URI. The check
was narrowed to URI-bearing HTML tags/attributes, a regression test was added,
and the full live comparison was rerun. On the final run, Rust's Wikipedia and
MDN configurations matched automatically; all other pages used the generic
fallback. This live result is discovery evidence, not the frozen best-evidence
gate, because there is no page-specific required/forbidden oracle.
