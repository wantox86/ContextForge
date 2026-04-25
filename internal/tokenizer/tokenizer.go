// Package tokenizer estimates token counts and warns when model limits are exceeded.
package tokenizer

// TokenLimit defines a model's context and warning threshold.
type TokenLimit struct {
	Model  string
	Max    int
	WarnAt int
}

// ModelLimits lists supported models with their token boundaries.
// Estimate: 1 token ≈ 4 characters.
var ModelLimits = []TokenLimit{
	{Model: "claude-sonnet", Max: 200000, WarnAt: 8000},
	{Model: "github-copilot", Max: 8000, WarnAt: 6000},
	{Model: "cursor", Max: 8000, WarnAt: 6000},
	{Model: "ollama-8k", Max: 8192, WarnAt: 6000},
	{Model: "ollama-32k", Max: 32768, WarnAt: 24000},
}

// Tokenizer provides token estimation utilities.
type Tokenizer struct{}

// New creates a new Tokenizer.
func New() *Tokenizer { return &Tokenizer{} }

// Count returns the estimated token count for the given content.
func (t *Tokenizer) Count(content string) int {
	return (len(content) + 3) / 4
}

// Warning describes a token limit breach for a specific model.
type Warning struct {
	Model   string
	Tokens  int
	WarnAt  int
	Max     int
	Exceeds bool // true if over Max, false if over WarnAt only
}

// Check returns warnings for all models where the content exceeds the WarnAt threshold.
func (t *Tokenizer) Check(content string) []Warning {
	tokens := t.Count(content)
	warnings := make([]Warning, 0)
	for _, limit := range ModelLimits {
		if tokens >= limit.WarnAt {
			warnings = append(warnings, Warning{
				Model:   limit.Model,
				Tokens:  tokens,
				WarnAt:  limit.WarnAt,
				Max:     limit.Max,
				Exceeds: tokens > limit.Max,
			})
		}
	}
	return warnings
}

// CheckForModel returns the warning for a specific model, or nil if within safe limits.
func (t *Tokenizer) CheckForModel(content, model string) *Warning {
	tokens := t.Count(content)
	for _, limit := range ModelLimits {
		if limit.Model == model && tokens >= limit.WarnAt {
			w := Warning{
				Model:   limit.Model,
				Tokens:  tokens,
				WarnAt:  limit.WarnAt,
				Max:     limit.Max,
				Exceeds: tokens > limit.Max,
			}
			return &w
		}
	}
	return nil
}
