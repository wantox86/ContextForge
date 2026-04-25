package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	db, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	return db
}

func TestTemplate_CRUD(t *testing.T) {
	db := newTestDB(t)

	tpl := &Template{
		Name:        "Test Template",
		Description: "desc",
		Stack:       []string{"golang"},
		Rules: []Rule{
			{Category: "coding_style", Title: "Use gofmt", Content: "Always run gofmt", Priority: 1, Enabled: true},
		},
	}

	if err := db.CreateTemplate(tpl); err != nil {
		t.Fatalf("create: %v", err)
	}
	if tpl.ID == 0 {
		t.Fatal("expected non-zero ID after create")
	}

	got, err := db.GetTemplate(tpl.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != tpl.Name {
		t.Errorf("name mismatch: %s vs %s", got.Name, tpl.Name)
	}
	if len(got.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(got.Rules))
	}

	got.Description = "updated"
	if err := db.UpdateTemplate(got); err != nil {
		t.Fatalf("update: %v", err)
	}

	list, err := db.ListTemplates()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 template, got %d", len(list))
	}

	if err := db.DeleteTemplate(tpl.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = db.GetTemplate(tpl.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestProject_CRUD(t *testing.T) {
	db := newTestDB(t)

	dir, _ := os.MkdirTemp("", "proj")
	defer os.RemoveAll(dir)

	p := &Project{Name: "MyApp", Path: dir, DetectedStack: []string{"golang"}}
	if err := db.CreateProject(p); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := db.GetProjectByPath(dir)
	if err != nil {
		t.Fatalf("get by path: %v", err)
	}
	if got.Name != "MyApp" {
		t.Errorf("name mismatch: %s", got.Name)
	}

	if err := db.DeleteProject(p.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
