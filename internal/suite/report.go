package suite

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Print(w io.Writer, report Report) {
	fmt.Fprintf(w, "readabilities-suite profile=%s: %d passed, %d failed\n\n", report.Profile, report.Passed, report.Failed)
	fmt.Fprintln(w, "participant         pass fail required noise structure markdown order metadata golden security md_bytes md_words noise_hits dup_lines elapsed_ms")
	for _, metric := range report.Metrics {
		fmt.Fprintf(w, "%-19s %4d %4d %7d/%-3d %3d/%-3d %5d/%-3d %5d/%-3d %3d/%-3d %5d/%-3d %5d/%-3d %8d %8d %8d %10d %9d %10d\n",
			metric.Participant, metric.Pass, metric.Fail,
			metric.RequiredFound, metric.RequiredTotal,
			metric.ForbiddenRejected, metric.ForbiddenTotal,
			metric.StructureFound, metric.StructureTotal,
			metric.MarkdownFound, metric.MarkdownTotal,
			metric.OrderFound, metric.OrderTotal,
			metric.MetadataFound, metric.MetadataTotal,
			metric.GoldenPass, metric.GoldenTotal,
			metric.SecurityViolations, metric.MarkdownBytes, metric.MarkdownWords,
			metric.NoiseHits, metric.DuplicateLines, metric.ElapsedMS)
	}
	fmt.Fprintf(w, "\nfrozen-corpus best eligibility: %t (%s)\n", report.Best.Eligible, report.Best.Note)
	if report.Failed > 0 {
		fmt.Fprintln(w, "\nfailures:")
		for _, result := range report.Cases {
			if result.Status == "fail" {
				fmt.Fprintf(w, "- %s/%s: %s\n", result.CaseID, result.Participant, result.Proof)
			}
		}
	}
}

func WriteMarkdown(path string, report Report) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Readabilities comparison report\n\nGenerated: %s\n\nProfile: `%s`\n\nResult: **%d pass / %d fail**\n\n", report.GeneratedAt, report.Profile, report.Passed, report.Failed)
	b.WriteString("| Participant | Pass | Fail | Recall | Noise rejection | Structure | Markdown fidelity | Order | Metadata | Defuddle Golden | Security | Markdown bytes | Words | Noise hints | Duplicate lines | Elapsed ms |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, m := range report.Metrics {
		fmt.Fprintf(&b, "| %s | %d | %d | %d/%d | %d/%d | %d/%d | %d/%d | %d/%d | %d/%d | %d/%d | %d | %d | %d | %d | %d | %d |\n", m.Participant, m.Pass, m.Fail, m.RequiredFound, m.RequiredTotal, m.ForbiddenRejected, m.ForbiddenTotal, m.StructureFound, m.StructureTotal, m.MarkdownFound, m.MarkdownTotal, m.OrderFound, m.OrderTotal, m.MetadataFound, m.MetadataTotal, m.GoldenPass, m.GoldenTotal, m.SecurityViolations, m.MarkdownBytes, m.MarkdownWords, m.NoiseHits, m.DuplicateLines, m.ElapsedMS)
	}
	fmt.Fprintf(&b, "\n## Best-evidence gate\n\n- Eligible on frozen corpus: `%t`\n- Rust security zero: `%t`\n- Rust Defuddle Golden conformant: `%t`\n- Rust recall at least best participant: `%t`\n- Rust noise rejection at least best participant: `%t`\n- Rust structure at least best participant: `%t`\n- Rust Markdown fidelity at least best participant: `%t`\n- Rust content order at least best participant: `%t`\n- Rust metadata at least best participant: `%t`\n- Recall leaders: `%s`\n- Noise leaders: `%s`\n- Structure leaders: `%s`\n- Markdown leaders: `%s`\n- Order leaders: `%s`\n- Metadata leaders: `%s`\n\n%s\n", report.Best.Eligible, report.Best.RustSecurityZero, report.Best.RustGoldenConformant, report.Best.RustRecallAtLeastBest, report.Best.RustNoiseAtLeastBest, report.Best.RustStructureAtLeastBest, report.Best.RustMarkdownAtLeastBest, report.Best.RustOrderAtLeastBest, report.Best.RustMetadataAtLeastBest, strings.Join(report.Best.RecallLeaders, ", "), strings.Join(report.Best.NoiseLeaders, ", "), strings.Join(report.Best.StructureLeaders, ", "), strings.Join(report.Best.MarkdownLeaders, ", "), strings.Join(report.Best.OrderLeaders, ", "), strings.Join(report.Best.MetadataLeaders, ", "), report.Best.Note)
	b.WriteString("\n## Per-case Markdown evidence\n\n")
	b.WriteString("| Case | Participant | Status | Recall | Noise | Structure | Markdown | Order | Golden | Bytes | Words | Headings | Links | Code | Noise hints | Duplicates | Site config | SHA-256 |\n")
	b.WriteString("|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---|\n")
	for _, result := range report.Cases {
		digest := result.Native.ContentSHA256
		if len(digest) > 12 {
			digest = digest[:12]
		}
		golden := "-"
		if result.Golden != nil {
			golden = fmt.Sprintf("%.0f%%", result.Golden.TokenCoverage*100)
			if !result.Golden.Pass {
				golden += " fail"
			}
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %d/%d | %d/%d | %d/%d | %d/%d | %d/%d | %s | %d | %d | %d | %d | %d | %d | %d | %s | `%s` |\n", result.CaseID, result.Participant, result.Status, passCount(result.Required), len(result.Required), passCount(result.Forbidden), len(result.Forbidden), passCount(result.Structure), len(result.Structure), passCount(result.Markdown), len(result.Markdown), passCount(result.Order), len(result.Order), golden, result.Native.ContentBytes, result.Stats.Words, result.Stats.Headings, result.Stats.Links, result.Stats.CodeFences, result.Stats.NoiseHits, result.Stats.DuplicateLines, result.Native.SiteConfig, digest)
	}
	if report.Failed > 0 {
		b.WriteString("\n## Failures\n\n")
		for _, result := range report.Cases {
			if result.Status == "fail" {
				fmt.Fprintf(&b, "- `%s/%s`: %s\n", result.CaseID, result.Participant, result.Proof)
			}
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
