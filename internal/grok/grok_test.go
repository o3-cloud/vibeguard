package grok

import (
	"testing"
)

func TestNew_EmptyPatterns(t *testing.T) {
	m, err := New([]string{})
	if err != nil {
		t.Fatalf("New() returned error for empty patterns: %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil matcher")
	}
	if len(m.Patterns()) != 0 {
		t.Errorf("expected 0 patterns, got %d", len(m.Patterns()))
	}
}

func TestNew_NilPatterns(t *testing.T) {
	m, err := New(nil)
	if err != nil {
		t.Fatalf("New() returned error for nil patterns: %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil matcher")
	}
}

func TestNew_ValidPattern(t *testing.T) {
	patterns := []string{"%{IP:ip_addr}"}
	m, err := New(patterns)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil matcher")
	}
	if len(m.Patterns()) != 1 {
		t.Errorf("expected 1 pattern, got %d", len(m.Patterns()))
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	patterns := []string{"%{INVALID_PATTERN_THAT_DOES_NOT_EXIST:value}"}
	_, err := New(patterns)
	if err == nil {
		t.Fatal("New() should return error for invalid pattern")
	}
}

func TestNew_MultiplePatterns(t *testing.T) {
	patterns := []string{
		"%{IP:ip_addr}",
		"%{NUMBER:port}",
	}
	m, err := New(patterns)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if len(m.Patterns()) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(m.Patterns()))
	}
}

func TestMatch_EmptyPatterns(t *testing.T) {
	m, _ := New([]string{})
	result, err := m.Match("any input")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestMatch_EmptyInput(t *testing.T) {
	m, _ := New([]string{"%{IP:ip_addr}"})
	result, err := m.Match("")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	// Empty input won't match the IP pattern, so result should be empty
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %v", result)
	}
}

func TestMatch_IPAddress(t *testing.T) {
	m, _ := New([]string{"%{IP:ip_addr}"})
	result, err := m.Match("Server started at 192.168.1.100")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["ip_addr"] != "192.168.1.100" {
		t.Errorf("expected ip_addr=192.168.1.100, got %q", result["ip_addr"])
	}
}

func TestMatch_MultipleCaptures(t *testing.T) {
	m, _ := New([]string{"%{IP:ip_addr}:%{NUMBER:port}"})
	result, err := m.Match("Listening on 127.0.0.1:8080")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["ip_addr"] != "127.0.0.1" {
		t.Errorf("expected ip_addr=127.0.0.1, got %q", result["ip_addr"])
	}
	if result["port"] != "8080" {
		t.Errorf("expected port=8080, got %q", result["port"])
	}
}

func TestMatch_MultiplePatterns_Merge(t *testing.T) {
	patterns := []string{
		"%{IP:ip_addr}",
		"port (?P<port>[0-9]+)",
	}
	m, _ := New(patterns)
	result, err := m.Match("Server at 10.0.0.1 listening on port 3000")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["ip_addr"] != "10.0.0.1" {
		t.Errorf("expected ip_addr=10.0.0.1, got %q", result["ip_addr"])
	}
	if result["port"] != "3000" {
		t.Errorf("expected port=3000, got %q", result["port"])
	}
}

func TestMatch_NoMatch(t *testing.T) {
	m, _ := New([]string{"%{IP:ip_addr}"})
	result, err := m.Match("no ip address here")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	// No match means empty result (unmatched variables are not included)
	if _, exists := result["ip_addr"]; exists {
		t.Errorf("expected ip_addr to not exist in result, got %v", result)
	}
}

func TestMatch_CustomPattern(t *testing.T) {
	// Test with a custom regex pattern using named captures
	m, err := New([]string{"coverage: (?P<coverage>[0-9.]+)%"})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	result, err := m.Match("Total coverage: 85.5%")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["coverage"] != "85.5" {
		t.Errorf("expected coverage=85.5, got %q", result["coverage"])
	}
}

func TestMatch_NumberPattern(t *testing.T) {
	m, _ := New([]string{"%{NUMBER:count} tests"})
	result, err := m.Match("Ran 42 tests")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["count"] != "42" {
		t.Errorf("expected count=42, got %q", result["count"])
	}
}

func TestMatch_WordPattern(t *testing.T) {
	m, _ := New([]string{"%{WORD:status}"})
	result, err := m.Match("Status: PASSED")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	// WORD matches the first word it finds
	if result["status"] != "Status" {
		t.Errorf("expected status=Status, got %q", result["status"])
	}
}

func TestMatch_MultilineInput(t *testing.T) {
	input := `Running tests...
Test 1: PASS
Test 2: FAIL
Total: 2 tests, 1 passed, 1 failed
Coverage: 75.0%`

	m, _ := New([]string{
		"(?P<total>[0-9]+) tests",
		"(?P<coverage>[0-9.]+)%",
	})
	result, err := m.Match(input)
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["total"] != "2" {
		t.Errorf("expected total=2, got %q", result["total"])
	}
	if result["coverage"] != "75.0" {
		t.Errorf("expected coverage=75.0, got %q", result["coverage"])
	}
}

func TestMatch_RealWorldGoTestOutput(t *testing.T) {
	input := `ok      github.com/example/project        0.125s  coverage: 82.3% of statements`

	m, _ := New([]string{
		"coverage: (?P<coverage>[0-9.]+)%",
	})
	result, err := m.Match(input)
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	if result["coverage"] != "82.3" {
		t.Errorf("expected coverage=82.3, got %q", result["coverage"])
	}
}

func TestMatch_LaterPatternOverridesEarlier(t *testing.T) {
	// If same key is captured by multiple patterns, later one wins
	patterns := []string{
		"(?P<value>first)",
		"(?P<value>second)",
	}
	m, _ := New(patterns)
	result, err := m.Match("first second")
	if err != nil {
		t.Fatalf("Match() returned error: %v", err)
	}
	// Second pattern should override
	if result["value"] != "second" {
		t.Errorf("expected value=second (later pattern override), got %q", result["value"])
	}
}

func TestPatterns_ReturnsConfiguredPatterns(t *testing.T) {
	patterns := []string{"%{IP:ip}", "%{NUMBER:num}"}
	m, _ := New(patterns)

	got := m.Patterns()
	if len(got) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(got))
	}
	if got[0] != patterns[0] || got[1] != patterns[1] {
		t.Errorf("patterns mismatch: got %v, want %v", got, patterns)
	}
}

func TestMatch_CommonLogPatterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		input    string
		expected map[string]string
	}{
		{
			name:     "IPv6 address",
			patterns: []string{"%{IPV6:ipv6}"},
			input:    "Connecting to ::1",
			expected: map[string]string{"ipv6": "::1"},
		},
		{
			name:     "HTTP status code",
			patterns: []string{"HTTP/1.1\" %{NUMBER:status}"},
			input:    `GET /api HTTP/1.1" 200 OK`,
			expected: map[string]string{"status": "200"},
		},
		{
			name:     "UUID pattern",
			patterns: []string{"%{UUID:id}"},
			input:    "Request ID: 550e8400-e29b-41d4-a716-446655440000",
			expected: map[string]string{"id": "550e8400-e29b-41d4-a716-446655440000"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := New(tc.patterns)
			if err != nil {
				t.Fatalf("New() returned error: %v", err)
			}
			result, err := m.Match(tc.input)
			if err != nil {
				t.Fatalf("Match() returned error: %v", err)
			}
			for k, v := range tc.expected {
				if result[k] != v {
					t.Errorf("expected %s=%s, got %q", k, v, result[k])
				}
			}
		})
	}
}
