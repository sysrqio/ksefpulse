package reconcile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeMissingUPO(t *testing.T) {
	st := &StateFile{
		Version: 1,
		Invoices: []Invoice{
			{ID: "inv-001", Status: "submitted"},
		},
	}
	rep := Analyze(st)
	if rep.ReadinessScore >= 70 {
		t.Fatalf("expected low score, got %d", rep.ReadinessScore)
	}
	if len(rep.Findings) != 1 || rep.Findings[0].Code != "MISSING_UPO" {
		t.Fatalf("unexpected findings: %+v", rep.Findings)
	}
}

func TestLoadFixtureOrphan(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "state-orphan.json")
	st, err := LoadState(path)
	if err != nil {
		t.Fatal(err)
	}
	rep := Analyze(st)
	if rep.ReadinessScore > 50 {
		t.Fatalf("expected critical findings, score=%d", rep.ReadinessScore)
	}
}

func TestWriteReport(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "report.json")
	rep := Analyze(&StateFile{Version: 1, Invoices: []Invoice{}})
	if err := WriteReport(out, rep); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}
