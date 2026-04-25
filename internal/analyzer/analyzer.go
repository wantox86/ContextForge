// Package analyzer detects the tech stack of a given project directory.
package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DetectionResult holds the result of a stack detection scan.
type DetectionResult struct {
	Stack      []string `json:"stack"`
	Confidence float64  `json:"confidence"`
	Indicators []string `json:"indicators"` // files/patterns that triggered detection
}

// primaryIndicators maps a filename to its stack tags.
var primaryIndicators = map[string][]string{
	"go.mod":            {"golang"},
	"pom.xml":           {"java", "maven"},
	"build.gradle.kts":  {"kotlin", "gradle"},
	"build.gradle":      {"java", "gradle"},
	"package.json":      {"nodejs"},
	"requirements.txt":  {"python"},
	"pyproject.toml":    {"python"},
	"Cargo.toml":        {"rust"},
	"pubspec.yaml":      {"flutter"},
}

// secondaryIndicators maps a content substring to additional stack tags,
// keyed by the file to inspect.
var secondaryIndicators = map[string]map[string][]string{
	"pom.xml": {
		"spring-boot-starter": {"spring-boot"},
	},
	"build.gradle.kts": {
		"org.springframework.boot": {"spring-boot"},
		"io.ktor":                  {"ktor"},
	},
	"build.gradle": {
		"org.springframework.boot": {"spring-boot"},
	},
	"package.json": {
		`"vue"`:     {"vue"},
		`"react"`:   {"react"},
		`"next"`:    {"nextjs"},
		`"express"`: {"express"},
	},
	"requirements.txt": {
		"fastapi":  {"fastapi"},
		"django":   {"django"},
		"flask":    {"flask"},
	},
	"pyproject.toml": {
		"fastapi":  {"fastapi"},
		"django":   {"django"},
		"flask":    {"flask"},
	},
}

// Analyzer scans a project directory and detects its tech stack.
type Analyzer struct{}

// New creates a new Analyzer.
func New() *Analyzer {
	return &Analyzer{}
}

// Analyze scans the given directory and returns a DetectionResult.
func (a *Analyzer) Analyze(dir string) (DetectionResult, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return DetectionResult{}, fmt.Errorf("analyzer: read dir: %w", err)
	}

	stackSet := map[string]bool{}
	var indicators []string

	// Check for android/ subfolder
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == "android" {
			stackSet["android"] = true
			indicators = append(indicators, "android/")
		}
	}

	// Primary: match filenames
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if tags, ok := primaryIndicators[name]; ok {
			for _, tag := range tags {
				stackSet[tag] = true
			}
			indicators = append(indicators, name)

			// Secondary: inspect file content for framework hints
			if contentMap, ok := secondaryIndicators[name]; ok {
				content, err := os.ReadFile(filepath.Join(dir, name))
				if err == nil {
					contentStr := string(content)
					for needle, extraTags := range contentMap {
						if strings.Contains(contentStr, needle) {
							for _, tag := range extraTags {
								stackSet[tag] = true
							}
							indicators = append(indicators, fmt.Sprintf("%s(%s)", name, needle))
						}
					}
				}
			}
		}
	}

	stack := make([]string, 0, len(stackSet))
	for tag := range stackSet {
		stack = append(stack, tag)
	}

	confidence := 0.0
	if len(stack) > 0 {
		confidence = float64(len(indicators)) / float64(len(indicators)+2)
		if confidence > 1.0 {
			confidence = 1.0
		}
	}

	return DetectionResult{
		Stack:      stack,
		Confidence: confidence,
		Indicators: indicators,
	}, nil
}
