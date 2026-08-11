# readabilities-suite

Independent black-box comparison harness for three readable-content tools:

1. Defuddle (`defuddle parse ... --json --markdown`)
2. `readabilities-py` (public `readability_of` API through a thin native driver)
3. `readabilities-rs` (`readabilities-rs extract ... --format markdown`, with
   a separate JSON metadata/provenance projection)

The suite owns the corpus and fact oracle. Participant adapters only translate
input/output shapes; they do not extract content, delete noise, sanitize HTML,
or fabricate unsupported capabilities.

## Quick start

```sh
just build
./readabilities-suite doctor
./readabilities-suite list
./readabilities-suite run --profile offline
```

By default, the suite resolves sibling checkouts under
`~/projects/document-tools`, Defuddle under `~/defuddle`, and the Rust debug
binary under `../readabilities-rs/target/debug`. Override them with global
`--defuddle`, `--python`, or `--rust` flags.

## Quality model

Each stable `READ-NNN` case measures separate, explainable dimensions:

- required-text recall;
- forbidden-noise rejection;
- structure preservation;
- Markdown fidelity and required-fact ordering;
- metadata presence;
- Rust site-configuration selection;
- active-content security violations;
- native success/no-content/failure status and elapsed time.

There is deliberately no weighted “quality score.” Rust is eligible for the
frozen-corpus “best evidence” gate only when it has zero failed cases and zero
security violations and is no worse than the best participant on recall, noise
rejection, structure, Markdown fidelity, ordering, and metadata. Elapsed time
cannot compensate for a quality regression. Competitor divergences remain
visible but do not make an otherwise eligible Rust release command fail.

Current frozen-corpus result (2026-08-11):

| Participant | Cases | Required | Noise | Structure | Markdown | Order | Metadata | Security |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| Defuddle | 4/16 | 43/43 | 34/38 | 23/33 | 67/70 | 14/14 | 19/19 | 1 |
| readabilities-py | 3/16 | 43/43 | 32/38 | 22/33 | 60/70 | 14/14 | 19/19 | 2 |
| readabilities-rs | **16/16** | **43/43** | **38/38** | **33/33** | **70/70** | **14/14** | **19/19** | **0** |

See [reports/offline.md](reports/offline.md) for bounded per-case evidence.

## Live comparison

Live mode downloads each URL once and gives the same temporary snapshot to all
three participants. `just live-common` runs the maintained diagnostic set:

- Example Domain;
- English and Chinese Wikipedia;
- MDN;
- a NASA news release;
- the Rust Blog;
- Hacker News;
- the W3C ARIA example.

```sh
just live-common
```

Full page snapshots go to a temporary directory and are deleted by default.
Use `--capture-dir` only when an explicitly preserved snapshot is appropriate.
Reports contain bounded previews, byte counts, and SHA-256 digests—not full
participant bodies.

Live pages have no frozen fact oracle and can drift, so live mode always marks
the “best” gate ineligible. It is useful for discovering crashes, security
residue, pathological size differences, and new corpus candidates.

## Site-specific coverage

The offline release gate includes self-authored, deterministic fixtures for
Medium, Wikipedia, and MDN built-ins, a caller-provided configuration, and an
unconfigured domain. The suite asserts both cleaned Markdown and Rust
provenance. This prevents a site-specific improvement from hiding a weaker
generic fallback. A new optimization starts with a failing fixture and explicit
required/forbidden facts; a copied live page is not accepted as a stable oracle.

## Commands and exit behavior

- `doctor`: resolve and probe every participant; exit 3 if one is unavailable.
- `list`: show/filter stable cases.
- `run --profile offline`: exit 0 when the frozen Rust best-evidence gate is
  eligible, while still reporting competitor failures.
- `run --profile live`: exit 1 if any participant fails the smoke comparison.
- `version`: show tool/schema version.
- malformed usage exits 2.

`just check` runs formatting, vet, unit tests, a built-binary doctor probe, and
the full offline comparison. `just install` embeds the source path into the Go
binary so the installed harness can still locate its corpus and driver.

Design and negative-control evidence are in [docs/validation.md](docs/validation.md).

## License

MIT. See [LICENSE](LICENSE).
