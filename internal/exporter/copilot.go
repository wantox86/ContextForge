package exporter

import (
	"fmt"

	"ContextForge/internal/storage"
)

// CopilotExporter exports rules as a .github/copilot-instructions.md file.
type CopilotExporter struct{}

// NewCopilotExporter creates a CopilotExporter.
func NewCopilotExporter() *CopilotExporter { return &CopilotExporter{} }

// GetFileName returns ".github/copilot-instructions.md".
func (e *CopilotExporter) GetFileName() string { return ".github/copilot-instructions.md" }

// GetTarget returns "copilot".
func (e *CopilotExporter) GetTarget() string { return "copilot" }

// Export renders copilot-instructions.md content for the given project and rules.
func (e *CopilotExporter) Export(project storage.Project, rules []storage.Rule) (string, error) {
	header, err := renderHeader("GitHub Copilot — Coding Instructions", project)
	if err != nil {
		return "", fmt.Errorf("copilot exporter: %w", err)
	}
	body := renderRules(project, rules)
	return header + body, nil
}
