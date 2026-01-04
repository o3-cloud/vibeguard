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
//
// Pattern Compilation:
//   - Each pattern is compiled immediately to detect syntax errors early.
//   - This "fail-fast" approach catches invalid patterns during matcher creation rather than during matching.
//   - If any pattern fails to compile, creation returns an error without creating the matcher.
//
// Pattern Syntax:
//   - Supports grok built-in patterns: %{PATTERN_NAME:capture_name}
//   - Supports custom regex with named captures: (?P<capture_name>regex)
//   - Can mix both syntaxes in the same pattern string
//   - Special regex characters must be properly escaped
//
// Error Messages:
//   - If a pattern fails to compile, the error message includes the invalid pattern string.
//   - This helps users identify exactly which pattern has a syntax error.
//
// Examples:
//   - Valid: "%{NUMBER:count} tests"
//   - Valid: "coverage:\s+%{NUMBER:coverage}%"
//   - Valid: "(?P<status>\w+) test"
//   - Invalid: "%{NONEXISTENT_PATTERN:val}" (unknown pattern name)
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
//
// Pattern Matching Behavior:
//   - All patterns are applied sequentially and independently to the input.
//   - If multiple patterns capture the same field name, the later pattern's value overrides earlier ones.
//   - If a pattern doesn't match, its fields are simply not included in the result (no error).
//   - Unmatched patterns do not generate errors; only pattern compilation errors do.
//
// Pattern Syntax Support:
//   - Built-in patterns: %{NUMBER:name}, %{INT:name}, %{WORD:name}, %{IP:name}, %{IPV6:name}, %{UUID:name}, etc.
//   - Custom regex: (?P<name>pattern) for named capture groups
//   - Mixed patterns: Combine built-in and custom patterns in the same string
//   - Special characters must be escaped: use \( for literal parentheses, \[ for brackets, etc.
//
// Examples:
//   - Simple: "%{NUMBER:coverage}%"
//   - With context: "coverage: %{NUMBER:coverage}%"
//   - Custom regex: "(?P<count>[0-9]+) tests?"
//   - Mixed: "total:.*\(statements\)\s+%{NUMBER:coverage}%"
func (m *Matcher) Match(input string) (map[string]string, error) {
	result := make(map[string]string)

	if len(m.compiled) == 0 {
		return result, nil
	}

	for _, g := range m.compiled {
		values, _ := g.ParseString(input)
		// Merge extracted values into result
		// Later patterns can override values from earlier patterns
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
