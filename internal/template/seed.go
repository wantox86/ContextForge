// Package template provides loading and seeding of built-in templates.
package template

import (
	"embed"
	"encoding/json"
	"fmt"

	"ContextForge/internal/storage"
)

//go:embed *.json
var builtinFS embed.FS

type jsonRule struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Priority int    `json:"priority"`
	Enabled  bool   `json:"enabled"`
}

type jsonTemplate struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Stack       []string   `json:"stack"`
	Rules       []jsonRule `json:"rules"`
}

// SeedBuiltIns loads built-in templates from embedded JSON files into the database.
// Templates that already exist (by name + IsBuiltIn=true) are skipped.
func SeedBuiltIns(db *storage.DB) error {
	entries, err := builtinFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("template: read embedded dir: %w", err)
	}

	existing, err := db.ListTemplates()
	if err != nil {
		return fmt.Errorf("template: list templates: %w", err)
	}
	builtInNames := map[string]bool{}
	for _, t := range existing {
		if t.IsBuiltIn {
			builtInNames[t.Name] = true
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := builtinFS.ReadFile(entry.Name())
		if err != nil {
			return fmt.Errorf("template: read %s: %w", entry.Name(), err)
		}

		var jt jsonTemplate
		if err := json.Unmarshal(data, &jt); err != nil {
			return fmt.Errorf("template: parse %s: %w", entry.Name(), err)
		}

		if builtInNames[jt.Name] {
			continue // already seeded
		}

		rules := make([]storage.Rule, len(jt.Rules))
		for i, r := range jt.Rules {
			rules[i] = storage.Rule{
				Category: r.Category,
				Title:    r.Title,
				Content:  r.Content,
				Priority: r.Priority,
				Enabled:  r.Enabled,
			}
		}

		tpl := &storage.Template{
			Name:        jt.Name,
			Description: jt.Description,
			Stack:       jt.Stack,
			Rules:       rules,
			IsBuiltIn:   true,
		}
		if err := db.CreateTemplate(tpl); err != nil {
			return fmt.Errorf("template: seed %s: %w", jt.Name, err)
		}
	}

	return nil
}
