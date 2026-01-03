---
summary: Investigation of bug vibeguard-7vh regarding Go version extraction from go.mod - found working as designed
event_type: research
sources:
  - internal/cli/inspector/metadata.go
  - internal/cli/inspector/metadata_test.go
  - go.mod
tags:
  - bug-investigation
  - metadata-extractor
  - go-version
  - inspector
  - ai-assisted-setup
---

# Investigation: Go Version Extraction Bug (vibeguard-7vh)

## Bug Description

The reported bug claimed: "The MetadataExtractor for Go projects returns an empty Version field. It should extract the Go version from go.mod (e.g., 'go 1.21')."

## Investigation Findings

### Current Implementation

The `extractGoMetadata()` function in `internal/cli/inspector/metadata.go:118-155` handles Go project metadata extraction:

1. **Module name extraction**: Uses regex `^module\s+(.+)$` to extract the module path and stores it in `metadata.Name`

2. **Go version extraction**: Uses regex `^go\s+(\d+\.\d+(?:\.\d+)?)$` to extract the Go version and stores it in `metadata.Extra["go_version"]`

3. **Project version extraction**: Attempts to read from a `VERSION` file if present and stores in `metadata.Version`

### Test Verification

Ran the metadata extractor against vibeguard's own `go.mod`:

```
Name: "github.com/vibeguard/vibeguard"
Version: ""
Extra: map[go_version:1.24.4]
```

The Go version (`1.24.4`) IS being correctly extracted and stored in `Extra["go_version"]`.

### Semantic Design

The implementation follows correct semantic design:
- `Version` field = Project version (e.g., `1.0.0`, `2.3.1`)
- `Extra["go_version"]` = Language/toolchain version (e.g., `1.21`, `1.24.4`)

This is consistent with other language implementations:
- Node.js: `Extra["node_version"]` stores engine requirements
- Python: `Extra["python_version"]` stores `requires-python`
- Rust: `Extra["rust_edition"]` stores the Rust edition

### Conclusion

The bug was a misunderstanding. The Go version from `go.mod` IS being extracted correctly to `metadata.Extra["go_version"]`. The `Version` field is empty because Go projects don't specify a project version in `go.mod` - that information typically comes from git tags or a `VERSION` file.

## Resolution

Closed vibeguard-7vh as "working as designed" with explanation of the semantic difference between project version and language version fields.

## Lessons Learned

- The `Extra` map in `ProjectMetadata` serves as the correct location for language-specific metadata that doesn't fit the standard fields
- Documentation could be improved to clarify the distinction between `Version` (project) and `Extra["go_version"]` (language)
