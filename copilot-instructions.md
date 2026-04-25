# ContextForge — GitHub Copilot Instructions

## Project Summary

ContextForge is an open-source desktop app (Mac + Windows) built with Go + Wails v2 + Vue 3.
It helps developers generate AI coding instruction files (CLAUDE.md, copilot-instructions.md, .cursorrules, AGENTS.md) from visual templates — auto-detecting tech stack and exporting to multiple formats at once.

---

## Tech Stack

- **Backend:** Go 1.22+, Wails v2, GORM, SQLite
- **Frontend:** Vue 3, TypeScript, Pinia, Tailwind CSS, Vite
- **Build:** Wails CLI (`wails dev`, `wails build`)
- **Testing:** Go `testing` package, Vitest (frontend)

---

## Architecture

```
cmd/main.go                  → Wails entry point
app.go                       → Go methods exposed to Vue frontend via Wails bindings
internal/analyzer/           → Auto-detect tech stack from project directory
internal/template/           → Template CRUD, rule rendering, validation
internal/exporter/           → Export rules to CLAUDE.md, copilot-instructions.md, .cursorrules
internal/tokenizer/          → Estimate token count for model compatibility warning
internal/storage/            → SQLite via GORM (models + migrations)
frontend/src/components/     → Vue components (ProjectAnalyzer, RuleBuilder, PreviewPanel, ExportPanel)
frontend/src/stores/         → Pinia stores (project, template, export)
assets/templates/            → Built-in template JSON files
```

---

## Go Conventions

- Use `fmt.Errorf("package: %w", err)` for error wrapping — always wrap with context
- No global state — pass dependencies via struct constructors
- Use interfaces for all exporters (`Exporter` interface in `internal/exporter/base.go`)
- Keep functions under 50 lines — extract helpers if longer
- Use `filepath.Join()` for all file paths — never string concatenation
- All exported functions must have godoc comments
- Use `context.Context` for operations that may be cancelled

## Vue / TypeScript Conventions

- Use Composition API with `<script setup lang="ts">` — no Options API
- All API calls to Go backend go through Pinia stores only — never directly in components
- Type all props with TypeScript interfaces — no `any`
- Component file names: `PascalCase.vue`
- Emit events for parent communication — don't mutate parent state directly
- Use `kebab-case` for custom event names

## General Rules

- All user-facing strings in English
- Error messages must tell the user what to do, not just what failed
- Never hardcode file paths — use constants or config
- Write tests for all `internal/` packages
- No TODO comments in committed code — open a GitHub issue instead

---

## Key Interfaces

```go
// All exporters implement this
type Exporter interface {
    Export(project Project, rules []Rule) (string, error)
    GetFileName() string  // "CLAUDE.md", ".github/copilot-instructions.md", etc
    GetTarget() string    // "claude", "copilot", "cursor", "agents"
}

// Stack detector result
type DetectionResult struct {
    Stack       []string // ["kotlin", "spring-boot", "gradle"]
    Confidence  float64
    Indicators  []string // files that triggered detection
}
```

---

## Data Models (GORM)

```go
type Template struct {
    ID          uint      `gorm:"primarykey"`
    Name        string    `gorm:"not null"`
    Description string
    Stack       []string  `gorm:"serializer:json"`
    Rules       []Rule    `gorm:"foreignKey:TemplateID"`
    IsBuiltIn   bool      `gorm:"default:false"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Rule struct {
    ID         uint   `gorm:"primarykey"`
    TemplateID uint
    Category   string // coding_style | testing | security | naming | architecture
    Title      string `gorm:"not null"`
    Content    string `gorm:"not null"`
    Priority   int
    Enabled    bool   `gorm:"default:true"`
}

type Project struct {
    ID               uint     `gorm:"primarykey"`
    Name             string
    Path             string   `gorm:"unique"`
    DetectedStack    []string `gorm:"serializer:json"`
    AppliedTemplates []uint   `gorm:"serializer:json"`
    CreatedAt        time.Time
}
```

---

## Stack Detection Logic

Detect by scanning for indicator files in project root:

| File/Pattern | Stack Tags |
|---|---|
| `go.mod` | `golang` |
| `pom.xml` | `java`, `maven` |
| `build.gradle.kts` | `kotlin`, `gradle` |
| `build.gradle` | `java`, `gradle` |
| `package.json` | `nodejs` |
| `requirements.txt` | `python` |
| `pyproject.toml` | `python` |
| `Cargo.toml` | `rust` |
| `pubspec.yaml` | `flutter` |
| `android/` folder | `android` |

Secondary: inspect file contents for framework indicators (e.g. `spring-boot-starter` in pom.xml → add `spring-boot` tag).

---

## Built-in Templates (JSON format)

Located in `assets/templates/`. Each file follows this schema:

```json
{
  "name": "Template Name",
  "description": "Short description",
  "stack": ["tag1", "tag2"],
  "rules": [
    {
      "category": "coding_style|testing|security|naming|architecture",
      "title": "Short rule title",
      "content": "Full instruction text for the AI",
      "priority": 1,
      "enabled": true
    }
  ]
}
```

Available built-in templates:
- `spring-boot-kotlin.json`
- `golang-gin.json`
- `python-fastapi.json`
- `vue3.json`
- `android-kotlin.json`

---

## Token Counter

```go
// Rough estimate: 1 token ≈ 4 characters
var ModelLimits = []TokenLimit{
    {"claude-sonnet", 200000, 8000},
    {"github-copilot", 8000, 6000},
    {"cursor", 8000, 6000},
    {"ollama-8k", 8192, 6000},
    {"ollama-32k", 32768, 24000},
}
```

Show warning in UI when estimated token count exceeds `WarnAt` threshold.

---

## Wails Bindings (app.go)

Go methods exposed to Vue frontend. All methods go in `app.go`:

```go
// Analyze a project directory
func (a *App) AnalyzeProject(path string) (DetectionResult, error)

// Template operations
func (a *App) GetTemplates() ([]Template, error)
func (a *App) GetTemplatesByStack(stack []string) ([]Template, error)
func (a *App) SaveTemplate(template Template) error

// Export
func (a *App) ExportToFile(projectID uint, target string, outputPath string) error
func (a *App) PreviewExport(projectID uint, target string) (string, error)

// Token counting
func (a *App) CountTokens(content string) int
```

---

## MVP Scope (v0.1.0)

In scope:
- Project analyzer with auto stack detection
- 5 built-in templates
- Rule enable/disable + inline edit
- Export to CLAUDE.md and copilot-instructions.md
- Token count estimate with warning
- Mac (.app) + Windows (.exe) builds

Out of scope for v0.1.0:
- Community template sharing
- AI-assisted rule suggestion
- Plugin system
- Cursor / AGENTS.md export (v0.2.0)

---

## Do Not

- Do not use global variables or `init()` functions
- Do not call Go backend directly from Vue components — always via Pinia store
- Do not use `os.Exit()` inside library code — return errors instead
- Do not use `interface{}` or `any` without a comment explaining why
- Do not commit `.env` files or secrets of any kind
- Do not use string concatenation for file paths — use `filepath.Join()`
- Do not add dependencies without discussing in an issue first
