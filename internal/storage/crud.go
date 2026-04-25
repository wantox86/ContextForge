// Package storage provides CRUD operations for Template, Rule, and Project models.
package storage

import "fmt"

// ---- Template CRUD ----

// CreateTemplate inserts a new template with its rules.
func (db *DB) CreateTemplate(t *Template) error {
	if err := db.conn.Create(t).Error; err != nil {
		return fmt.Errorf("storage: create template: %w", err)
	}
	return nil
}

// GetTemplate returns a template by ID, preloading its rules.
func (db *DB) GetTemplate(id uint) (*Template, error) {
	var t Template
	if err := db.conn.Preload("Rules").First(&t, id).Error; err != nil {
		return nil, fmt.Errorf("storage: get template %d: %w", id, err)
	}
	return &t, nil
}

// ListTemplates returns all templates with their rules.
func (db *DB) ListTemplates() ([]Template, error) {
	var templates []Template
	if err := db.conn.Preload("Rules").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("storage: list templates: %w", err)
	}
	return templates, nil
}

// UpdateTemplate saves changes to an existing template.
func (db *DB) UpdateTemplate(t *Template) error {
	if err := db.conn.Save(t).Error; err != nil {
		return fmt.Errorf("storage: update template %d: %w", t.ID, err)
	}
	return nil
}

// DeleteTemplate removes a template and its associated rules.
func (db *DB) DeleteTemplate(id uint) error {
	if err := db.conn.Where("template_id = ?", id).Delete(&Rule{}).Error; err != nil {
		return fmt.Errorf("storage: delete rules for template %d: %w", id, err)
	}
	if err := db.conn.Delete(&Template{}, id).Error; err != nil {
		return fmt.Errorf("storage: delete template %d: %w", id, err)
	}
	return nil
}

// ---- Rule CRUD ----

// UpdateRule saves changes to an existing rule.
func (db *DB) UpdateRule(r *Rule) error {
	if err := db.conn.Save(r).Error; err != nil {
		return fmt.Errorf("storage: update rule %d: %w", r.ID, err)
	}
	return nil
}

// DeleteRule removes a rule by ID.
func (db *DB) DeleteRule(id uint) error {
	if err := db.conn.Delete(&Rule{}, id).Error; err != nil {
		return fmt.Errorf("storage: delete rule %d: %w", id, err)
	}
	return nil
}

// ---- Project CRUD ----

// CreateProject inserts a new project.
func (db *DB) CreateProject(p *Project) error {
	if err := db.conn.Create(p).Error; err != nil {
		return fmt.Errorf("storage: create project: %w", err)
	}
	return nil
}

// GetProject returns a project by ID.
func (db *DB) GetProject(id uint) (*Project, error) {
	var p Project
	if err := db.conn.First(&p, id).Error; err != nil {
		return nil, fmt.Errorf("storage: get project %d: %w", id, err)
	}
	return &p, nil
}

// GetTemplatesByStack is also used internally; expose it on DB for reuse.
func (db *DB) GetTemplatesByStack(stack []string) ([]Template, error) {
	all, err := db.ListTemplates()
	if err != nil {
		return nil, err
	}
	matched := make([]Template, 0)
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
func (db *DB) GetProjectByPath(path string) (*Project, error) {
	var p Project
	if err := db.conn.Where("path = ?", path).First(&p).Error; err != nil {
		return nil, fmt.Errorf("storage: get project by path: %w", err)
	}
	return &p, nil
}

// ListProjects returns all projects.
func (db *DB) ListProjects() ([]Project, error) {
	var projects []Project
	if err := db.conn.Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("storage: list projects: %w", err)
	}
	return projects, nil
}

// UpdateProject saves changes to an existing project.
func (db *DB) UpdateProject(p *Project) error {
	if err := db.conn.Save(p).Error; err != nil {
		return fmt.Errorf("storage: update project %d: %w", p.ID, err)
	}
	return nil
}

// DeleteProject removes a project by ID.
func (db *DB) DeleteProject(id uint) error {
	if err := db.conn.Delete(&Project{}, id).Error; err != nil {
		return fmt.Errorf("storage: delete project %d: %w", id, err)
	}
	return nil
}
