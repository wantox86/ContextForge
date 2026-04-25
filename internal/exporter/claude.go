package exporter

import (
	"fmt"

	"ContextForge/internal/storage"
)

// ClaudeExporter exports rules as a CLAUDE.md file for use with Claude AI.
type ClaudeExporter struct{}

// NewClaudeExporter creates a ClaudeExporter.
func NewClaudeExporter() *ClaudeExporter { return &ClaudeExporter{} }

// GetFileName returns "CLAUDE.md".
func (e *ClaudeExporter) GetFileName() string { return "CLAUDE.md" }

// GetTarget returns "claude".
func (e *ClaudeExporter) GetTarget() string { return "claude" }

// Export renders CLAUDE.md content for the given project and rules.
func (e *ClaudeExporter) Export(project storage.Project, rules []storage.Rule) (string, error) {
	header, err := renderHeader("Claude AI — Coding Instructions", project)
	if err != nil {
		return "", fmt.Errorf("claude exporter: %w", err)
	}
	body := renderRules(project, rules)
	return header + body, nil
}
