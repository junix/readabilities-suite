// Package corpus loads the suite-owned fixtures and fact oracles.
package corpus

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const Version = "readabilities-suite.corpus.v3"
const GoldenSchemaVersion = 1

type Case struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Kind                   string   `json:"kind"`
	Fixture                string   `json:"fixture"`
	BaseURL                string   `json:"base_url,omitempty"`
	ExpectedStatuses       []string `json:"expected_statuses"`
	RequiredText           []string `json:"required_text,omitempty"`
	ForbiddenText          []string `json:"forbidden_text,omitempty"`
	RequiredStructure      []string `json:"required_structure,omitempty"`
	RequiredMarkdown       []string `json:"required_markdown,omitempty"`
	ForbiddenMarkdown      []string `json:"forbidden_markdown,omitempty"`
	MetadataRequired       []string `json:"metadata_required,omitempty"`
	ExpectedRustSiteConfig *string  `json:"expected_rust_site_config,omitempty"`
	RustSiteConfigFile     string   `json:"rust_site_config_file,omitempty"`
	Tags                   []string `json:"tags,omitempty"`
	Golden                 *Golden  `json:"-"`
}

type File struct {
	Version string `json:"version"`
	Cases   []Case `json:"cases"`
}

// Golden is the frozen output from the selected reference extractor.
// Content is intentionally excluded from ordinary suite reports and list
// output; the dedicated corpus artifact holds the reviewed reference body.
type Golden struct {
	SchemaVersion int               `json:"schema_version"`
	Source        string            `json:"source"`
	CaseID        string            `json:"case_id"`
	FixtureSHA256 string            `json:"fixture_sha256"`
	Status        string            `json:"status"`
	Content       string            `json:"content"`
	ContentSHA256 string            `json:"content_sha256"`
	Metadata      map[string]string `json:"metadata"`
}

func Load(root string) ([]Case, error) {
	return load(root, true)
}

// LoadForGoldenRecording loads fixture definitions without loading their old
// Golden artifacts. This lets an explicit reviewed `record-golden --force`
// replace a Golden whose fixture digest is intentionally stale; normal suite
// runs always use Load and retain strict digest validation.
func LoadForGoldenRecording(root string) ([]Case, error) {
	return load(root, false)
}

func load(root string, loadGoldens bool) ([]Case, error) {
	path := filepath.Join(root, "corpus", "cases.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read corpus: %w", err)
	}
	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse corpus: %w", err)
	}
	if file.Version != Version {
		return nil, fmt.Errorf("corpus version %q, want %q", file.Version, Version)
	}
	seen := map[string]bool{}
	for i := range file.Cases {
		c := &file.Cases[i]
		if c.ID == "" || seen[c.ID] {
			return nil, fmt.Errorf("case %d has empty or duplicate id %q", i, c.ID)
		}
		seen[c.ID] = true
		if len(c.ExpectedStatuses) == 0 {
			return nil, fmt.Errorf("case %s has no expected statuses", c.ID)
		}
		if c.Fixture != "" {
			c.Fixture = filepath.Join(root, "corpus", "fixtures", c.Fixture)
			if _, err := os.Stat(c.Fixture); err != nil {
				return nil, fmt.Errorf("case %s fixture: %w", c.ID, err)
			}
		}
		if c.RustSiteConfigFile != "" {
			c.RustSiteConfigFile = filepath.Join(root, "corpus", "site-configs", c.RustSiteConfigFile)
			if _, err := os.Stat(c.RustSiteConfigFile); err != nil {
				return nil, fmt.Errorf("case %s Rust site config: %w", c.ID, err)
			}
		}
		if loadGoldens {
			golden, err := loadGolden(root, *c)
			if err != nil {
				return nil, err
			}
			c.Golden = golden
		}
	}
	return file.Cases, nil
}

func RequireGoldens(cases []Case) error {
	for _, c := range cases {
		if c.Kind == "offline" && c.Golden == nil {
			return fmt.Errorf("case %s has no Defuddle Golden Oracle; run record-golden", c.ID)
		}
	}
	return nil
}

func GoldenPath(root, caseID string) string {
	return filepath.Join(root, "corpus", "golden", caseID+".json")
}

func NewGolden(source string, c Case, status, content string, metadata map[string]string) (Golden, error) {
	fixture, err := os.ReadFile(c.Fixture)
	if err != nil {
		return Golden{}, fmt.Errorf("read fixture %s: %w", c.ID, err)
	}
	fixtureDigest := sha256.Sum256(fixture)
	contentDigest := sha256.Sum256([]byte(content))
	return Golden{
		SchemaVersion: GoldenSchemaVersion,
		Source:        source,
		CaseID:        c.ID,
		FixtureSHA256: fmt.Sprintf("%x", fixtureDigest),
		Status:        status,
		Content:       content,
		ContentSHA256: fmt.Sprintf("%x", contentDigest),
		Metadata:      metadata,
	}, nil
}

func WriteGolden(path string, golden Golden, overwrite bool) error {
	if _, err := os.Stat(path); err == nil && !overwrite {
		return fmt.Errorf("Golden %s exists; rerun with --force after review", path)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	data, err := json.MarshalIndent(golden, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func loadGolden(root string, c Case) (*Golden, error) {
	path := GoldenPath(root, c.ID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("case %s read Golden: %w", c.ID, err)
	}
	var golden Golden
	if err := json.Unmarshal(data, &golden); err != nil {
		return nil, fmt.Errorf("case %s parse Golden: %w", c.ID, err)
	}
	if golden.SchemaVersion != GoldenSchemaVersion || golden.CaseID != c.ID || golden.Source == "" || golden.Status == "" {
		return nil, fmt.Errorf("case %s has invalid Golden metadata", c.ID)
	}
	fixture, err := os.ReadFile(c.Fixture)
	if err != nil {
		return nil, fmt.Errorf("case %s read fixture for Golden: %w", c.ID, err)
	}
	fixtureDigest := sha256.Sum256(fixture)
	if golden.FixtureSHA256 != fmt.Sprintf("%x", fixtureDigest) {
		return nil, fmt.Errorf("case %s Golden does not match fixture; regenerate after review", c.ID)
	}
	contentDigest := sha256.Sum256([]byte(golden.Content))
	if golden.ContentSHA256 != fmt.Sprintf("%x", contentDigest) {
		return nil, fmt.Errorf("case %s Golden content checksum is invalid", c.ID)
	}
	return &golden, nil
}

func Filter(cases []Case, selectors, tags []string) ([]Case, error) {
	selected := make([]Case, 0, len(cases))
	for _, c := range cases {
		if !matches(c, selectors, tags) {
			continue
		}
		selected = append(selected, c)
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("selection matched no cases")
	}
	return selected, nil
}

func matches(c Case, selectors, tags []string) bool {
	if len(selectors) > 0 {
		ok := false
		for _, selector := range selectors {
			if strings.EqualFold(c.ID, selector) || strings.EqualFold(c.Name, selector) {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	for _, wanted := range tags {
		found := false
		for _, tag := range c.Tags {
			if strings.EqualFold(tag, wanted) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
