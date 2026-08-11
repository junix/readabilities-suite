set shell := ["bash", "-uo", "pipefail", "-c"]

source_path := justfile_directory()
build_flags := "-X main.sourcePath=" + source_path + " -X main.version=0.1.0"
os_name := if os() == "macos" { "macos" } else { "linux" }
arch_name := if arch() == "aarch64" { "arm64" } else { "x86" }
default_install_bin := home_directory() / "sync" / (os_name + "-" + arch_name + "-bin")
install_bin := env("SYNC_BIN_DIR", default_install_bin)

default: build

build:
  go build -ldflags "{{ build_flags }}" -o ./readabilities-suite ./cmd/readabilities-suite
  go build ./...

test:
  go test ./...

check:
  #!/usr/bin/env bash
  set -euo pipefail
  unformatted="$(gofmt -l cmd internal)"
  if [[ -n "$unformatted" ]]; then
    echo "unformatted Go files:" >&2
    echo "$unformatted" >&2
    exit 1
  fi
  go vet ./...
  go test ./...
  just build
  ./readabilities-suite doctor
  ./readabilities-suite run --profile offline >/dev/null

# Non-gating diagnostics against maintained common public sites.
live-common: build
  ./readabilities-suite run --profile live --timeout 90s \
    --url https://example.com/ \
    --url https://en.wikipedia.org/wiki/Readability \
    --url https://zh.wikipedia.org/wiki/%E5%8F%AF%E8%AF%BB%E6%80%A7 \
    --url https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article \
    --url https://www.nasa.gov/news-release/nasa-to-cover-three-us-spacewalks-host-preview-news-conference/ \
    --url https://blog.rust-lang.org/2025/02/20/Rust-1.85.0/ \
    --url https://news.ycombinator.com/ \
    --url https://www.w3.org/WAI/ARIA/apg/patterns/accordion/examples/accordion/ \
    --report reports/live-2026-08-11.json \
    --markdown-report reports/live-2026-08-11.md

tidy:
  go mod tidy

install: build
  #!/usr/bin/env bash
  set -euo pipefail
  mkdir -p "{{ install_bin }}"
  cp ./readabilities-suite "{{ install_bin }}/readabilities-suite"
  echo "Installed {{ install_bin }}/readabilities-suite"

clean:
  rm -f ./readabilities-suite
