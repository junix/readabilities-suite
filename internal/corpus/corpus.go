// Package corpus loads the suite-owned fixtures and fact oracles.
package corpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const Version = "readabilities-suite.corpus.v2"

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
}

type File struct {
	Version string `json:"version"`
	Cases   []Case `json:"cases"`
}

func Load(root string) ([]Case, error) {
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
	}
	return file.Cases, nil
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
