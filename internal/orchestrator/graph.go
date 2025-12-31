package orchestrator

import (
	"fmt"

	"github.com/vibeguard/vibeguard/internal/config"
)

// DependencyGraph represents the check execution order.
type DependencyGraph struct {
	levels [][]string // Each level contains check IDs that can run in parallel
}

// BuildGraph creates a dependency graph from checks using topological sort.
// Returns execution levels where each level can be run in parallel.
func BuildGraph(checks []config.Check) (*DependencyGraph, error) {
	// Build lookup maps
	checkByID := make(map[string]*config.Check)
	for i := range checks {
		checkByID[checks[i].ID] = &checks[i]
	}

	// Validate all dependencies exist
	for _, check := range checks {
		for _, dep := range check.Requires {
			if _, ok := checkByID[dep]; !ok {
				return nil, fmt.Errorf("check %q requires unknown check: %s", check.ID, dep)
			}
		}
	}

	// Topological sort using Kahn's algorithm
	// Initialize in-degree for each node
	inDegree := make(map[string]int)
	for _, check := range checks {
		if _, ok := inDegree[check.ID]; !ok {
			inDegree[check.ID] = 0
		}
	}

	// Build adjacency list (dependency -> dependents)
	dependents := make(map[string][]string)
	for _, check := range checks {
		for _, dep := range check.Requires {
			dependents[dep] = append(dependents[dep], check.ID)
			inDegree[check.ID]++
		}
	}

	// Process levels
	var levels [][]string
	processed := make(map[string]bool)

	for len(processed) < len(checks) {
		// Find all checks with no unprocessed dependencies
		var level []string
		for _, check := range checks {
			if processed[check.ID] {
				continue
			}
			if inDegree[check.ID] == 0 {
				level = append(level, check.ID)
			}
		}

		if len(level) == 0 {
			// Circular dependency detected
			var unprocessed []string
			for _, check := range checks {
				if !processed[check.ID] {
					unprocessed = append(unprocessed, check.ID)
				}
			}
			return nil, fmt.Errorf("circular dependency detected among checks: %v", unprocessed)
		}

		levels = append(levels, level)

		// Mark as processed and update in-degrees
		for _, id := range level {
			processed[id] = true
			for _, dependent := range dependents[id] {
				inDegree[dependent]--
			}
		}
	}

	return &DependencyGraph{levels: levels}, nil
}

// Levels returns the execution levels.
func (g *DependencyGraph) Levels() [][]string {
	return g.levels
}
