package config

import (
	"bytes"
	"strings"
	"text/template"
)

// Interpolate replaces {{.VAR}} placeholders in the config with variable values.
func (c *Config) Interpolate() {
	for i := range c.Checks {
		c.Checks[i].Run = c.interpolateString(c.Checks[i].Run)
		c.Checks[i].Assert = c.interpolateString(c.Checks[i].Assert)
		c.Checks[i].Suggestion = c.interpolateString(c.Checks[i].Suggestion)
		c.Checks[i].Fix = c.interpolateString(c.Checks[i].Fix)
		c.Checks[i].File = c.interpolateString(c.Checks[i].File)

		for j := range c.Checks[i].Grok {
			c.Checks[i].Grok[j] = c.interpolateString(c.Checks[i].Grok[j])
		}
	}
}

// interpolateString replaces {{.VAR}} with variable values.
func (c *Config) interpolateString(s string) string {
	if s == "" {
		return s
	}

	result := s
	for key, value := range c.Vars {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// InterpolateWithExtracted renders a Go template string with config vars and
// extracted values from grok patterns available as {{.varname}}.
// Config vars take precedence over extracted values if there's a conflict.
func InterpolateWithExtracted(templateStr string, vars map[string]string, extracted map[string]string) string {
	if templateStr == "" {
		return templateStr
	}

	// Merge extracted values first, then config vars (so vars take precedence)
	data := make(map[string]string)
	for key, value := range extracted {
		data[key] = value
	}
	for key, value := range vars {
		data[key] = value
	}

	// Parse and execute the template
	tmpl, err := template.New("suggestion").Parse(templateStr)
	if err != nil {
		// If template parsing fails, return original string
		return templateStr
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// If template execution fails, return original string
		return templateStr
	}

	return buf.String()
}
