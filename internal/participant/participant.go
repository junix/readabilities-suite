// Package participant resolves and drives the three real public surfaces.
package participant

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	Defuddle = "defuddle"
	Python   = "readabilities-py"
	Rust     = "readabilities-rs"
)

type Config struct {
	SuiteRoot string
	Explicit  map[string]string
}

type Target struct {
	ID        string   `json:"id"`
	Command   string   `json:"command,omitempty"`
	Prefix    []string `json:"prefix,omitempty"`
	WorkDir   string   `json:"work_dir,omitempty"`
	Available bool     `json:"available"`
	Reason    string   `json:"reason,omitempty"`
}

type NativeResult struct {
	Participant    string            `json:"participant"`
	Status         string            `json:"status"`
	Content        string            `json:"-"`
	ContentBytes   int               `json:"content_bytes"`
	ContentSHA256  string            `json:"content_sha256,omitempty"`
	ContentPreview string            `json:"content_preview,omitempty"`
	ContentFormat  string            `json:"content_format,omitempty"`
	SiteConfig     string            `json:"site_config,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	ElapsedMS      int64             `json:"elapsed_ms"`
	Warnings       []string          `json:"warnings,omitempty"`
	ExitCode       int               `json:"exit_code"`
	Stderr         string            `json:"stderr,omitempty"`
}

func Resolve(cfg Config) []Target {
	return []Target{
		resolveDefuddle(cfg),
		resolvePython(cfg),
		resolveRust(cfg),
	}
}

func resolveDefuddle(cfg Config) Target {
	path := cfg.Explicit[Defuddle]
	if path == "" {
		path = os.Getenv("DEFUDDLE_BIN")
	}
	if path == "" {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, "defuddle", "dist", "cli.js")
		}
	}
	t := Target{ID: Defuddle}
	if path == "" {
		t.Reason = "no Defuddle CLI configured"
		return t
	}
	if _, err := os.Stat(path); err != nil {
		t.Reason = err.Error()
		return t
	}
	node, err := exec.LookPath("node")
	if err != nil {
		t.Reason = "node not found"
		return t
	}
	t.Command, t.Prefix, t.Available = node, []string{path}, true
	return t
}

func resolvePython(cfg Config) Target {
	dir := cfg.Explicit[Python]
	if dir == "" {
		dir = os.Getenv("READABILITIES_PY_DIR")
	}
	if dir == "" {
		dir = filepath.Clean(filepath.Join(cfg.SuiteRoot, "..", "readabilities-py"))
	}
	t := Target{ID: Python, WorkDir: dir}
	if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err != nil {
		t.Reason = err.Error()
		return t
	}
	uv, err := exec.LookPath("uv")
	if err != nil {
		t.Reason = "uv not found"
		return t
	}
	driver := filepath.Join(cfg.SuiteRoot, "drivers", "readabilities_py.py")
	if _, err := os.Stat(driver); err != nil {
		t.Reason = err.Error()
		return t
	}
	t.Command, t.Prefix, t.Available = uv, []string{"run", "python", driver}, true
	return t
}

func resolveRust(cfg Config) Target {
	path := cfg.Explicit[Rust]
	if path == "" {
		path = os.Getenv("READABILITIES_RS_BIN")
	}
	if path == "" {
		path = filepath.Clean(filepath.Join(cfg.SuiteRoot, "..", "readabilities-rs", "target", "debug", "readabilities-rs"))
	}
	t := Target{ID: Rust, Command: path}
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		if err != nil {
			t.Reason = err.Error()
		} else {
			t.Reason = "configured Rust path is a directory"
		}
		return t
	}
	t.Available = true
	return t
}

func Run(ctx context.Context, target Target, fixture, baseURL, rustSiteConfig string) NativeResult {
	if target.ID == Rust {
		return runRust(ctx, target, fixture, baseURL, rustSiteConfig)
	}
	started := time.Now()
	args := append([]string{}, target.Prefix...)
	switch target.ID {
	case Defuddle:
		args = append(args, "parse", fixture, "--json", "--markdown")
	case Python:
		args = append(args, fixture, baseURL)
	}
	cmd := exec.CommandContext(ctx, target.Command, args...)
	cmd.Dir = target.WorkDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	result := NativeResult{
		Participant: target.ID,
		ElapsedMS:   time.Since(started).Milliseconds(),
		Metadata:    map[string]string{},
		Stderr:      bounded(stderr.String(), 600),
	}
	if err != nil {
		result.Status = "failure"
		if exit, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exit.ExitCode()
		} else {
			result.ExitCode = -1
			result.Warnings = []string{bounded(err.Error(), 300)}
		}
		return result
	}
	if err := decode(target.ID, stdout.Bytes(), &result); err != nil {
		result.Status = "failure"
		result.ExitCode = -1
		result.Warnings = []string{bounded(err.Error(), 300)}
		return result
	}
	result.summarizeContent()
	result.ExitCode = 0
	return result
}

func runRust(ctx context.Context, target Target, fixture, baseURL, siteConfig string) NativeResult {
	started := time.Now()
	markdownArgs := []string{"extract", fixture, "--format", "markdown"}
	if baseURL != "" {
		markdownArgs = append(markdownArgs, "--url", baseURL)
	}
	if siteConfig != "" {
		markdownArgs = append(markdownArgs, "--site-config", siteConfig)
	}
	markdownCmd := exec.CommandContext(ctx, target.Command, markdownArgs...)
	markdownCmd.Dir = target.WorkDir
	var markdownOut, markdownErr bytes.Buffer
	markdownCmd.Stdout, markdownCmd.Stderr = &markdownOut, &markdownErr
	err := markdownCmd.Run()
	result := NativeResult{
		Participant:   Rust,
		ElapsedMS:     time.Since(started).Milliseconds(),
		Metadata:      map[string]string{},
		Stderr:        bounded(markdownErr.String(), 600),
		Content:       strings.TrimSpace(markdownOut.String()),
		ContentFormat: "markdown",
	}
	if err != nil {
		result.Status = "failure"
		if exit, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exit.ExitCode()
		} else {
			result.ExitCode = -1
			result.Warnings = []string{bounded(err.Error(), 300)}
		}
		return result
	}

	metadataArgs := []string{"extract", fixture, "--json"}
	if baseURL != "" {
		metadataArgs = append(metadataArgs, "--url", baseURL)
	}
	if siteConfig != "" {
		metadataArgs = append(metadataArgs, "--site-config", siteConfig)
	}
	metadataCmd := exec.CommandContext(ctx, target.Command, metadataArgs...)
	metadataCmd.Dir = target.WorkDir
	metadataOut, err := metadataCmd.Output()
	if err != nil {
		result.Status, result.ExitCode = "failure", -1
		result.Warnings = []string{bounded(fmt.Sprintf("Rust metadata projection: %v", err), 300)}
		return result
	}
	var value map[string]any
	if err := json.Unmarshal(metadataOut, &value); err != nil {
		result.Status, result.ExitCode = "failure", -1
		result.Warnings = []string{bounded(fmt.Sprintf("decode Rust metadata JSON: %v", err), 300)}
		return result
	}
	copyStringMap(result.Metadata, value["metadata"])
	if provenance, ok := value["provenance"].(map[string]any); ok {
		result.SiteConfig = stringValue(provenance["site_config"])
	}
	result.Status = "success"
	if result.Content == "" {
		result.Status = "no_content"
	}
	result.ExitCode = 0
	result.summarizeContent()
	return result
}

// summarizeContent records bounded evidence while keeping extracted page bodies
// out of persistent JSON and Markdown reports.
func (result *NativeResult) summarizeContent() {
	result.ContentBytes = len(result.Content)
	if result.Content == "" {
		return
	}
	digest := sha256.Sum256([]byte(result.Content))
	result.ContentSHA256 = fmt.Sprintf("%x", digest)
	result.ContentPreview = bounded(NormalizedText(result.Content), 240)
}

func Probe(ctx context.Context, target Target) error {
	args := append([]string{}, target.Prefix...)
	switch target.ID {
	case Defuddle:
		args = append(args, "--version")
	case Python:
		args = append(args, "--doctor")
	case Rust:
		args = append(args, "doctor", "--json")
	}
	cmd := exec.CommandContext(ctx, target.Command, args...)
	cmd.Dir = target.WorkDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s probe: %w: %s", target.ID, err, bounded(string(output), 300))
	}
	return nil
}

func decode(id string, data []byte, result *NativeResult) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("decode %s JSON: %w", id, err)
	}
	if id == Python {
		result.Status = stringValue(value["status"])
		result.Content = stringValue(value["content"])
		result.ContentFormat = stringValue(value["content_format"])
		copyStringMap(result.Metadata, value["metadata"])
		return nil
	}
	result.Status = "success"
	result.Content = stringValue(value["content"])
	result.ContentFormat = "markdown"
	{
		for _, key := range []string{"title", "author", "published", "site", "language"} {
			if v := stringValue(value[key]); v != "" {
				result.Metadata[key] = v
			}
		}
	}
	if strings.TrimSpace(result.Content) == "" {
		result.Status = "no_content"
	}
	return nil
}

func copyStringMap(dst map[string]string, raw any) {
	values, ok := raw.(map[string]any)
	if !ok {
		return
	}
	for key, value := range values {
		if text := stringValue(value); text != "" {
			dst[key] = text
		}
	}
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return fmt.Sprint(value)
}

var tagPattern = regexp.MustCompile(`(?s)<[^>]*>`)

func NormalizedText(content string) string {
	withoutTags := tagPattern.ReplaceAllString(content, " ")
	return strings.Join(strings.Fields(html.UnescapeString(withoutTags)), " ")
}

func bounded(value string, max int) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max]) + "…"
}
