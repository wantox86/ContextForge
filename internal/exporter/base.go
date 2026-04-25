// Package exporter defines the Exporter interface and shared rendering helpers.
package exporter

import "ContextForge/internal/storage"

// Exporter renders project rules into a target-specific instruction file.
type Exporter interface {
	// Export renders the content string for the given project and rules.
	Export(project storage.Project, rules []storage.Rule) (string, error)
	// GetFileName returns the output filename (e.g. "CLAUDE.md").
	GetFileName() string
	// GetTarget returns the target identifier (e.g. "claude", "copilot").
	GetTarget() string
}
