# ADR-003: Adopt Go as the Primary Implementation Language

## Context and Problem Statement

VibeGuard is a policy enforcement system designed to be a lightweight, composable tool that integrates with CI/CD pipelines, agent loops, and Cloud Code workflows. The choice of implementation language directly impacts:

- **Performance and resource usage** — VibeGuard should be fast and have minimal overhead
- **Deployment and distribution** — Single binary delivery is preferred over runtime dependencies
- **Developer experience** — The language should support clear policy definition and extensibility
- **Integration with existing tools** — Shell integration, cross-platform compatibility, and ability to invoke external commands are critical
- **Maintainability** — The codebase should be easy for teams to understand and extend

A decision on the primary implementation language is necessary to guide architecture and initial development.

## Considered Options

### Option A: Go

**Characteristics:**
- Compiled to a single, statically-linked binary
- Fast compilation and execution
- Built-in concurrency primitives (goroutines, channels)
- Excellent standard library with strong support for:
  - File I/O and shell integration
  - JSON parsing and structured data handling
  - YAML parsing via community packages
  - HTTP clients and servers
- Cross-platform support (Linux, macOS, Windows)
- Simple syntax with minimal boilerplate
- Strong ecosystem for CLI tools (Cobra, etc.)
- Aligns with cloud-native and DevOps tooling trends

**Trade-offs:**
- Less dynamic than scripting languages (requires compilation)
- Smaller community compared to Python or JavaScript
- Learning curve for developers unfamiliar with static typing

### Option B: Python

**Characteristics:**
- Rapid development cycle
- Excellent for scripting and policy definition
- Rich ecosystem for data processing, testing, and configuration (YAML, JSON)
- Strong support for LLM integration via libraries (OpenAI SDK, etc.)
- Familiar to many developers
- Good readiness for dynamic policy evaluation and judgment logic

**Trade-offs:**
- Runtime dependency required (Python interpreter)
- Slower execution compared to compiled languages
- Distribution complexity (packaging, virtual environments)
- Performance overhead for policy evaluation loops
- Larger memory footprint

### Option C: Rust

**Characteristics:**
- Extreme performance and safety guarantees
- Compiled to a single binary
- Strong type system and error handling
- Excellent for systems programming

**Trade-offs:**
- Steep learning curve and verbose syntax
- Slower development velocity
- Overkill for a policy enforcement tool
- Smaller DevOps community adoption

### Option D: TypeScript/Node.js

**Characteristics:**
- Familiar to JavaScript developers
- Good for rapid prototyping
- Strong ecosystem for LLM integration
- Good tooling and package management

**Trade-offs:**
- Runtime dependency (Node.js)
- Larger binary size and memory footprint
- Slower startup time (problematic for CLI tools invoked frequently)
- Less common in DevOps/infrastructure tooling

## Decision Outcome

**Chosen option:** Option A (Go)

**Rationale:**

Go is the best fit for VibeGuard because:

1. **Single Binary Distribution** — VibeGuard can be deployed as a single executable with no runtime dependencies. This is critical for seamless integration into CI/CD pipelines and agent loops where minimal overhead is expected.

2. **Performance** — Go's compiled nature ensures fast startup time and efficient policy evaluation. VibeGuard is designed to be invoked frequently, potentially multiple times per CI run. Go's performance is essential here.

3. **Policy Orchestration** — Go's excellent standard library and strong support for YAML/JSON parsing make it ideal for implementing the declarative policy runners and orchestration patterns documented in our implementation patterns. The simplicity of Go's syntax also aids readability for non-trivial policy evaluation logic.

4. **LLM Integration** — Go has mature libraries for calling external APIs and handling structured I/O. Judge invocations can be cleanly implemented without the overhead of a runtime environment.

5. **CLI Tooling Ecosystem** — Go dominates the CLI/DevOps tooling space. Frameworks like Cobra provide excellent support for command-line argument parsing, help generation, and structured output. This aligns with VibeGuard's vision as a Cloud Code companion command.

6. **Cross-Platform Compatibility** — Go's cross-compilation story is exceptional. Building for Linux, macOS, and Windows from a single codebase requires minimal configuration.

7. **Future Growth** — Go's concurrency model (goroutines) provides a natural path if VibeGuard needs to parallelize policy evaluation in the future without significant architectural changes.

**Tradeoffs Accepted:**

- We accept the need for a compilation step during development. This is a minor cost for the benefits of a single binary and performance.
- We accept that Go is less dynamic than Python, but this aligns with VibeGuard's philosophy of deterministic, declarative policy definition. The structured approach is a feature, not a limitation.

## Consequences

### Positive Outcomes

- **Frictionless Deployment:** VibeGuard can be distributed as a single binary, requiring no setup or dependency management from users.
- **Low Overhead Integration:** Fast startup time and minimal resource usage make VibeGuard a low-friction addition to CI/CD and agent workflows.
- **Clear, Maintainable Code:** Go's simplicity and explicit error handling make the codebase easy to understand and extend.
- **Ecosystem Alignment:** Go is the de facto standard for modern DevOps and infrastructure tooling, making VibeGuard fit naturally into that landscape.
- **Strong Tooling:** Go's built-in testing, benchmarking, and profiling tools support high-quality implementation from the start.

### Negative Outcomes

- **Compilation Overhead:** Development iteration requires compilation. Mitigated by fast Go compile times.
- **Smaller Community for Policy Engines:** Go is less commonly used for dynamic policy systems compared to Python. This means fewer pre-built libraries for advanced judge implementations. Mitigated by keeping policy evaluation simple and declarative initially.
- **Learning Curve for Team:** Developers unfamiliar with Go will need onboarding. Mitigated by Go's straightforward syntax and excellent documentation.

### Neutral Impacts

- **Testing Strategy:** Go's built-in testing framework is excellent and will be used throughout the project.
- **Documentation Requirements:** Implementation patterns and policy schema must be clearly documented to guide users and contributors.
- **Dependency Management:** Go modules are the standard; minimal additional complexity.

## Related Decisions

- **ADR-002: Adopt Conventional Commits** — Commit messages will follow conventional commits for semantic versioning and changelog generation
- **VibeGuard Implementation Patterns (Log Entry)** — The choice of Go supports all five implementation patterns (Declarative Policy Runner, Judge-as-a-Policy, Cloud Code Native, Event-Based Policy Graph, Git-Aware Guardrails)

## Next Steps

1. Set up Go module structure and project layout
2. Define initial project dependencies (YAML parser, testing framework, etc.)
3. Prototype the Declarative Policy Runner (Pattern 1) as the initial implementation
4. Establish coding conventions and style guidelines for the project
