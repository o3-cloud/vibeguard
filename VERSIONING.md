# Versioning Policy

VibeGuard follows **Semantic Versioning (SemVer)** with a clear stability commitment for enterprise adoption. This document outlines the versioning scheme, breaking change policy, deprecation guidelines, and support lifecycle.

## Semantic Versioning

VibeGuard releases follow the [Semantic Versioning 2.0.0](https://semver.org/) specification: `MAJOR.MINOR.PATCH`.

### Version Components

- **MAJOR** (X.0.0): Incremented for breaking changes to the public API or CLI interface
- **MINOR** (0.X.0): Incremented for backward-compatible feature additions
- **PATCH** (0.0.X): Incremented for bug fixes and non-breaking improvements

### Pre-release and Build Metadata

Pre-release versions use the format `vX.Y.Z-alpha`, `vX.Y.Z-beta`, `vX.Y.Z-rc1`, etc.

Examples:
- `v0.1.0-alpha` - Alpha release
- `v0.1.0-beta.1` - Beta release iteration
- `v0.1.0-rc.1` - Release candidate
- `v0.1.0` - Stable release

## Breaking Changes

A **breaking change** is any modification that requires users to update their code or configuration to maintain compatibility.

### Examples of Breaking Changes

**CLI Interface:**
- Renaming or removing a command (`check`, `init`, `list`, `validate`)
- Changing the meaning of a flag or removing a flag entirely
- Changing required command-line arguments
- Altering default behavior in a non-backward-compatible way

**Configuration (vibeguard.yaml):**
- Removing or renaming top-level configuration fields
- Changing the type of a field (string → boolean)
- Removing required check properties
- Altering assertion syntax or expression operators in incompatible ways

**Check System:**
- Changing how check dependencies work in a breaking way
- Removing or fundamentally changing built-in check types
- Altering exit code behavior from existing check types

**Exit Codes:**
- Changing the meaning of existing exit codes (2, 3, etc.)
- Adding mandatory exit code requirements in a breaking way

**JSON Output Schema:**
- Removing fields from JSON output structure
- Changing the type of existing JSON fields
- Altering the structure of nested JSON objects

### Breaking Change Documentation

All breaking changes must be:

1. **Documented in CHANGELOG.md** under the relevant version with a `BREAKING CHANGE:` footer in commit messages
2. **Documented in Release Notes** with migration guide
3. **Announced in Major Version Release** with clear migration instructions
4. **Documented in this file** for cumulative reference

### Breaking Change Example

```
feat(cli): rename --output to --format

BREAKING CHANGE: The --output flag has been renamed to --format for consistency
with industry standards. Update any scripts using --output to use --format instead.

Migration: Change `vibeguard check --output json` to `vibeguard check --format json`
```

## Deprecation Policy

A **deprecation** is a planned removal of a feature that allows users time to migrate.

### Deprecation Lifecycle

1. **Announcement Phase** (Deprecation Release)
   - Feature is marked as deprecated in documentation
   - Warning messages guide users to the replacement
   - Feature remains fully functional

2. **Support Phase** (Minimum 2 Minor Versions)
   - Feature works without changes for at least 2 minor version releases
   - Clear migration path provided in documentation
   - Example: Deprecated in v0.5.0, can be removed in v0.7.0 or later

3. **Removal Phase** (Major Version Release)
   - Feature is removed in the next major version
   - Removal constitutes a breaking change
   - Removed features are documented with migration instructions

### Deprecation Messaging

Deprecated features should provide helpful guidance:

```
Warning: The --verbose flag is deprecated and will be removed in v1.0.0.
Please use --output verbose instead.
```

### Deprecation Example Timeline

- **v0.5.0**: `--verbose` flag marked deprecated (recommend `--output verbose`)
- **v0.5.1 - v0.6.x**: Feature supported, deprecation warnings issued
- **v1.0.0**: Feature removed, breaking change documented

## Stability Levels

### Current Stability: Pre-release (v0.x.x)

The project is in **active development** with rapid feature addition and potential breaking changes.

**Expectations:**
- API and CLI are subject to change without notice
- Breaking changes may occur in minor version releases
- Configuration format may evolve
- Not recommended for production critical systems without acceptance of instability

### Planned Stability: Stable (v1.0.0+)

Upon reaching v1.0.0:

**Expectations:**
- Public API/CLI stability maintained across minor versions
- Breaking changes only in major releases
- Configuration format stability
- Security patches backported to recent minor versions
- Production-ready stability guarantees

## Support Lifecycle

### Release Support

VibeGuard follows an **N-1 support model** once stable (v1.0.0+):

- **Current Release** (vX.Y.Z): Receives features, fixes, and security patches
- **Previous Minor Release** (vX.(Y-1).Z): Receives critical fixes and security patches only
- **Earlier Releases**: No longer supported; users must upgrade

Example at v1.3.0:
- v1.3.x: All updates
- v1.2.x: Critical fixes and security patches only
- v1.1.x or earlier: No support

### Security Patches

Security vulnerabilities in supported versions are patched within 7 days of responsible disclosure.

- Patches are released as PATCH versions
- Security advisories are published in GitHub Security Advisories
- All supported versions may receive patches

### Long-Term Support

There are no Long-Term Support (LTS) versions planned at this time. Users are encouraged to stay current with releases.

## Version Numbering Examples

```
v0.1.0      First alpha release
v0.1.1      Bug fix
v0.2.0      New features, backward compatible
v0.2.1      Bug fix
v0.3.0      More features and improvements
v1.0.0      First stable release
v1.0.1      Critical bug fix
v1.1.0      New features, backward compatible
v2.0.0      Major breaking changes, backward incompatible with v1.x
```

## Configuration Versioning

VibeGuard configuration files do not have explicit version numbers. However, the project maintains backward compatibility for config files across minor versions.

**Config Format Changes:**
- Minor versions: New optional fields may be added (backward compatible)
- Major versions: Required config format changes may occur

Configuration migration guides are provided in release notes.

## Release Checklist

Before releasing a new version:

1. ✓ Update `internal/version/version.go` with new version
2. ✓ Update `CHANGELOG.md` with release notes
3. ✓ Verify all breaking changes are documented
4. ✓ Run full test suite (unit, integration, mutation)
5. ✓ Verify code coverage meets or exceeds 70%
6. ✓ Tag release with semantic version (`git tag vX.Y.Z`)
7. ✓ Build and test binary on all target platforms
8. ✓ Create GitHub release with CHANGELOG entry and migration guides
9. ✓ Announce on appropriate channels

## References

- [Semantic Versioning 2.0.0](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Conventional Commits](https://www.conventionalcommits.org/) - See ADR-002
- [CHANGELOG.md](CHANGELOG.md) - VibeGuard Release History
