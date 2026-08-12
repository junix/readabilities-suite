// Package cli implements doctor/list/run/version for the independent suite.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/junix/readabilities-suite/internal/corpus"
	"github.com/junix/readabilities-suite/internal/participant"
	"github.com/junix/readabilities-suite/internal/suite"
)

const (
	exitOK          = 0
	exitFailed      = 1
	exitUsage       = 2
	exitUnavailable = 3
)

const defuddleGoldenSource = "defuddle/0.19.2-markdown-v1"

type BuildInfo struct{ Version, SourcePath string }

type globals struct {
	json     bool
	explicit map[string]string
	root     string
}

type listFlag []string

func (l *listFlag) String() string { return strings.Join(*l, ",") }
func (l *listFlag) Set(value string) error {
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			*l = append(*l, item)
		}
	}
	return nil
}

func Run(args []string, stdout, stderr io.Writer, build BuildInfo) int {
	global, rest, err := parseGlobals(args, build)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(rest) == 0 {
		printHelp(stdout)
		return exitUsage
	}
	switch rest[0] {
	case "help", "-h", "--help":
		printHelp(stdout)
		return exitOK
	case "version", "--version", "-V":
		return runVersion(rest[1:], global, build, stdout, stderr)
	case "doctor":
		return runDoctor(rest[1:], global, stdout, stderr)
	case "list":
		return runList(rest[1:], global, stdout, stderr)
	case "run", "test":
		return runCases(rest[1:], global, stdout, stderr)
	case "record-golden":
		return recordGolden(rest[1:], global, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", rest[0])
		printHelp(stderr)
		return exitUsage
	}
}

func parseGlobals(args []string, build BuildInfo) (globals, []string, error) {
	root := build.SourcePath
	if root == "" || root == "(unknown)" {
		if cwd, err := os.Getwd(); err == nil {
			root = cwd
		}
	}
	g := globals{explicit: map[string]string{}, root: root}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			g.json = true
		case "--root":
			if i+1 >= len(args) {
				return g, nil, errors.New("--root requires a value")
			}
			i++
			g.root = args[i]
		case "--rust":
			if i+1 >= len(args) {
				return g, nil, errors.New("--rust requires a value")
			}
			i++
			g.explicit[participant.Rust] = args[i]
		case "--defuddle":
			if i+1 >= len(args) {
				return g, nil, errors.New("--defuddle requires a value")
			}
			i++
			g.explicit[participant.Defuddle] = args[i]
		case "--python":
			if i+1 >= len(args) {
				return g, nil, errors.New("--python requires a value")
			}
			i++
			g.explicit[participant.Python] = args[i]
		default:
			return g, args[i:], nil
		}
	}
	return g, nil, nil
}

func targets(g globals) []participant.Target {
	return participant.Resolve(participant.Config{SuiteRoot: g.root, Explicit: g.explicit})
}

func runDoctor(args []string, g globals, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("readabilities-suite doctor", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", g.json, "emit JSON")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	targets := targets(g)
	allReady := true
	for i := range targets {
		if !targets[i].Available {
			allReady = false
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := participant.Probe(ctx, targets[i])
		cancel()
		if err != nil {
			targets[i].Available = false
			targets[i].Reason = err.Error()
			allReady = false
		}
	}
	if *jsonOutput {
		_ = json.NewEncoder(stdout).Encode(map[string]any{"schema_version": 1, "ready": allReady, "targets": targets})
	} else {
		for _, target := range targets {
			status := "ok"
			if !target.Available {
				status = "unavailable: " + target.Reason
			}
			fmt.Fprintf(stdout, "%-18s %s\n", target.ID, status)
		}
	}
	if !allReady {
		return exitUnavailable
	}
	return exitOK
}

func runList(args []string, g globals, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("readabilities-suite list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", g.json, "emit JSON")
	var selectors, tags listFlag
	fs.Var(&selectors, "case", "case id/name; repeatable")
	fs.Var(&tags, "tag", "tag; repeatable")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	cases, err := corpus.Load(g.root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	cases, err = corpus.Filter(cases, selectors, tags)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if *jsonOutput {
		_ = json.NewEncoder(stdout).Encode(cases)
		return exitOK
	}
	for _, c := range cases {
		fmt.Fprintf(stdout, "%s\t%-24s\t%s\n", c.ID, c.Name, strings.Join(c.Tags, ","))
	}
	return exitOK
}

func runCases(args []string, g globals, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("readabilities-suite run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	profile := fs.String("profile", "offline", "offline or live")
	jsonOutput := fs.Bool("json", g.json, "emit JSON")
	jobs := fs.Int("jobs", 3, "parallel participant jobs")
	fs.IntVar(jobs, "j", 3, "alias for --jobs")
	timeout := fs.Duration("timeout", 60*time.Second, "timeout per participant/case")
	reportPath := fs.String("report", "", "write JSON report")
	markdownPath := fs.String("markdown-report", "", "write Markdown report")
	captureDir := fs.String("capture-dir", "", "preserve live snapshots in this directory")
	var selectors, tags, participants, urls listFlag
	fs.Var(&selectors, "case", "case id/name; repeatable")
	fs.Var(&tags, "tag", "tag; repeatable")
	fs.Var(&participants, "participant", "participant id; repeatable")
	fs.Var(&urls, "url", "live URL; repeatable")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	var cases []corpus.Case
	var cleanup func()
	var err error
	switch *profile {
	case "offline":
		cases, err = corpus.Load(g.root)
		if err == nil {
			cases, err = corpus.Filter(cases, selectors, tags)
		}
		if err == nil {
			err = corpus.RequireGoldens(cases)
		}
	case "live":
		if len(urls) == 0 {
			fmt.Fprintln(stderr, "live profile requires at least one --url")
			return exitUsage
		}
		cases, cleanup, err = captureLive(urls, *captureDir, *timeout)
	default:
		err = fmt.Errorf("unknown profile %q", *profile)
	}
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	report, err := suite.Run(context.Background(), targets(g), cases, suite.Options{Jobs: *jobs, Timeout: *timeout, Profile: *profile, Participants: participants})
	if err != nil {
		fmt.Fprintln(stderr, err)
		if strings.Contains(err.Error(), "unavailable") {
			return exitUnavailable
		}
		return exitUsage
	}
	if *reportPath != "" {
		if err := suite.WriteJSON(*reportPath, report); err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
	}
	if *markdownPath != "" {
		if err := suite.WriteMarkdown(*markdownPath, report); err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
	}
	if *jsonOutput {
		_ = json.NewEncoder(stdout).Encode(report)
	} else {
		suite.Print(stdout, report)
	}
	if *profile == "offline" && !report.Best.Eligible {
		return exitFailed
	}
	if *profile == "live" && report.Failed > 0 {
		return exitFailed
	}
	return exitOK
}

func recordGolden(args []string, g globals, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("readabilities-suite record-golden", flag.ContinueOnError)
	fs.SetOutput(stderr)
	force := fs.Bool("force", false, "overwrite existing reviewed Golden files")
	timeout := fs.Duration("timeout", 60*time.Second, "timeout per Defuddle extraction")
	var selectors, tags listFlag
	fs.Var(&selectors, "case", "case id/name; repeatable")
	fs.Var(&tags, "tag", "tag; repeatable")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	cases, err := corpus.LoadForGoldenRecording(g.root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	cases, err = corpus.Filter(cases, selectors, tags)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	var defuddle participant.Target
	for _, target := range targets(g) {
		if target.ID == participant.Defuddle {
			defuddle = target
			break
		}
	}
	if !defuddle.Available {
		fmt.Fprintf(stderr, "Defuddle Golden Oracle unavailable: %s\n", defuddle.Reason)
		return exitUnavailable
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	err = participant.Probe(ctx, defuddle)
	cancel()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUnavailable
	}
	for _, c := range cases {
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		native := participant.Run(ctx, defuddle, c.Fixture, c.BaseURL, "")
		cancel()
		if native.Status == "failure" {
			fmt.Fprintf(stderr, "%s: Defuddle Golden extraction failed: %s\n", c.ID, native.Stderr)
			return exitFailed
		}
		golden, err := corpus.NewGolden(defuddleGoldenSource, c, native.Status, native.Content, native.Metadata)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
		if err := corpus.WriteGolden(corpus.GoldenPath(g.root, c.ID), golden, *force); err != nil {
			fmt.Fprintln(stderr, err)
			return exitUsage
		}
		fmt.Fprintf(stdout, "recorded %s status=%s sha256=%s\n", c.ID, native.Status, golden.ContentSHA256[:12])
	}
	return exitOK
}

func runVersion(args []string, g globals, build BuildInfo, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("readabilities-suite version", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", g.json, "emit JSON")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *jsonOutput {
		_ = json.NewEncoder(stdout).Encode(map[string]any{"schema_version": 1, "name": "readabilities-suite", "version": build.Version})
	} else {
		fmt.Fprintln(stdout, build.Version)
	}
	return exitOK
}

func captureLive(urls []string, outputDir string, timeout time.Duration) ([]corpus.Case, func(), error) {
	cleanup := func() {}
	if outputDir == "" {
		dir, err := os.MkdirTemp("", "readabilities-suite-live-")
		if err != nil {
			return nil, nil, err
		}
		outputDir = dir
		cleanup = func() { _ = os.RemoveAll(dir) }
	} else if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, nil, err
	}
	client := &http.Client{Timeout: timeout}
	var cases []corpus.Case
	for i, rawURL := range urls {
		resp, err := fetchLiveWithRetry(client, rawURL)
		if err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("capture %s: %w", rawURL, err)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			cleanup()
			return nil, nil, fmt.Errorf("capture %s: HTTP %s", rawURL, resp.Status)
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024+1))
		resp.Body.Close()
		if err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("capture %s: %w", rawURL, err)
		}
		if len(body) > 8*1024*1024 {
			cleanup()
			return nil, nil, fmt.Errorf("capture %s exceeded 8 MiB", rawURL)
		}
		path := filepath.Join(outputDir, fmt.Sprintf("LIVE-%03d.html", i+1))
		if err := os.WriteFile(path, body, 0o644); err != nil {
			cleanup()
			return nil, nil, err
		}
		cases = append(cases, corpus.Case{ID: fmt.Sprintf("LIVE-%03d", i+1), Name: rawURL, Kind: "live", Fixture: path, BaseURL: rawURL, ExpectedStatuses: []string{"success", "no_content"}, Tags: []string{"live"}})
	}
	return cases, cleanup, nil
}

func fetchLiveWithRetry(client *http.Client, rawURL string) (*http.Response, error) {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "readabilities-suite/0.1 common-snapshot")
		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if attempt < 3 {
			time.Sleep(time.Duration(attempt) * 250 * time.Millisecond)
		}
	}
	return nil, lastErr
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, `readabilities-suite [--rust PATH] [--defuddle JS] [--python DIR] <command>

Commands:
  doctor [--json]                resolve and probe all three participants
  list [--json] [--case ID]      list stable READ-NNN cases
  run [--profile offline|live]   compare normalized facts and security
  record-golden [--force]        freeze the selected Defuddle Oracle for offline cases
  version [--json]               print suite version

Offline is the release gate. Live captures one common snapshot per --url and is discovery-only.`)
}
