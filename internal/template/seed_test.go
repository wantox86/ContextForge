package template

import (
	"path/filepath"
	"testing"

	"ContextForge/internal/storage"
)

func TestSeedBuiltIns(t *testing.T) {
	db, err := storage.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("create db: %v", err)
	}

	if err := SeedBuiltIns(db); err != nil {
		t.Fatalf("seed: %v", err)
	}

	templates, err := db.ListTemplates()
	if err != nil {
		t.Fatalf("list templates: %v", err)
	}

	if len(templates) != 5 {
		t.Errorf("expected 5 built-in templates, got %d", len(templates))
	}

	for _, tpl := range templates {
		if !tpl.IsBuiltIn {
			t.Errorf("template %q should be marked IsBuiltIn", tpl.Name)
		}
		if len(tpl.Rules) == 0 {
			t.Errorf("template %q has no rules", tpl.Name)
		}
		if len(tpl.Stack) == 0 {
			t.Errorf("template %q has no stack tags", tpl.Name)
		}
	}
}

func TestSeedBuiltIns_Idempotent(t *testing.T) {
	db, err := storage.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("create db: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := SeedBuiltIns(db); err != nil {
			t.Fatalf("seed run %d: %v", i+1, err)
		}
	}

	templates, err := db.ListTemplates()
	if err != nil {
		t.Fatalf("list templates: %v", err)
	}
	if len(templates) != 5 {
		t.Errorf("idempotent seed: expected 5 templates, got %d", len(templates))
	}
}
