// Package grok provides pattern matching using grok patterns.
// Grok patterns allow extracting structured data from unstructured text
// like logs, command output, etc.
package grok

import (
	"fmt"

	"github.com/elastic/go-grok"
)

// Matcher extracts values from text using grok patterns.
type Matcher struct {
	patterns []string
	compiled []*grok.Grok
}

// New creates a new Matcher with the given patterns.
// Each pattern is compiled eagerly to detect errors early.
func New(patterns []string) (*Matcher, error) {
	if len(patterns) == 0 {
		return &Matcher{
			patterns: patterns,
			compiled: nil,
		}, nil
	}

	compiled := make([]*grok.Grok, 0, len(patterns))
	for _, pattern := range patterns {
		g := grok.New()
		if err := g.Compile(pattern, true); err != nil {
			return nil, fmt.Errorf("failed to compile grok pattern %q: %w", pattern, err)
		}
		compiled = append(compiled, g)
	}

	return &Matcher{
		patterns: patterns,
		compiled: compiled,
	}, nil
}

// Match applies grok patterns to the input and returns extracted values.
// If multiple patterns are provided, all are applied and their results merged.
// Later patterns can override values from earlier patterns.
// Unmatched pattern variables are set to empty strings.
func (m *Matcher) Match(input string) (map[string]string, error) {
	result := make(map[string]string)

	if len(m.compiled) == 0 {
		return result, nil
	}

	for i, g := range m.compiled {
		values, err := g.ParseString(input)
		if err != nil {
			pattern := m.patterns[i]
			// Truncate output for readability if too long
			output := input
			if len(output) > 100 {
				output = output[:97] + "..."
			}
			return nil, fmt.Errorf("grok pattern %d failed to parse\n  pattern: %q\n  output: %q\n  error: %w", i, pattern, output, err)
		}
		// Merge extracted values into result
		for k, v := range values {
			result[k] = v
		}
	}

	return result, nil
}

// Patterns returns the patterns configured for this matcher.
func (m *Matcher) Patterns() []string {
	return m.patterns
}
