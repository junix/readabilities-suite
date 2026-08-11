package participant

import "testing"

func TestNormalizedTextOnlyProjectsText(t *testing.T) {
	got := NormalizedText("<h1>A &amp; B</h1><p> one\n two </p>")
	if got != "A & B one two" {
		t.Fatalf("unexpected normalized text: %q", got)
	}
}

func TestNativeResultSummaryDoesNotSerializeBody(t *testing.T) {
	result := NativeResult{Content: "UNBOUNDED-PRIVATE-CONTENT"}
	result.summarizeContent()
	if result.ContentBytes != len(result.Content) || result.ContentSHA256 == "" {
		t.Fatalf("summary is incomplete: %#v", result)
	}
	if result.ContentPreview != result.Content {
		t.Fatalf("short preview changed: %q", result.ContentPreview)
	}
}
