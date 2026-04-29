package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"ContextForge/internal/analyzer"
	"ContextForge/internal/exporter"
	"ContextForge/internal/storage"
	"ContextForge/internal/template"
	"ContextForge/internal/tokenizer"
)

// App holds application-level dependencies exposed to the Vue frontend via Wails bindings.
type App struct {
	ctx context.Context
	db  *storage.DB
}

// NewApp creates a new App with uninitialised dependencies (startup wires them).
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. Opens the database and seeds built-in templates.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir, err := os.UserConfigDir()
	if err != nil {
		runtime.LogErrorf(ctx, "app: get config dir: %v", err)
		return
	}

	dbDir := filepath.Join(dataDir, "ContextForge")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		runtime.LogErrorf(ctx, "app: create data dir: %v", err)
		return
	}

	db, err := storage.New(filepath.Join(dbDir, "contextforge.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "app: init db: %v", err)
		return
	}
	a.db = db
	if err := template.SeedBuiltIns(db); err != nil {
		runtime.LogWarningf(ctx, "app: seed built-ins: %v", err)
	}
	runtime.LogInfo(ctx, "app: database ready")
}

// ---- Project analysis ----

// AnalyzeProject scans the given directory, detects its tech stack, and
// persists (or updates) the project in the database.
func (a *App) AnalyzeProject(path string) (analyzer.DetectionResult, error) {
	if a.db == nil {
		return analyzer.DetectionResult{}, fmt.Errorf("app: database not initialised")
	}

	an := analyzer.New()
	result, err := an.Analyze(path)
	if err != nil {
		return analyzer.DetectionResult{}, fmt.Errorf("app: analyze project: %w", err)
	}

	// Upsert project record
	existing, err := a.db.GetProjectByPath(path)
	if err != nil {
		// Project not found — create it
		name := filepath.Base(path)
		proj := &storage.Project{Name: name, Path: path, DetectedStack: result.Stack}
		if createErr := a.db.CreateProject(proj); createErr != nil {
			runtime.LogWarningf(a.ctx, "app: could not persist project: %v", createErr)
		}
	} else {
		existing.DetectedStack = result.Stack
		if updateErr := a.db.UpdateProject(existing); updateErr != nil {
			runtime.LogWarningf(a.ctx, "app: could not update project: %v", updateErr)
		}
	}

	return result, nil
}

// ---- Template operations ----

// GetTemplates returns all templates including their rules.
func (a *App) GetTemplates() ([]storage.Template, error) {
	if a.db == nil {
		return nil, fmt.Errorf("app: database not initialised")
	}
	var templates []storage.Template
	if err := a.db.Conn().Preload("Rules").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("app: get templates: %w", err)
	}
	return templates, nil
}

// GetTemplatesByStack returns templates that match any of the given stack tags.
func (a *App) GetTemplatesByStack(stack []string) ([]storage.Template, error) {
	all, err := a.GetTemplates()
	if err != nil {
		return nil, err
	}
	matched := make([]storage.Template, 0)
	for _, t := range all {
		for _, tag := range t.Stack {
			for _, s := range stack {
				if tag == s {
					matched = append(matched, t)
					goto next
				}
			}
		}
	next:
	}
	return matched, nil
}

// GetProjectByPath returns the persisted project record for the given path.
func (a *App) GetProjectByPath(path string) (storage.Project, error) {
	if a.db == nil {
		return storage.Project{}, fmt.Errorf("app: database not initialised")
	}
	p, err := a.db.GetProjectByPath(path)
	if err != nil {
		return storage.Project{}, fmt.Errorf("app: get project by path: %w", err)
	}
	return *p, nil
}

// SaveTemplate creates or updates a template (upsert by ID).
func (a *App) SaveTemplate(template storage.Template) error {
	if a.db == nil {
		return fmt.Errorf("app: database not initialised")
	}
	if err := a.db.Conn().Save(&template).Error; err != nil {
		return fmt.Errorf("app: save template: %w", err)
	}
	return nil
}

// ---- Token counting ----

// CountTokens estimates the token count for the given content (1 token ≈ 4 chars).
func (a *App) CountTokens(content string) int {
	tk := tokenizer.New()
	return tk.Count(content)
}

// CheckTokenWarnings returns warnings for all models where content exceeds safe limits.
func (a *App) CheckTokenWarnings(content string) []tokenizer.Warning {
	tk := tokenizer.New()
	return tk.Check(content)
}

// ---- Export ----

// resolveExporter returns the Exporter for the given target identifier.
func resolveExporter(target string) (exporter.Exporter, error) {
	switch target {
	case "claude":
		return exporter.NewClaudeExporter(), nil
	case "copilot":
		return exporter.NewCopilotExporter(), nil
	default:
		return nil, fmt.Errorf("app: unknown export target %q", target)
	}
}

// PreviewExport renders the export content for a project without writing to disk.
func (a *App) PreviewExport(projectID uint, target string) (string, error) {
	if a.db == nil {
		return "", fmt.Errorf("app: database not initialised")
	}
	exp, err := resolveExporter(target)
	if err != nil {
		return "", err
	}
	project, err := a.db.GetProject(projectID)
	if err != nil {
		return "", fmt.Errorf("app: get project: %w", err)
	}
	templates, err := a.db.GetTemplatesByStack(project.DetectedStack)
	if err != nil {
		return "", fmt.Errorf("app: get templates: %w", err)
	}
	var rules []storage.Rule
	for _, t := range templates {
		rules = append(rules, t.Rules...)
	}
	return exp.Export(*project, rules)
}

// ExportToFile renders and writes the export file into outputPath directory.
// The .github/ subdirectory is created automatically if needed.
func (a *App) ExportToFile(projectID uint, target string, outputPath string) error {
	content, err := a.PreviewExport(projectID, target)
	if err != nil {
		return err
	}
	exp, err := resolveExporter(target)
	if err != nil {
		return err
	}
	fullPath := filepath.Join(outputPath, exp.GetFileName())
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("app: create output dir: %w", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("app: write export file: %w", err)
	}
	return nil
}
