// Package suite evaluates normalized facts without changing participant content.
package suite

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/junix/readabilities-suite/internal/corpus"
	"github.com/junix/readabilities-suite/internal/participant"
)

const SchemaVersion = 2

type Options struct {
	Jobs         int
	Timeout      time.Duration
	Profile      string
	Participants []string
}

type FactResult struct {
	Fact string `json:"fact"`
	Pass bool   `json:"pass"`
}

type CaseResult struct {
	CaseID       string                   `json:"case_id"`
	CaseName     string                   `json:"case_name"`
	Participant  string                   `json:"participant"`
	Status       string                   `json:"status"`
	NativeStatus string                   `json:"native_status"`
	Required     []FactResult             `json:"required,omitempty"`
	Forbidden    []FactResult             `json:"forbidden,omitempty"`
	Structure    []FactResult             `json:"structure,omitempty"`
	Markdown     []FactResult             `json:"markdown_fidelity,omitempty"`
	Order        []FactResult             `json:"content_order,omitempty"`
	Metadata     []FactResult             `json:"metadata,omitempty"`
	SiteConfig   []FactResult             `json:"site_config,omitempty"`
	Security     []string                 `json:"security_violations,omitempty"`
	Stats        MarkdownStats            `json:"markdown_stats"`
	ElapsedMS    int64                    `json:"elapsed_ms"`
	Proof        string                   `json:"proof,omitempty"`
	Native       participant.NativeResult `json:"native"`
}

type MarkdownStats struct {
	Lines          int `json:"lines"`
	Words          int `json:"words"`
	Headings       int `json:"headings"`
	Links          int `json:"links"`
	Images         int `json:"images"`
	CodeFences     int `json:"code_fences"`
	Tables         int `json:"tables"`
	NoiseHits      int `json:"noise_hits"`
	DuplicateLines int `json:"duplicate_lines"`
}

type Metrics struct {
	Participant        string `json:"participant"`
	Pass               int    `json:"pass"`
	Fail               int    `json:"fail"`
	RequiredFound      int    `json:"required_found"`
	RequiredTotal      int    `json:"required_total"`
	ForbiddenRejected  int    `json:"forbidden_rejected"`
	ForbiddenTotal     int    `json:"forbidden_total"`
	StructureFound     int    `json:"structure_found"`
	StructureTotal     int    `json:"structure_total"`
	MarkdownFound      int    `json:"markdown_found"`
	MarkdownTotal      int    `json:"markdown_total"`
	OrderFound         int    `json:"order_found"`
	OrderTotal         int    `json:"order_total"`
	MetadataFound      int    `json:"metadata_found"`
	MetadataTotal      int    `json:"metadata_total"`
	SecurityViolations int    `json:"security_violations"`
	ElapsedMS          int64  `json:"elapsed_ms"`
	MarkdownBytes      int    `json:"markdown_bytes"`
	MarkdownWords      int    `json:"markdown_words"`
	NoiseHits          int    `json:"noise_hits"`
	DuplicateLines     int    `json:"duplicate_lines"`
}

type BestEvidence struct {
	Eligible                 bool     `json:"eligible"`
	RustSecurityZero         bool     `json:"rust_security_zero"`
	RustRecallAtLeastBest    bool     `json:"rust_recall_at_least_best"`
	RustNoiseAtLeastBest     bool     `json:"rust_noise_at_least_best"`
	RustStructureAtLeastBest bool     `json:"rust_structure_at_least_best"`
	RustMarkdownAtLeastBest  bool     `json:"rust_markdown_at_least_best"`
	RustOrderAtLeastBest     bool     `json:"rust_order_at_least_best"`
	RustMetadataAtLeastBest  bool     `json:"rust_metadata_at_least_best"`
	RecallLeaders            []string `json:"recall_leaders"`
	NoiseLeaders             []string `json:"noise_leaders"`
	StructureLeaders         []string `json:"structure_leaders"`
	MetadataLeaders          []string `json:"metadata_leaders"`
	MarkdownLeaders          []string `json:"markdown_leaders"`
	OrderLeaders             []string `json:"order_leaders"`
	Note                     string   `json:"note"`
}

type Report struct {
	SchemaVersion int                  `json:"schema_version"`
	GeneratedAt   string               `json:"generated_at"`
	Profile       string               `json:"profile"`
	Targets       []participant.Target `json:"targets"`
	Total         int                  `json:"total"`
	Passed        int                  `json:"passed"`
	Failed        int                  `json:"failed"`
	Metrics       []Metrics            `json:"metrics"`
	Best          BestEvidence         `json:"best_evidence"`
	Cases         []CaseResult         `json:"cases"`
}

func Run(ctx context.Context, targets []participant.Target, cases []corpus.Case, opts Options) (Report, error) {
	if opts.Jobs < 1 {
		return Report{}, fmt.Errorf("jobs must be positive")
	}
	if opts.Timeout <= 0 {
		return Report{}, fmt.Errorf("timeout must be positive")
	}
	targets = filterTargets(targets, opts.Participants)
	if len(targets) == 0 {
		return Report{}, fmt.Errorf("no participants selected")
	}
	for _, target := range targets {
		if !target.Available {
			return Report{}, fmt.Errorf("required participant %s unavailable: %s", target.ID, target.Reason)
		}
	}

	type job struct {
		index   int
		caseDef corpus.Case
		target  participant.Target
	}
	jobs := make(chan job)
	results := make([]CaseResult, len(cases)*len(targets))
	var wg sync.WaitGroup
	workers := opts.Jobs
	if workers > len(results) {
		workers = len(results)
	}
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for work := range jobs {
				caseCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
				native := participant.Run(caseCtx, work.target, work.caseDef.Fixture, work.caseDef.BaseURL, work.caseDef.RustSiteConfigFile)
				cancel()
				results[work.index] = evaluate(work.caseDef, native)
			}
		}()
	}
	index := 0
	for _, c := range cases {
		for _, target := range targets {
			jobs <- job{index: index, caseDef: c, target: target}
			index++
		}
	}
	close(jobs)
	wg.Wait()

	report := Report{
		SchemaVersion: SchemaVersion,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Profile:       opts.Profile,
		Targets:       targets,
		Total:         len(results),
		Cases:         results,
	}
	for _, result := range results {
		if result.Status == "pass" {
			report.Passed++
		} else {
			report.Failed++
		}
	}
	report.Metrics = aggregate(targets, results)
	report.Best = bestEvidence(report.Metrics)
	if opts.Profile != "offline" {
		report.Best.Eligible = false
		report.Best.Note = "Live pages have no frozen fact oracle and cannot establish release eligibility."
	}
	return report, nil
}

func evaluate(c corpus.Case, native participant.NativeResult) CaseResult {
	result := CaseResult{
		CaseID: c.ID, CaseName: c.Name, Participant: native.Participant,
		Status: "pass", NativeStatus: native.Status, ElapsedMS: native.ElapsedMS, Native: native,
		Stats: markdownStats(native.Content),
	}
	if !contains(c.ExpectedStatuses, native.Status) {
		result.Status = "fail"
		result.Proof = fmt.Sprintf("expected status %v, actual %q; stderr=%q", c.ExpectedStatuses, native.Status, native.Stderr)
		return result
	}
	if native.Status != "success" {
		result.Proof = "validation/no-content path remained distinguishable"
		return result
	}
	formatPass := native.ContentFormat == "markdown"
	result.Markdown = append(result.Markdown, FactResult{Fact: "native Markdown output", Pass: formatPass})
	if !formatPass {
		result.Status = "fail"
	}
	for _, check := range markdownSyntaxChecks(native.Content) {
		result.Markdown = append(result.Markdown, check)
		if !check.Pass {
			result.Status = "fail"
		}
	}
	text := participant.NormalizedText(native.Content)
	for _, fact := range c.RequiredText {
		pass := strings.Contains(text, fact)
		result.Required = append(result.Required, FactResult{Fact: fact, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	for _, fact := range c.ForbiddenText {
		pass := !strings.Contains(text, fact)
		result.Forbidden = append(result.Forbidden, FactResult{Fact: fact, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	for _, structure := range c.RequiredStructure {
		pass := hasStructure(native.Content, structure)
		result.Structure = append(result.Structure, FactResult{Fact: structure, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	for _, fact := range c.RequiredMarkdown {
		pass := strings.Contains(native.Content, fact)
		result.Markdown = append(result.Markdown, FactResult{Fact: "contains " + fact, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	for _, fact := range c.ForbiddenMarkdown {
		pass := !strings.Contains(native.Content, fact)
		result.Markdown = append(result.Markdown, FactResult{Fact: "rejects " + fact, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	if len(c.RequiredText) > 1 {
		pass := ordered(text, c.RequiredText)
		result.Order = append(result.Order, FactResult{Fact: "required facts preserve source order", Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	for _, field := range c.MetadataRequired {
		pass := strings.TrimSpace(native.Metadata[field]) != ""
		result.Metadata = append(result.Metadata, FactResult{Fact: field, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	if native.Participant == participant.Rust && c.ExpectedRustSiteConfig != nil {
		expected := *c.ExpectedRustSiteConfig
		pass := native.SiteConfig == expected
		fact := "no site config; generic fallback"
		if expected != "" {
			fact = "site config " + expected
		}
		result.SiteConfig = append(result.SiteConfig, FactResult{Fact: fact, Pass: pass})
		if !pass {
			result.Status = "fail"
		}
	}
	result.Security = securityViolations(native.Content)
	if len(result.Security) > 0 {
		result.Status = "fail"
	}
	if result.Status == "fail" {
		result.Proof = boundedProof(c, result, text)
	} else {
		result.Proof = fmt.Sprintf("required=%d forbidden=%d structure=%d markdown=%d order=%d metadata=%d site=%d security=0", len(result.Required), len(result.Forbidden), len(result.Structure), len(result.Markdown), len(result.Order), len(result.Metadata), len(result.SiteConfig))
	}
	return result
}

func hasStructure(content, name string) bool {
	patterns := map[string]*regexp.Regexp{
		"heading":    regexp.MustCompile(`(?m)^#{1,6}\s+\S`),
		"code":       regexp.MustCompile(`(?m)^(?:` + "```" + `|~~~)`),
		"table":      regexp.MustCompile(`(?m)^\|.*\|\s*$`),
		"image":      regexp.MustCompile(`!\[[^\]]*\]\([^\n)]+\)`),
		"footnote":   regexp.MustCompile(`(?m)(?:\[\^[^\]]+\]|<sup>[^<]+</sup>)`),
		"math":       regexp.MustCompile(`(?s)\$\$.+?\$\$|\$[^\n$]+\$`),
		"list":       regexp.MustCompile(`(?m)^\s*(?:[-+*]|\d+\.)\s+\S`),
		"blockquote": regexp.MustCompile(`(?m)^>\s+\S`),
	}
	pattern, ok := patterns[name]
	return ok && pattern.MatchString(content)
}

var eventAttribute = regexp.MustCompile(`(?i)\son[a-z]+\s*=`)
var dangerousURIAttribute = regexp.MustCompile(`(?i)(?:href|src|action|formaction)\s*=\s*["']?\s*(?:javascript:|data:text/html)`)
var dangerousMarkdownURI = regexp.MustCompile(`(?i)!?\[[^\]]*\]\(\s*(?:javascript:|data:text/html)`)
var collapsedImageText = regexp.MustCompile(`!\[[^\]]*\]\([^\n)]+\)[\pL\pN]`)

func securityViolations(content string) []string {
	lower := strings.ToLower(content)
	checks := []struct{ name, needle string }{
		{"script_element", "<script"},
		{"iframe_element", "<iframe"},
	}
	var violations []string
	for _, check := range checks {
		if strings.Contains(lower, check.needle) {
			violations = append(violations, check.name)
		}
	}
	if eventAttribute.MatchString(content) {
		violations = append(violations, "event_handler_attribute")
	}
	if dangerousURIAttribute.MatchString(content) {
		violations = append(violations, "dangerous_uri_attribute")
	}
	if dangerousMarkdownURI.MatchString(content) {
		violations = append(violations, "dangerous_markdown_uri")
	}
	return violations
}

func ordered(text string, facts []string) bool {
	position := 0
	for _, fact := range facts {
		index := strings.Index(text[position:], fact)
		if index < 0 {
			return false
		}
		position += index + len(fact)
	}
	return true
}

func markdownSyntaxChecks(content string) []FactResult {
	return []FactResult{
		{Fact: "balanced fenced code blocks", Pass: balancedFences(content)},
		{Fact: "image and following prose remain separated", Pass: !collapsedImageText.MatchString(content)},
		{Fact: "no NUL bytes", Pass: !strings.ContainsRune(content, '\x00')},
	}
}

func balancedFences(content string) bool {
	marker := ""
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		candidate := ""
		switch {
		case strings.HasPrefix(trimmed, "```"):
			candidate = "```"
		case strings.HasPrefix(trimmed, "~~~"):
			candidate = "~~~"
		}
		if candidate == "" {
			continue
		}
		if marker == "" {
			marker = candidate
		} else if marker == candidate {
			marker = ""
		}
	}
	return marker == ""
}

func markdownStats(content string) MarkdownStats {
	if strings.TrimSpace(content) == "" {
		return MarkdownStats{}
	}
	headingPattern := regexp.MustCompile(`(?m)^#{1,6}\s+\S`)
	linkPattern := regexp.MustCompile(`\[[^\]]+\]\([^\n)]+\)`)
	imagePattern := regexp.MustCompile(`!\[[^\]]*\]\([^\n)]+\)`)
	fencePattern := regexp.MustCompile(`(?m)^(?:` + "```" + `|~~~)`)
	tablePattern := regexp.MustCompile(`(?m)^\|\s*:?-{3,}`)
	lower := strings.ToLower(content)
	noiseHits := 0
	for _, phrase := range []string{
		"cookie policy", "accept cookies", "sign in", "subscribe", "newsletter",
		"related articles", "related stories", "privacy policy", "all rights reserved",
		"was this page helpful", "previous page", "next page",
	} {
		noiseHits += strings.Count(lower, phrase)
	}
	seen := map[string]int{}
	duplicateLines := 0
	for _, line := range strings.Split(content, "\n") {
		line = strings.Join(strings.Fields(line), " ")
		if len([]rune(line)) < 40 || strings.HasPrefix(line, "|") {
			continue
		}
		seen[line]++
		if seen[line] == 2 {
			duplicateLines++
		}
	}
	return MarkdownStats{
		Lines:          len(strings.Split(content, "\n")),
		Words:          len(strings.Fields(participant.NormalizedText(content))),
		Headings:       len(headingPattern.FindAllStringIndex(content, -1)),
		Links:          len(linkPattern.FindAllStringIndex(content, -1)),
		Images:         len(imagePattern.FindAllStringIndex(content, -1)),
		CodeFences:     len(fencePattern.FindAllStringIndex(content, -1)) / 2,
		Tables:         len(tablePattern.FindAllStringIndex(content, -1)),
		NoiseHits:      noiseHits,
		DuplicateLines: duplicateLines,
	}
}

func aggregate(targets []participant.Target, results []CaseResult) []Metrics {
	metrics := make(map[string]*Metrics, len(targets))
	for _, target := range targets {
		metrics[target.ID] = &Metrics{Participant: target.ID}
	}
	for _, result := range results {
		m := metrics[result.Participant]
		if result.Status == "pass" {
			m.Pass++
		} else {
			m.Fail++
		}
		m.ElapsedMS += result.ElapsedMS
		m.RequiredFound += passCount(result.Required)
		m.RequiredTotal += len(result.Required)
		m.ForbiddenRejected += passCount(result.Forbidden)
		m.ForbiddenTotal += len(result.Forbidden)
		m.StructureFound += passCount(result.Structure)
		m.StructureTotal += len(result.Structure)
		m.MarkdownFound += passCount(result.Markdown)
		m.MarkdownTotal += len(result.Markdown)
		m.OrderFound += passCount(result.Order)
		m.OrderTotal += len(result.Order)
		m.MetadataFound += passCount(result.Metadata)
		m.MetadataTotal += len(result.Metadata)
		m.SecurityViolations += len(result.Security)
		m.MarkdownBytes += result.Native.ContentBytes
		m.MarkdownWords += result.Stats.Words
		m.NoiseHits += result.Stats.NoiseHits
		m.DuplicateLines += result.Stats.DuplicateLines
	}
	out := make([]Metrics, 0, len(targets))
	for _, target := range targets {
		out = append(out, *metrics[target.ID])
	}
	return out
}

func bestEvidence(metrics []Metrics) BestEvidence {
	byID := map[string]Metrics{}
	for _, metric := range metrics {
		byID[metric.Participant] = metric
	}
	rust, rustOK := byID[participant.Rust]
	evidence := BestEvidence{
		RustSecurityZero: rustOK && rust.SecurityViolations == 0,
		Note:             "Quality dimensions are release gates; elapsed time cannot compensate for lost content or retained noise. Eligibility is a corpus claim, not universal superiority.",
	}
	evidence.RustRecallAtLeastBest = rustOK && atLeastBest(metrics, rust.RequiredFound, rust.RequiredTotal, func(m Metrics) (int, int) { return m.RequiredFound, m.RequiredTotal })
	evidence.RustNoiseAtLeastBest = rustOK && atLeastBest(metrics, rust.ForbiddenRejected, rust.ForbiddenTotal, func(m Metrics) (int, int) { return m.ForbiddenRejected, m.ForbiddenTotal })
	evidence.RustStructureAtLeastBest = rustOK && atLeastBest(metrics, rust.StructureFound, rust.StructureTotal, func(m Metrics) (int, int) { return m.StructureFound, m.StructureTotal })
	evidence.RustMarkdownAtLeastBest = rustOK && atLeastBest(metrics, rust.MarkdownFound, rust.MarkdownTotal, func(m Metrics) (int, int) { return m.MarkdownFound, m.MarkdownTotal })
	evidence.RustOrderAtLeastBest = rustOK && atLeastBest(metrics, rust.OrderFound, rust.OrderTotal, func(m Metrics) (int, int) { return m.OrderFound, m.OrderTotal })
	evidence.RustMetadataAtLeastBest = rustOK && atLeastBest(metrics, rust.MetadataFound, rust.MetadataTotal, func(m Metrics) (int, int) { return m.MetadataFound, m.MetadataTotal })
	evidence.Eligible = evidence.RustSecurityZero && evidence.RustRecallAtLeastBest && evidence.RustNoiseAtLeastBest && evidence.RustStructureAtLeastBest && evidence.RustMarkdownAtLeastBest && evidence.RustOrderAtLeastBest && evidence.RustMetadataAtLeastBest && rust.Fail == 0
	evidence.RecallLeaders = leaders(metrics, func(m Metrics) (int, int) { return m.RequiredFound, m.RequiredTotal })
	evidence.NoiseLeaders = leaders(metrics, func(m Metrics) (int, int) { return m.ForbiddenRejected, m.ForbiddenTotal })
	evidence.StructureLeaders = leaders(metrics, func(m Metrics) (int, int) { return m.StructureFound, m.StructureTotal })
	evidence.MetadataLeaders = leaders(metrics, func(m Metrics) (int, int) { return m.MetadataFound, m.MetadataTotal })
	evidence.MarkdownLeaders = leaders(metrics, func(m Metrics) (int, int) { return m.MarkdownFound, m.MarkdownTotal })
	evidence.OrderLeaders = leaders(metrics, func(m Metrics) (int, int) { return m.OrderFound, m.OrderTotal })
	return evidence
}

func atLeastBest(metrics []Metrics, rustN, rustD int, value func(Metrics) (int, int)) bool {
	for _, metric := range metrics {
		n, d := value(metric)
		if !ratioAtLeast(rustN, rustD, n, d) {
			return false
		}
	}
	return true
}

func leaders(metrics []Metrics, value func(Metrics) (int, int)) []string {
	bestN, bestD := -1, 1
	var result []string
	for _, metric := range metrics {
		n, d := value(metric)
		if bestN < 0 || ratioAtLeast(n, d, bestN, bestD) && !ratioAtLeast(bestN, bestD, n, d) {
			bestN, bestD, result = n, d, []string{metric.Participant}
		} else if ratioAtLeast(n, d, bestN, bestD) && ratioAtLeast(bestN, bestD, n, d) {
			result = append(result, metric.Participant)
		}
	}
	sort.Strings(result)
	return result
}

func ratioAtLeast(a, b, c, d int) bool {
	if b == 0 {
		a, b = 1, 1
	}
	if d == 0 {
		c, d = 1, 1
	}
	return a*d >= c*b
}

func filterTargets(targets []participant.Target, ids []string) []participant.Target {
	if len(ids) == 0 {
		return targets
	}
	var out []participant.Target
	for _, target := range targets {
		if contains(ids, target.ID) {
			out = append(out, target)
		}
	}
	return out
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if strings.EqualFold(candidate, value) {
			return true
		}
	}
	return false
}

func passCount(values []FactResult) int {
	count := 0
	for _, value := range values {
		if value.Pass {
			count++
		}
	}
	return count
}

func boundedProof(c corpus.Case, result CaseResult, text string) string {
	var failed []string
	for _, group := range [][]FactResult{result.Required, result.Forbidden, result.Structure, result.Markdown, result.Order, result.Metadata, result.SiteConfig} {
		for _, fact := range group {
			if !fact.Pass {
				failed = append(failed, fact.Fact)
			}
		}
	}
	preview := []rune(text)
	if len(preview) > 240 {
		preview = preview[:240]
	}
	return fmt.Sprintf("case=%s failed_facts=%q security=%v actual_preview=%q", c.ID, failed, result.Security, string(preview))
}

func WriteJSON(path string, report Report) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if path == "-" {
		fmt.Println(string(data))
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
