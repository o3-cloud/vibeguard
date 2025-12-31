package config

import (
	"strings"
)

// Interpolate replaces {{.VAR}} placeholders in the config with variable values.
func (c *Config) Interpolate() {
	for i := range c.Checks {
		c.Checks[i].Run = c.interpolateString(c.Checks[i].Run)
		c.Checks[i].Assert = c.interpolateString(c.Checks[i].Assert)
		c.Checks[i].Suggestion = c.interpolateString(c.Checks[i].Suggestion)
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

// InterpolateWithExtracted replaces {{.VAR}} in a suggestion string with
// both config vars and extracted values from grok patterns.
func InterpolateWithExtracted(template string, vars map[string]string, extracted map[string]string) string {
	if template == "" {
		return template
	}

	result := template

	// Apply config vars first
	for key, value := range vars {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Apply extracted values
	for key, value := range extracted {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
