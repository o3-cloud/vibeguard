# Security Policy

VibeGuard is a policy enforcement system designed for integration with CI/CD pipelines, agent loops, and development workflows. This document describes our security model, threat boundaries, and responsible disclosure process.

## Security Model

### Trust Boundaries

VibeGuard has clearly defined trust boundaries that determine what is considered a security vulnerability:

#### 1. Configuration Author ≡ Trust Boundary

The person writing `vibeguard.yaml` is the **primary trust boundary**. They have full control over:
- Variable definitions (`vars:` section)
- Command execution (`run:` field)
- Policy definitions and assertions

**Security Implication:** If the configuration author is malicious, they can execute arbitrary commands directly. Therefore, command injection through variables is not a vulnerability—the config author can already achieve the same outcome through direct command specification.

#### 2. Configuration File as Input

The `vibeguard.yaml` file must be treated as a trusted configuration source in your repository. Access control should be enforced at the repository level:
- Require code reviews for changes to `vibeguard.yaml`
- Use branch protection rules to prevent unauthorized modifications
- Audit who has commit access to the repository

#### 3. External Data: Display-Only

Grok-extracted values from command output are **never used for command execution**. These values are only used for:
- Rendering suggestions and fix messages
- Displaying check results to users
- JSON output formatting

This means command output cannot influence subsequent command execution.

### Threat Model

#### Vulnerabilities Within Scope

1. **Command Injection from Untrusted Sources**
   - Not applicable: VibeGuard does not accept untrusted input from environment variables, command-line arguments, or external sources at runtime
   - All inputs are statically defined in the YAML configuration

2. **Denial of Service (Resource Exhaustion)**
   - Long-running commands can consume resources
   - Infinite loops in assertions are possible
   - **Mitigation:** Use timeout configuration (`timeout:` field) to limit execution time

3. **Information Disclosure**
   - Command output may contain sensitive information
   - **Mitigation:** Be cautious about which environment variables are exposed in check output; use filtering in commands

4. **Path Traversal**
   - The `file:` field in checks reads from the filesystem
   - **Mitigation:** Only use relative paths or paths within your project; restrict file reads to intended directories

#### Vulnerabilities Out of Scope

1. **Configuration Injection**
   - Variables are defined statically in YAML, not interpolated from untrusted sources
   - Escaping shell metacharacters would break legitimate use cases (e.g., path patterns with dots and slashes)

2. **Command Injection via Output**
   - Grok-extracted values are only used for display and suggestions
   - They are never passed back into shell command execution

3. **Privilege Escalation**
   - VibeGuard runs with the privileges of the user executing it
   - No escalation mechanism exists

### Variable Interpolation Security

VibeGuard supports variable interpolation with the syntax `{{.variable}}` for injecting configuration variables into check commands. This is **secure by design**:

```yaml
vars:
  packages: "./cmd/... ./internal/..."
checks:
  - id: test
    run: go test {{.packages}}
```

**Why it's secure:**

1. **Static Definition:** Variables are defined once in the `vars:` section and interpolated as strings
2. **No External Input:** Runtime input from environment variables or command-line arguments is not supported
3. **Config Author Control:** The person writing the config controls both variables and the commands that use them
4. **Legitimate Use Cases:** Shell metacharacters in paths (dots, slashes) are intentionally preserved

### Exit Codes and CI/CD Integration

VibeGuard uses exit codes to communicate policy enforcement results:

| Exit Code | Meaning | CI/CD Behavior |
|-----------|---------|----------------|
| 0 | All checks passed | Pipeline continues |
| 1 | One or more checks failed | Pipeline fails (can retry) |
| 2 | Configuration error | Configuration must be fixed before retry |
| 3 | Execution error | Transient error, safe to retry |

This exit code model allows CI/CD systems to distinguish between:
- **Retryable errors** (exit 3) — transient failures, safe to retry
- **Configuration errors** (exit 2) — must fix config and resubmit
- **Policy violations** (exit 1) — developer action required

## Responsible Disclosure

If you discover a security vulnerability in VibeGuard, please report it responsibly:

### Reporting Process

1. **Do not open a public GitHub issue** for security vulnerabilities
2. **Email your findings** to the maintainers with:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Your suggested fix (if available)

3. **Timeline:**
   - We will acknowledge receipt within 48 hours
   - We aim to release a fix within 7 days for critical vulnerabilities
   - We will credit you in the release notes (unless you prefer anonymity)

### What We Consider a Vulnerability

- **In Scope:**
  - Remote code execution vulnerabilities
  - Privilege escalation
  - Authentication/authorization bypasses
  - Data exposure
  - Denial of service in the core VibeGuard process

- **Out of Scope:**
  - Issues in third-party dependencies (report to upstream)
  - Issues arising from insecure configuration by users
  - Performance issues or feature requests
  - Documentation clarity

## Security Best Practices

### For VibeGuard Users

1. **Protect Your Configuration**
   - Treat `vibeguard.yaml` as sensitive configuration
   - Review changes in pull requests
   - Use branch protection rules

2. **Limit Check Scope**
   - Only run checks on code/artifacts you trust
   - Be cautious about checks that read arbitrary files
   - Use timeouts to prevent resource exhaustion

3. **Monitor Output**
   - Check that commands in your checks don't expose sensitive information
   - Review environment variables passed to checks
   - Use `--quiet` mode to suppress detailed output in public CI logs

4. **Update Regularly**
   - Keep VibeGuard up to date
   - Monitor releases for security updates
   - Pin to specific versions in CI/CD systems

### For VibeGuard Developers

1. **Assume Hostile Configuration**
   - Treat user-provided YAML as potentially malicious
   - Validate all inputs thoroughly
   - Use safe defaults

2. **Test Edge Cases**
   - Test with extremely large inputs (memory safety)
   - Test with deeply nested structures (stack overflow prevention)
   - Test with special characters and escape sequences

3. **Code Review**
   - Security-sensitive code changes require review
   - Pay special attention to parser, executor, and interpolation code
   - Check for command injection risks

4. **Automated Testing**
   - Maintain comprehensive test coverage (70%+ minimum)
   - Use mutation testing to verify assertion quality
   - Run static analysis tools (golangci-lint)

## Implementation Details

### Shell Execution Model

VibeGuard executes checks using `/bin/sh -c <command>` with inherited environment variables. This model:

- Provides compatibility across Unix-like systems
- Allows complex command pipelines
- Requires careful handling of untrusted input (which we don't accept)

### File Reading

The `file:` field in checks reads files from the filesystem. Security considerations:

- Only relative paths or absolute paths starting with project root should be used
- File permissions are enforced by the OS
- No file encryption or access control is implemented by VibeGuard

### Dependencies

VibeGuard has minimal dependencies:

- **Standard Library Only** for core functionality
- **go-grok** for pattern matching (reviewed for security)
- **go.yaml** for configuration parsing

All dependencies are vendored and reviewed.

## Security Audit and Testing

VibeGuard undergoes regular security reviews:

- **Code Review:** All changes reviewed for security implications
- **Static Analysis:** golangci-lint with security-focused linters
- **Dependency Scanning:** Regular checks for known vulnerabilities
- **Mutation Testing:** Verification that test suite catches real defects

See [docs/log/](docs/log/) for security review findings.

## Related Documentation

- [Architecture Decision Records](docs/adr/) — Design decisions including security considerations
- [CONTRIBUTING.md](CONTRIBUTING.md) — Guidelines for secure development
- [CONVENTIONS.md](CONVENTIONS.md) — Code quality and style standards
- [Internal Security Reviews](docs/log/) — Security analysis and findings

## Questions?

For security questions that don't constitute a vulnerability report:
- Review this document and related ADRs
- Check existing issues and discussions
- Open a GitHub discussion in the security category
