---
summary: Implemented metadata extraction for AI-assisted setup inspector (vibeguard-9mi.3)
event_type: code
sources:
  - internal/cli/inspector/metadata.go
  - internal/cli/inspector/metadata_test.go
  - internal/cli/inspector/detector.go
  - internal/cli/inspector/tools.go
tags:
  - ai-assisted-setup
  - inspector
  - metadata-extraction
  - go
  - node
  - python
  - rust
  - ruby
  - java
---

# Metadata Extraction Implementation

Completed task **vibeguard-9mi.3**: Phase 1: Repository Inspector - Metadata Extraction

## Implementation Summary

Added comprehensive metadata extraction capabilities to the inspector package, enabling the AI-assisted setup feature to understand project configuration files across multiple languages.

### New Types

**ProjectMetadata** - Holds extracted metadata from project configuration files:
- Name, Version, Description, License, Repository, Author
- Keywords array
- Extra map for language-specific fields (e.g., go_version, node_version, rust_edition)

**ProjectStructure** - Describes the layout of a project:
- EntryPoints (main files like cmd/main.go, src/index.js)
- SourceDirs (src, lib, pkg, internal)
- TestDirs (test, tests, __tests__, spec)
- ConfigFiles found
- HasMonorepo detection
- BuildOutputDir

**MetadataExtractor** - Main extraction logic with methods:
- `Extract(projectType)` - Extracts metadata based on detected project type
- `ExtractStructure(projectType)` - Analyzes project directory structure

### Supported Languages

| Language | Config Files | Metadata Extracted |
|----------|-------------|-------------------|
| Go | go.mod, VERSION | module name, go version |
| Node | package.json | name, version, description, author, license, repository, engines, keywords |
| Python | pyproject.toml, setup.py | name, version, description, author, license, python version |
| Rust | Cargo.toml | name, version, description, license, repository, edition |
| Ruby | *.gemspec | name, version, description, author, license |
| Java | pom.xml, build.gradle | name/artifactId, version, description, groupId |

### Monorepo Detection

The implementation detects monorepo patterns including:
- npm/yarn workspaces in package.json
- pnpm-workspace.yaml
- lerna.json
- Cargo workspace in Cargo.toml
- packages/, apps/, libs/ directories with multiple sub-packages

## Test Coverage

Added 21 test cases covering:
- Metadata extraction for all 6 supported languages
- Various config file formats (pyproject.toml vs setup.py, pom.xml vs build.gradle)
- Project structure detection for Go, Node, Python, Rust
- Monorepo pattern detection (5 different patterns)
- Edge cases: missing files, empty projects, unknown project types

All tests pass with the full test suite remaining green.

## Design Decisions

1. **Simple TOML/XML parsing via regex** - Avoided adding dependencies for full parsers since we only need specific fields. This keeps the binary small.

2. **Graceful degradation** - Missing files return empty metadata without errors, allowing partial extraction.

3. **Language-specific extra fields** - Used a map for language-specific metadata (go_version, rust_edition, etc.) rather than bloating the main struct.

4. **Monorepo detection as structure property** - Included in ProjectStructure since it affects how AI agents should recommend configurations.

## Next Steps

This completes Phase 1 of the AI-assisted setup feature. The next tasks are:
- **vibeguard-9mi.4**: Unit tests for repository inspector (already partially done)
- **vibeguard-9mi.5**: Recommendation engine to suggest checks based on detected metadata

## Related

- ADR-003: Go as primary implementation language
- Task vibeguard-9mi: AI Agent-Assisted Setup Feature epic
