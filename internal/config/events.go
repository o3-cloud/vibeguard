package config

import "gopkg.in/yaml.v3"

// EventHandler defines the actions to take when a check completes with specific outcomes.
// Each event (success, failure, timeout) can contain either prompt ID references or inline content.
type EventHandler struct {
	// Success contains prompts to display when check passes (exit code 0, assertions true).
	Success EventValue `yaml:"success,omitempty"`
	// Failure contains prompts to display when check fails (exit code non-zero, assertions false).
	Failure EventValue `yaml:"failure,omitempty"`
	// Timeout contains prompts to display when check exceeds timeout.
	// Takes precedence over failure event.
	Timeout EventValue `yaml:"timeout,omitempty"`
}

// EventValue represents either prompt ID references or inline content.
// - Array of strings: treated as prompt ID references
// - Single string: treated as inline content (not a prompt ID)
type EventValue struct {
	// IDs stores prompt ID references when parsed from array syntax
	IDs []string
	// Content stores inline content when parsed from string syntax
	Content string
	// IsInline indicates whether this value is inline content (true) or prompt IDs (false)
	IsInline bool
}

// UnmarshalYAML implements custom YAML unmarshaling for EventValue.
// Handles both array syntax (prompt ID references) and string syntax (inline content).
func (ev *EventValue) UnmarshalYAML(value *yaml.Node) error {
	// Try array of strings first (prompt IDs)
	var ids []string
	if err := value.Decode(&ids); err == nil {
		*ev = EventValue{IDs: ids, IsInline: false}
		return nil
	}

	// Try single string (inline content)
	var content string
	if err := value.Decode(&content); err != nil {
		return err
	}
	*ev = EventValue{Content: content, IsInline: true}
	return nil
}

// MarshalYAML implements custom YAML marshaling for EventValue.
func (ev EventValue) MarshalYAML() (interface{}, error) {
	if ev.IsInline {
		return ev.Content, nil
	}
	return ev.IDs, nil
}

// UnmarshalYAML implements custom YAML unmarshaling for EventHandler.
// This handles the conversion from YAML node to EventHandler with proper EventValue handling.
func (eh *EventHandler) UnmarshalYAML(value *yaml.Node) error {
	// Create a temporary struct to unmarshal into
	type eventHandlerAlias struct {
		Success EventValue `yaml:"success,omitempty"`
		Failure EventValue `yaml:"failure,omitempty"`
		Timeout EventValue `yaml:"timeout,omitempty"`
	}

	var alias eventHandlerAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}

	eh.Success = alias.Success
	eh.Failure = alias.Failure
	eh.Timeout = alias.Timeout
	return nil
}
