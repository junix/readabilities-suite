package suite

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/junix/readabilities-suite/internal/corpus"
	"github.com/junix/readabilities-suite/internal/participant"
)

func TestEvaluateKeepsRecallNoiseStructureAndSecuritySeparate(t *testing.T) {
	caseDef := corpus.Case{
		ID:                "READ-TEST",
		ExpectedStatuses:  []string{"success"},
		RequiredText:      []string{"required fact"},
		ForbiddenText:     []string{"forbidden noise"},
		RequiredStructure: []string{"heading"},
	}
	native := participant.NativeResult{
		Participant:   participant.Rust,
		Status:        "success",
		Content:       "# Title\n\nrequired fact\n\n<script>bad()</script>",
		ContentFormat: "markdown",
	}
	result := evaluate(caseDef, native)
	if result.Status != "fail" {
		t.Fatal("security violation must fail the case")
	}
	if passCount(result.Required) != 1 || passCount(result.Forbidden) != 1 || passCount(result.Structure) != 1 {
		t.Fatalf("independent dimensions were conflated: %#v", result)
	}
	if len(result.Security) != 1 || result.Security[0] != "script_element" {
		t.Fatalf("missing security evidence: %#v", result.Security)
	}
}

func TestBestEvidenceRequiresRustToHaveNoFailedCases(t *testing.T) {
	metrics := []Metrics{
		{Participant: participant.Defuddle, RequiredFound: 2, RequiredTotal: 2, ForbiddenRejected: 2, ForbiddenTotal: 2},
		{Participant: participant.Rust, Fail: 1, RequiredFound: 2, RequiredTotal: 2, ForbiddenRejected: 2, ForbiddenTotal: 2},
	}
	if bestEvidence(metrics).Eligible {
		t.Fatal("Rust with a failed case cannot pass the best-evidence gate")
	}
}

func TestJSONReportOmitsFullParticipantBodies(t *testing.T) {
	report := Report{Cases: []CaseResult{{Native: participant.NativeResult{
		Content:        "UNBOUNDED-PRIVATE-CONTENT",
		ContentPreview: "bounded preview",
	}}}}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "UNBOUNDED-PRIVATE-CONTENT") {
		t.Fatalf("full participant body leaked into report: %s", data)
	}
	if !strings.Contains(string(data), "bounded preview") {
		t.Fatalf("bounded evidence was omitted: %s", data)
	}
}

func TestSecurityCheckDoesNotConfuseProseWithAURI(t *testing.T) {
	safe := `<p>JavaScript: this is literal prose.</p><a href="https://example.test">safe</a>`
	if got := securityViolations(safe); len(got) != 0 {
		t.Fatalf("safe prose was treated as active content: %v", got)
	}
	dangerous := `[unsafe](javascript:steal())`
	if got := securityViolations(dangerous); len(got) != 1 || got[0] != "dangerous_markdown_uri" {
		t.Fatalf("dangerous URI was not detected: %v", got)
	}
}

func TestGoldenOracleFindsMissingReferenceContent(t *testing.T) {
	caseDef := corpus.Case{
		ID:               "READ-TEST",
		ExpectedStatuses: []string{"success"},
		Golden: &corpus.Golden{
			Source:  "defuddle/0.19.2-markdown-v1",
			Status:  "success",
			Content: "A retained reference paragraph has several important words.",
		},
	}
	native := participant.NativeResult{
		Participant:   participant.Rust,
		Status:        "success",
		Content:       "A retained paragraph.",
		ContentFormat: "markdown",
	}
	result := evaluate(caseDef, native)
	if result.Status != "fail" || result.Golden == nil || result.Golden.Pass {
		t.Fatalf("Golden omission must fail with evidence: %#v", result)
	}
	if result.Golden.TokenCoverage >= minimumGoldenTokenCoverage || len(result.Golden.MissingTokens) == 0 {
		t.Fatalf("missing Golden evidence: %#v", result.Golden)
	}
}

func TestLiveOracleComparisonUsesTheSameDefuddleSnapshot(t *testing.T) {
	cases := []corpus.Case{{ID: "LIVE-001", ExpectedStatuses: []string{"success", "no_content"}}}
	results := []CaseResult{
		{
			CaseID: "LIVE-001", Participant: participant.Defuddle, Status: "pass",
			Native: participant.NativeResult{Participant: participant.Defuddle, Status: "no_content"},
		},
		{
			CaseID: "LIVE-001", Participant: participant.Rust, Status: "pass",
			Native: participant.NativeResult{Participant: participant.Rust, Status: "success", Content: "link index"},
		},
	}
	attachLiveGolden(cases, results)
	if results[1].Status != "fail" || results[1].Golden == nil || results[1].Golden.Pass {
		t.Fatalf("live Defuddle no_content must reject Rust false positive: %#v", results[1])
	}
}
