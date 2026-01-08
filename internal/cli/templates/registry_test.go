package templates

import (
	"testing"
)

func TestListReturnsAllTemplates(t *testing.T) {
	templates := List()

	// Should have all 10 templates
	expectedNames := []string{
		"bun-typescript",
		"generic",
		"go-minimal",
		"go-standard",
		"node-javascript",
		"node-typescript",
		"python-pip",
		"python-poetry",
		"python-uv",
		"rust-cargo",
	}

	if len(templates) != len(expectedNames) {
		t.Errorf("expected %d templates, got %d", len(expectedNames), len(templates))
	}

	// Verify templates are sorted
	for i, tmpl := range templates {
		if tmpl.Name != expectedNames[i] {
			t.Errorf("expected template[%d] to be %q, got %q", i, expectedNames[i], tmpl.Name)
		}
	}
}

func TestNamesReturnsAllNames(t *testing.T) {
	names := Names()

	expectedNames := []string{
		"bun-typescript",
		"generic",
		"go-minimal",
		"go-standard",
		"node-javascript",
		"node-typescript",
		"python-pip",
		"python-poetry",
		"python-uv",
		"rust-cargo",
	}

	if len(names) != len(expectedNames) {
		t.Errorf("expected %d names, got %d", len(expectedNames), len(names))
	}

	for i, name := range names {
		if name != expectedNames[i] {
			t.Errorf("expected name[%d] to be %q, got %q", i, expectedNames[i], name)
		}
	}
}

func TestGetExistingTemplate(t *testing.T) {
	tmpl, err := Get("go-standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tmpl.Name != "go-standard" {
		t.Errorf("expected name %q, got %q", "go-standard", tmpl.Name)
	}

	if tmpl.Description == "" {
		t.Error("expected non-empty description")
	}

	if tmpl.Content == "" {
		t.Error("expected non-empty content")
	}

	// Verify content is valid YAML with version
	if !contains(tmpl.Content, "version:") {
		t.Error("expected content to contain version field")
	}

	if !contains(tmpl.Content, "checks:") {
		t.Error("expected content to contain checks field")
	}
}

func TestGetNonExistentTemplate(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestExistsForExistingTemplate(t *testing.T) {
	if !Exists("go-standard") {
		t.Error("expected go-standard to exist")
	}
}

func TestExistsForNonExistentTemplate(t *testing.T) {
	if Exists("nonexistent") {
		t.Error("expected nonexistent to not exist")
	}
}

func TestAllTemplatesHaveRequiredFields(t *testing.T) {
	templates := List()

	for _, tmpl := range templates {
		t.Run(tmpl.Name, func(t *testing.T) {
			if tmpl.Name == "" {
				t.Error("template name is empty")
			}
			if tmpl.Description == "" {
				t.Errorf("template %q has empty description", tmpl.Name)
			}
			if tmpl.Content == "" {
				t.Errorf("template %q has empty content", tmpl.Name)
			}
			if !contains(tmpl.Content, "version:") {
				t.Errorf("template %q content missing version field", tmpl.Name)
			}
			if !contains(tmpl.Content, "checks:") {
				t.Errorf("template %q content missing checks field", tmpl.Name)
			}
		})
	}
}

func TestGoStandardTemplateContent(t *testing.T) {
	tmpl, err := Get("go-standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have standard Go checks
	expectedChecks := []string{"fmt", "vet", "lint", "test", "coverage", "build"}
	for _, check := range expectedChecks {
		if !contains(tmpl.Content, "id: "+check) {
			t.Errorf("go-standard template missing check %q", check)
		}
	}
}

func TestNodeTypescriptTemplateContent(t *testing.T) {
	tmpl, err := Get("node-typescript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have TypeScript-specific checks
	expectedChecks := []string{"format", "lint", "typecheck", "test", "build"}
	for _, check := range expectedChecks {
		if !contains(tmpl.Content, "id: "+check) {
			t.Errorf("node-typescript template missing check %q", check)
		}
	}
}

func TestBunTypescriptTemplateContent(t *testing.T) {
	tmpl, err := Get("bun-typescript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have TypeScript-specific checks for Bun
	expectedChecks := []string{"format", "lint", "typecheck", "test", "build"}
	for _, check := range expectedChecks {
		if !contains(tmpl.Content, "id: "+check) {
			t.Errorf("bun-typescript template missing check %q", check)
		}
	}
}

func TestRustCargoTemplateContent(t *testing.T) {
	tmpl, err := Get("rust-cargo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have Rust/Cargo-specific checks
	expectedChecks := []string{"fmt", "clippy", "test", "build"}
	for _, check := range expectedChecks {
		if !contains(tmpl.Content, "id: "+check) {
			t.Errorf("rust-cargo template missing check %q", check)
		}
	}
}

func TestPythonUvTemplateContent(t *testing.T) {
	tmpl, err := Get("python-uv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have Python uv-specific checks
	expectedChecks := []string{"format", "lint", "typecheck", "test", "coverage"}
	for _, check := range expectedChecks {
		if !contains(tmpl.Content, "id: "+check) {
			t.Errorf("python-uv template missing check %q", check)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
