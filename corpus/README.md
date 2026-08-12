# Corpus policy

`cases.json` is versioned as `readabilities-suite.corpus.v3`. The
`READ-001` through `READ-017` fixtures were written specifically for this suite
and are MIT-licensed with the repository; they are not copied web pages.

The selected Golden Oracle is the local Defuddle 0.19.2 CLI
(`defuddle parse ... --json --markdown`). Its input is the same static HTML
snapshot and base URL given to Rust; it neither fetches a second copy nor
transforms Rust output. The suite records native status, normalized content
facts, structure, metadata, and bounded difference evidence against that
baseline. Rust's final active-content sanitizer is a separate hard gate.

The cases cover noisy long-form articles, short content, Chinese text,
technical documentation, news, blogs, forum replies, MathML/footnotes, images,
active-content security, validation failure, Medium/Wikipedia/MDN layouts,
external site configuration, and a generic no-match fallback. Each oracle uses
stable required/forbidden text and Markdown facts rather than expected
serialized HTML bytes. Site cases also assert the Rust configuration provenance.

Real public pages are discovery inputs only. Full snapshots are temporary by
default and must not be committed without a separate redistribution review.
When a live issue becomes a regression case, reduce it to a self-authored or
minimal reproducible fixture with explicit required and forbidden facts. For
example, `READ-017` is a self-authored link-index reduction of Hacker News.
Defuddle treats it as readable content, so Rust must retain it as a successful
result while still sanitizing the output.
