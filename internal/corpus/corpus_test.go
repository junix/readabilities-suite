package corpus

import "testing"

func TestFilterRequiresEveryRequestedTag(t *testing.T) {
	cases := []Case{
		{ID: "READ-001", Name: "one", Tags: []string{"security", "article"}},
		{ID: "READ-002", Name: "two", Tags: []string{"article"}},
	}
	got, err := Filter(cases, nil, []string{"article", "security"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "READ-001" {
		t.Fatalf("unexpected selection: %#v", got)
	}
}

func TestFilterRejectsEmptySelection(t *testing.T) {
	_, err := Filter([]Case{{ID: "READ-001", Name: "one"}}, []string{"missing"}, nil)
	if err == nil {
		t.Fatal("expected an empty-selection error")
	}
}
