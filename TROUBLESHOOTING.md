# VibeGuard Troubleshooting Guide

This guide helps you resolve common issues when using VibeGuard.

## Common Issues

### Configuration Errors

#### Error: "Configuration file not found"

**Causes:**
- No configuration file exists in the current directory or parent directories
- Configuration file has an incorrect name

**Solutions:**
1. Check that you have a `vibeguard.yaml` or `vibeguard.yml` file in your project root:
   ```bash
   ls -la vibeguard.yaml
   ```

2. Create a default configuration if none exists:
   ```bash
   vibeguard init
   ```

3. Specify an explicit config path:
   ```bash
   vibeguard check --config /path/to/vibeguard.yaml
   ```

---

#### Error: "YAML validation failure" or "invalid configuration"

**Causes:**
- Malformed YAML syntax (incorrect indentation, missing colons)
- Invalid field types (string instead of boolean)
- Required fields are missing

**Solutions:**
1. Validate your YAML syntax:
   ```bash
   # Use a YAML validator or IDE with YAML support
   # Check for proper indentation (2 spaces, no tabs)
   ```

2. Review the configuration schema in [README.md](README.md#configuration-schema)

3. Use `vibeguard init --assist` to generate a valid configuration for your project:
   ```bash
   vibeguard init --assist
   ```

4. Check that all required top-level fields are present:
   - `version` (required)
   - `checks` (required, must be a list)

---

### Cyclic Dependency Errors

#### Error: "cyclic dependency detected"

**What this means:**
Checks have a circular dependency in their `requires` fields, preventing execution order determination.

**Example of cyclic dependency:**
```yaml
checks:
  - id: check-a
    requires: [check-b]
  - id: check-b
    requires: [check-c]
  - id: check-c
    requires: [check-a]  # Creates cycle: A → B → C → A
```

**Solutions:**
1. Visualize your dependency graph to identify cycles
2. Break the cycle by removing one dependency
3. Common fixes:
   - Move checks to run in parallel if they don't truly depend on each other
   - Combine dependent checks into a single check if appropriate
   - Remove transitive dependencies that are not strictly necessary

---

### Grok Pattern Errors

#### Error: "failed to compile grok pattern" or "grok pattern failed to parse"

**Causes:**
- Invalid grok pattern syntax
- Undefined grok pattern reference (e.g., `%{UNDEFINED_PATTERN}`)
- Pattern expects structured input but receives unstructured output

**Solutions:**

1. **Validate pattern syntax:**
   - Use valid grok pattern syntax: `%{PATTERN_NAME:field_name}`
   - Common patterns: `%{IP}`, `%{TIMESTAMP}`, `%{WORD}`, `%{NUMBER}`
   - Escape special regex characters in custom patterns

2. **Check the output being parsed:**
   - Ensure the command output matches your pattern expectations
   - Test with a simpler pattern first:
     ```yaml
     grok:
       - "%{GREEDYDATA:message}"
     ```

3. **Use built-in patterns:**
   - VibeGuard uses Elastic grok patterns
   - Reference: https://github.com/elastic/go-grok
   - Example working pattern:
     ```yaml
     checks:
       - id: logs
         grok:
           - "%{TIMESTAMP_ISO8601:timestamp} %{LOGLEVEL:level} %{GREEDYDATA:message}"
     ```

4. **Debugging tips:**
   - Add verbose output to see full error details:
     ```bash
     vibeguard check -v
     ```
   - Test your pattern separately with sample data
   - If multiple patterns are used, ensure later patterns don't override critical fields unintentionally

---

### Assertion Errors

#### Error: "assertion failed" or invalid assertion expression

**Causes:**
- Invalid assertion syntax
- Referencing non-existent grok fields
- Type mismatch (comparing string to number)
- Logic errors in conditional expressions

**Solutions:**

1. **Verify assertion syntax:**
   - Use valid comparison operators: `==`, `!=`, `>`, `<`, `>=`, `<=`, `contains`, `matches`
   - Field references must exist in grok output
   - Example:
     ```yaml
     assert:
       - "exit_code == 0"
       - "message contains 'success'"
     ```

2. **Ensure grok fields are captured:**
   - Verify your grok pattern captures the field you're asserting on
   - Test with verbose mode:
     ```bash
     vibeguard check -v
     ```

3. **Check data types:**
   - Numeric fields: `exit_code == 0` (no quotes)
   - String fields: `message == "success"` (with quotes)
   - Contains/matches work with strings

4. **Common assertion patterns:**
   ```yaml
   assert:
     - "exit_code == 0"           # Command succeeded
     - "output contains 'pass'"   # Output includes text
     - "duration < 5000"          # Performance check (ms)
     - "status != 'error'"        # Status validation
   ```

---

### Timeout Errors

#### Error: Exit code 4 or "Check execution timeout"

**What this means:**
A check command exceeded the configured timeout limit and was terminated.

**Causes:**
- Check command takes longer than the timeout duration
- Command is hanging or stuck in an infinite loop
- System is under heavy load

**Solutions:**

1. **Increase the timeout:**
   ```yaml
   checks:
     - id: slow-check
       command: "long-running-script.sh"
       timeout: 300s  # Increase from default (usually 30s)
   ```

2. **Optimize the check command:**
   - Profile the slow command to find bottlenecks
   - Add early exit conditions
   - Reduce dataset size or parallelism if applicable

3. **Check default timeout:**
   - View your check's timeout in the config or logs
   - VibeGuard has a global default timeout; individual checks can override it
   - Ensure timeout is reasonable for your environment

4. **Verify the command works:**
   ```bash
   time your-command-here
   ```

---

### Check Not Found Errors

#### Error: "check not found" or "unknown check ID"

**Causes:**
- Specified check ID doesn't exist in configuration
- Typo in check ID name
- Check was commented out or removed from config

**Solutions:**

1. **List available checks:**
   ```bash
   vibeguard list
   ```

2. **Verify check ID spelling:**
   - Check IDs are case-sensitive
   - Must match exactly as defined in `vibeguard.yaml`

3. **Check configuration:**
   - Verify the check exists in your config file
   - Ensure it's not commented out
   - Look for proper YAML structure

4. **Example correct usage:**
   ```bash
   vibeguard check fmt   # Runs check with id: fmt
   ```

---

## Exit Codes and Meanings

| Code | Name | Meaning | Action |
|------|------|---------|--------|
| 0 | Success | All checks passed | Proceed normally |
| 2 | ConfigError | Configuration error | Fix configuration file before running |
| 3 | Violation | Check violations detected | Review failures and fix issues |
| 4 | Timeout | Check execution timeout or error | Increase timeout or debug command |

---

## Debugging Strategies

### Enable Verbose Output

```bash
vibeguard check -v
```

Verbose mode shows:
- All check results (not just failures)
- Detailed error messages
- Timing information for each check

### Use JSON Output for Parsing

```bash
vibeguard check --json
```

JSON output includes:
- Structured error details
- Full stdout/stderr from commands
- Timing and exit code information

### Run Single Check for Isolation

```bash
vibeguard check check-id
```

Running a single check helps isolate issues and speeds up testing.

### Check Command Independently

Run the underlying command directly to verify it works:

```bash
# If your check runs a shell command, test it:
sh -c "your-command-here"

# Check exit code:
echo $?
```

---

## Common Patterns and Solutions

### "Check not executing" or "No output"

1. Verify the command exists and is in PATH:
   ```bash
   which command-name
   ```

2. Check file permissions:
   ```bash
   ls -la script.sh
   chmod +x script.sh
   ```

3. Test the command manually in your shell

### "Grok pattern matches but assertion fails"

1. Print extracted fields:
   ```bash
   vibeguard check -v
   ```

2. Check field names in grok pattern match assertion references

3. Verify data types (quote strings, don't quote numbers)

### "Dependency order seems wrong"

1. Check the dependency graph:
   ```bash
   vibeguard list  # Shows all checks and dependencies
   ```

2. Ensure dependencies form a valid DAG (no cycles)

3. Consider reordering checks if dependencies can be removed

---

## Getting Help

If you encounter issues not covered here:

1. **Check the documentation:**
   - [README.md](README.md) - Configuration schema and CLI reference
   - Configuration examples in the repo

2. **Enable verbose output for more details:**
   ```bash
   vibeguard check -v
   ```

3. **Review error messages carefully** - they usually point to the specific issue

4. **Test components independently:**
   - Test grok patterns separately
   - Run check commands manually
   - Validate YAML configuration
