package tokenizer

import (
	"testing"
)

func TestCount(t *testing.T) {
	tk := New()

	cases := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"test", 1},
		{"hello world", 3},
		{"abcd", 1},
		{"abcde", 2},
	}

	for _, c := range cases {
		got := tk.Count(c.input)
		if got != c.expected {
			t.Errorf("Count(%q) = %d, want %d", c.input, got, c.expected)
		}
	}
}

func TestCheck_NoWarning(t *testing.T) {
	tk := New()
	// Short content should not trigger any warning
	warnings := tk.Check("short")
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for short content, got %d", len(warnings))
	}
}

func TestCheck_CopilotWarning(t *testing.T) {
	tk := New()
	// 6000 tokens * 4 chars = 24000 chars — exceeds github-copilot WarnAt (6000)
	content := make([]byte, 24000)
	for i := range content {
		content[i] = 'a'
	}

	warnings := tk.Check(string(content))
	found := false
	for _, w := range warnings {
		if w.Model == "github-copilot" {
			found = true
			if w.Exceeds {
				t.Error("should not exceed Max at exactly WarnAt")
			}
		}
	}
	if !found {
		t.Error("expected github-copilot warning")
	}
}

func TestCheckForModel(t *testing.T) {
	tk := New()

	// Under threshold
	w := tk.CheckForModel("hello", "github-copilot")
	if w != nil {
		t.Error("expected nil warning for short content")
	}

	// Over threshold — 8001 tokens * 4 = 32004 chars to exceed Max (8000)
	content := make([]byte, 32004)
	for i := range content {
		content[i] = 'a'
	}
	w = tk.CheckForModel(string(content), "github-copilot")
	if w == nil {
		t.Fatal("expected warning for long content")
	}
	if !w.Exceeds {
		t.Errorf("expected Exceeds=true, got tokens=%d max=%d", w.Tokens, w.Max)
	}
}
