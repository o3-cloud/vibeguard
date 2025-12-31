package assist

import (
	"strings"
	"testing"
)

func TestNewValidationGuide(t *testing.T) {
	guide := NewValidationGuide()

	if guide == nil {
		t.Fatal("NewValidationGuide returned nil")
	}

	if guide.YAMLSyntax == "" {
		t.Error("YAMLSyntax should not be empty")
	}
	if guide.CheckStructure == "" {
		t.Error("CheckStructure should not be empty")
	}
	if guide.DependencyRules == "" {
		t.Error("DependencyRules should not be empty")
	}
	if guide.InterpolationRules == "" {
		t.Error("InterpolationRules should not be empty")
	}
	if guide.DoNotList == "" {
		t.Error("DoNotList should not be empty")
	}
}

func TestGetFullGuide(t *testing.T) {
	guide := NewValidationGuide()
	fullGuide := guide.GetFullGuide()

	if fullGuide == "" {
		t.Fatal("GetFullGuide returned empty string")
	}

	// Verify all sections are included
	sections := []string{
		"YAML Syntax Requirements",
		"Check Structure Requirements",
		"Dependency Validation Rules",
		"Variable Interpolation Rules",
		"Explicit DO NOT List",
	}

	for _, section := range sections {
		if !strings.Contains(fullGuide, section) {
			t.Errorf("Full guide missing section: %s", section)
		}
	}
}

func TestYAMLSyntaxRulesContent(t *testing.T) {
	// Verify key content is present
	requiredContent := []string{
		"version",
		"vars",
		"checks",
		"indentation",
		"quoted",
	}

	for _, content := range requiredContent {
		if !strings.Contains(strings.ToLower(YAMLSyntaxRules), content) {
			t.Errorf("YAMLSyntaxRules missing content about: %s", content)
		}
	}
}

func TestCheckStructureRulesContent(t *testing.T) {
	// Verify required fields are documented
	requiredFields := []string{
		"id",
		"run",
	}

	for _, field := range requiredFields {
		if !strings.Contains(CheckStructureRules, "**"+field+"**") {
			t.Errorf("CheckStructureRules missing required field documentation: %s", field)
		}
	}

	// Verify optional fields are documented
	optionalFields := []string{
		"grok",
		"assert",
		"severity",
		"suggestion",
		"requires",
		"timeout",
		"file",
	}

	for _, field := range optionalFields {
		if !strings.Contains(CheckStructureRules, "**"+field+"**") {
			t.Errorf("CheckStructureRules missing optional field documentation: %s", field)
		}
	}
}

func TestDependencyValidationRulesContent(t *testing.T) {
	// Verify key concepts are documented
	concepts := []string{
		"circular",
		"requires",
		"self-reference",
	}

	for _, concept := range concepts {
		if !strings.Contains(strings.ToLower(DependencyValidationRules), concept) {
			t.Errorf("DependencyValidationRules missing concept: %s", concept)
		}
	}

	// Should include examples
	if !strings.Contains(DependencyValidationRules, "```yaml") {
		t.Error("DependencyValidationRules should include YAML examples")
	}
}

func TestVariableInterpolationRulesContent(t *testing.T) {
	// Verify interpolation syntax is documented
	if !strings.Contains(VariableInterpolationRules, "{{.") {
		t.Error("VariableInterpolationRules should document {{.varname}} syntax")
	}

	// Verify fields where variables can be used
	fields := []string{"run", "assert", "suggestion", "file", "grok"}
	for _, field := range fields {
		if !strings.Contains(VariableInterpolationRules, "**"+field+"**") {
			t.Errorf("VariableInterpolationRules should mention %s field", field)
		}
	}
}

func TestExplicitDoNotListContent(t *testing.T) {
	// Verify key prohibitions
	prohibitions := []string{
		"DO NOT",
		"comment",
		"duplicate",
		"circular",
		"undefined",
	}

	for _, prohibition := range prohibitions {
		if !strings.Contains(ExplicitDoNotList, prohibition) {
			t.Errorf("ExplicitDoNotList missing prohibition about: %s", prohibition)
		}
	}
}

func TestAssertionOperators(t *testing.T) {
	expectedOperators := []string{"==", "!=", "<", "<=", ">", ">=", "&&", "||", "!"}

	for _, op := range expectedOperators {
		found := false
		for _, actual := range AssertionOperators {
			if actual == op {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AssertionOperators missing: %s", op)
		}
	}
}

func TestSpecialAssertionVariables(t *testing.T) {
	expectedVars := []string{"exit_code", "stdout", "stderr"}

	for _, v := range expectedVars {
		found := false
		for _, actual := range SpecialAssertionVariables {
			if actual == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SpecialAssertionVariables missing: %s", v)
		}
	}
}

func TestSupportedSeverities(t *testing.T) {
	if len(SupportedSeverities) != 2 {
		t.Errorf("Expected 2 severities, got %d", len(SupportedSeverities))
	}

	expectedSeverities := []string{"error", "warning"}
	for _, s := range expectedSeverities {
		found := false
		for _, actual := range SupportedSeverities {
			if actual == s {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SupportedSeverities missing: %s", s)
		}
	}
}

func TestGrokPatternExamples(t *testing.T) {
	// Verify common patterns are included
	expectedPatterns := []string{"NUMBER", "WORD", "INT", "GREEDYDATA"}

	for _, pattern := range expectedPatterns {
		if _, ok := GrokPatternExamples[pattern]; !ok {
			t.Errorf("GrokPatternExamples missing pattern: %s", pattern)
		}
	}

	// Verify patterns follow correct format
	for name, pattern := range GrokPatternExamples {
		if !strings.Contains(pattern, "%{"+name) {
			t.Errorf("GrokPatternExamples[%s] should contain %%{%s", name, name)
		}
		if !strings.Contains(pattern, ":varname}") {
			t.Errorf("GrokPatternExamples[%s] should contain :varname}", name)
		}
	}
}

func TestCommonTimeoutValues(t *testing.T) {
	// Verify common check types have timeout recommendations
	expectedTypes := []string{"format", "lint", "test", "build", "coverage"}

	for _, checkType := range expectedTypes {
		if _, ok := CommonTimeoutValues[checkType]; !ok {
			t.Errorf("CommonTimeoutValues missing type: %s", checkType)
		}
	}

	// Verify timeout format is valid Go duration
	for checkType, timeout := range CommonTimeoutValues {
		if !strings.HasSuffix(timeout, "s") && !strings.HasSuffix(timeout, "m") {
			t.Errorf("CommonTimeoutValues[%s] = %s should end with 's' or 'm'", checkType, timeout)
		}
	}
}

func TestGuideContainsExamples(t *testing.T) {
	guide := NewValidationGuide()
	fullGuide := guide.GetFullGuide()

	// Count code blocks
	codeBlocks := strings.Count(fullGuide, "```yaml")

	if codeBlocks < 5 {
		t.Errorf("Expected at least 5 YAML code examples, found %d", codeBlocks)
	}
}

func TestGuideReadability(t *testing.T) {
	guide := NewValidationGuide()
	fullGuide := guide.GetFullGuide()

	// Verify sections have headers
	headers := strings.Count(fullGuide, "## ")
	if headers < 5 {
		t.Errorf("Expected at least 5 section headers (##), found %d", headers)
	}

	// Verify subsections exist
	subheaders := strings.Count(fullGuide, "### ")
	if subheaders < 10 {
		t.Errorf("Expected at least 10 subsection headers (###), found %d", subheaders)
	}
}
