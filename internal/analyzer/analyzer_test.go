package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyze_Golang(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/app\n\ngo 1.22\n")

	a := New()
	result, err := a.Analyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsTag(result.Stack, "golang") {
		t.Errorf("expected 'golang' tag, got %v", result.Stack)
	}
}

func TestAnalyze_SpringBoot(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pom.xml"), "<dependency><artifactId>spring-boot-starter-web</artifactId></dependency>")

	a := New()
	result, err := a.Analyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsTag(result.Stack, "spring-boot") {
		t.Errorf("expected 'spring-boot' tag, got %v", result.Stack)
	}
	if !containsTag(result.Stack, "java") {
		t.Errorf("expected 'java' tag, got %v", result.Stack)
	}
}

func TestAnalyze_Vue(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "package.json"), `{"dependencies":{"vue":"^3.0.0"}}`)

	a := New()
	result, err := a.Analyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsTag(result.Stack, "vue") {
		t.Errorf("expected 'vue' tag, got %v", result.Stack)
	}
}

func TestAnalyze_Android(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "android"), 0755); err != nil {
		t.Fatal(err)
	}

	a := New()
	result, err := a.Analyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsTag(result.Stack, "android") {
		t.Errorf("expected 'android' tag, got %v", result.Stack)
	}
}

func TestAnalyze_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	a := New()
	result, err := a.Analyze(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Stack) != 0 {
		t.Errorf("expected empty stack, got %v", result.Stack)
	}
	if result.Confidence != 0.0 {
		t.Errorf("expected 0 confidence, got %v", result.Confidence)
	}
}

// helpers

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func containsTag(stack []string, tag string) bool {
	for _, s := range stack {
		if s == tag {
			return true
		}
	}
	return false
}
