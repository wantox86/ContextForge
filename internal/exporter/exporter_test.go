package exporter

import (
	"strings"
	"testing"

	"ContextForge/internal/storage"
)

var testProject = storage.Project{
	Name:          "MyApp",
	Path:          "/tmp/myapp",
	DetectedStack: []string{"golang"},
}

var testRules = []storage.Rule{
	{ID: 1, Category: "coding_style", Title: "Use gofmt", Content: "Always run gofmt before committing.", Priority: 1, Enabled: true},
	{ID: 2, Category: "testing", Title: "Write table tests", Content: "Prefer table-driven tests.", Priority: 1, Enabled: true},
	{ID: 3, Category: "security", Title: "No hardcoded secrets", Content: "Never hardcode API keys.", Priority: 1, Enabled: false},
}

func TestClaudeExporter(t *testing.T) {
	e := NewClaudeExporter()

	if e.GetTarget() != "claude" {
		t.Errorf("unexpected target: %s", e.GetTarget())
	}
	if e.GetFileName() != "CLAUDE.md" {
		t.Errorf("unexpected filename: %s", e.GetFileName())
	}

	out, err := e.Export(testProject, testRules)
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "Claude AI") {
		t.Error("missing Claude header")
	}
	if !strings.Contains(out, "Use gofmt") {
		t.Error("missing enabled rule")
	}
	// disabled rule should not appear
	if strings.Contains(out, "No hardcoded secrets") {
		t.Error("disabled rule should not appear in output")
	}
}

func TestCopilotExporter(t *testing.T) {
	e := NewCopilotExporter()

	if e.GetTarget() != "copilot" {
		t.Errorf("unexpected target: %s", e.GetTarget())
	}
	if e.GetFileName() != ".github/copilot-instructions.md" {
		t.Errorf("unexpected filename: %s", e.GetFileName())
	}

	out, err := e.Export(testProject, testRules)
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "GitHub Copilot") {
		t.Error("missing Copilot header")
	}
	if !strings.Contains(out, "Write table tests") {
		t.Error("missing testing rule")
	}
}

func TestExporterInterface(t *testing.T) {
	exporters := []Exporter{
		NewClaudeExporter(),
		NewCopilotExporter(),
	}
	for _, e := range exporters {
		if e.GetTarget() == "" {
			t.Errorf("empty target for %T", e)
		}
		if e.GetFileName() == "" {
			t.Errorf("empty filename for %T", e)
		}
		_, err := e.Export(testProject, testRules)
		if err != nil {
			t.Errorf("export failed for %T: %v", e, err)
		}
	}
}
