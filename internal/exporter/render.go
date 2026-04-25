package exporter

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"ContextForge/internal/storage"
)

// categoryOrder defines display order for rule categories.
var categoryOrder = []string{
	"architecture", "coding_style", "naming", "testing", "security",
}

// renderRules groups enabled rules by category and renders them as Markdown sections.
func renderRules(project storage.Project, rules []storage.Rule) string {
	// Filter enabled rules
	enabled := make([]storage.Rule, 0, len(rules))
	for _, r := range rules {
		if r.Enabled {
			enabled = append(enabled, r)
		}
	}

	// Group by category
	grouped := map[string][]storage.Rule{}
	for _, r := range enabled {
		grouped[r.Category] = append(grouped[r.Category], r)
	}

	// Sort rules within each category by priority
	for cat := range grouped {
		sort.Slice(grouped[cat], func(i, j int) bool {
			return grouped[cat][i].Priority < grouped[cat][j].Priority
		})
	}

	var sb strings.Builder

	// Render in defined category order, then any remaining categories
	rendered := map[string]bool{}
	ordered := append([]string{}, categoryOrder...)
	for cat := range grouped {
		found := false
		for _, o := range categoryOrder {
			if o == cat {
				found = true
				break
			}
		}
		if !found {
			ordered = append(ordered, cat)
		}
	}

	for _, cat := range ordered {
		rules, ok := grouped[cat]
		if !ok {
			continue
		}
		if rendered[cat] {
			continue
		}
		rendered[cat] = true

		title := strings.ReplaceAll(strings.Title(strings.ReplaceAll(cat, "_", " ")), " ", " ")
		sb.WriteString(fmt.Sprintf("## %s\n\n", title))
		for _, r := range rules {
			sb.WriteString(fmt.Sprintf("### %s\n\n%s\n\n", r.Title, r.Content))
		}
	}

	return sb.String()
}

var headerTmpl = template.Must(template.New("header").Parse(
	`# {{.Title}}

> Project: **{{.ProjectName}}** | Stack: {{.Stack}}

`))

type headerData struct {
	Title       string
	ProjectName string
	Stack       string
}

// renderHeader produces the top-level header block.
func renderHeader(title string, project storage.Project) (string, error) {
	var buf bytes.Buffer
	if err := headerTmpl.Execute(&buf, headerData{
		Title:       title,
		ProjectName: project.Name,
		Stack:       strings.Join(project.DetectedStack, ", "),
	}); err != nil {
		return "", fmt.Errorf("exporter: render header: %w", err)
	}
	return buf.String(), nil
}
