# Policy Enforcement Ecosystem Comparison

This guide compares VibeGuard with other popular policy enforcement tools to help you choose the right tool for your needs.

## Overview

Policy enforcement tools ensure code quality, security, and compliance across development workflows. The three main tools compared here serve different use cases:

| Tool | Primary Focus | Policy Language | Scope |
|------|---------------|-----------------|-------|
| **VibeGuard** | CI/CD pipelines, dev workflows | YAML + shell commands | Cross-platform, language-agnostic |
| **OPA (Open Policy Agent)** | Universal policy engine | Rego | Cross-platform, API-centric |
| **Kyverno** | Kubernetes admission control | YAML | Kubernetes-native |

## Quick Decision Guide

**Choose VibeGuard if:**
- You need lightweight CI/CD policy enforcement
- Your policies run shell commands (linters, tests, builds)
- You want zero-dependency single-binary deployment
- You're enforcing code quality across multiple languages
- You need fast startup for frequent invocation (agent loops, pre-commit hooks)

**Choose OPA if:**
- You need a universal policy engine across your entire stack
- Your policies evaluate structured data (JSON, API responses)
- You require fine-grained authorization decisions
- You're willing to learn Rego for maximum flexibility
- You need to integrate with Terraform, APIs, and microservices

**Choose Kyverno if:**
- You're focused exclusively on Kubernetes
- You need admission control, mutation, and resource generation
- You prefer staying within the Kubernetes YAML ecosystem
- You want pre-built policy sets for Pod Security Standards

## Detailed Comparison

### Policy Definition

#### VibeGuard
```yaml
version: "1"
checks:
  - id: test-coverage
    run: go test ./... -coverprofile=cover.out && go tool cover -func cover.out
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
    assert: "coverage >= 80"
    suggestion: "Coverage is {{.coverage}}%, need 80%"
```

VibeGuard policies wrap shell commands with optional output parsing (grok patterns) and assertions. This makes it easy to integrate any existing tool.

#### OPA (Rego)
```rego
package authz

default allow = false

allow {
    input.method == "GET"
    input.path == ["api", "public"]
}

allow {
    input.user.role == "admin"
}
```

OPA uses Rego, a purpose-built policy language for evaluating structured data. Powerful but requires learning a new language.

#### Kyverno
```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-labels
spec:
  validationFailureAction: enforce
  rules:
  - name: check-team-label
    match:
      resources:
        kinds:
        - Pod
    validate:
      message: "Pod must have a 'team' label"
      pattern:
        metadata:
          labels:
            team: "?*"
```

Kyverno uses Kubernetes-native YAML with pattern matching. Excellent for Kubernetes but limited to that ecosystem.

### Execution Model

| Aspect | VibeGuard | OPA | Kyverno |
|--------|-----------|-----|---------|
| **Runtime** | CLI binary | Daemon/library | Kubernetes controller |
| **Invocation** | On-demand | Query-based | Admission webhook |
| **Dependencies** | None (single binary) | Go runtime or container | Kubernetes cluster |
| **Startup time** | ~10ms | ~100ms | Always running |
| **Parallelism** | Built-in (configurable) | Request-based | Per-admission |

### Use Case Alignment

#### CI/CD Pipeline Enforcement

**VibeGuard** excels here:
```yaml
# vibeguard.yaml
checks:
  - id: lint
    run: golangci-lint run ./...
    severity: error

  - id: test
    run: go test ./...
    requires: [lint]

  - id: build
    run: go build ./cmd/...
    requires: [test]
```

**OPA** can do this but requires more setup:
```rego
# policy.rego
package ci

violation[msg] {
    input.lint_exit_code != 0
    msg := "Linting failed"
}
```

**Kyverno** is not designed for CI/CD pipelines.

#### Kubernetes Admission Control

**Kyverno** excels here:
```yaml
# Automatically add security context
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: add-security-context
spec:
  rules:
  - name: add-run-as-non-root
    match:
      resources:
        kinds:
        - Pod
    mutate:
      patchStrategicMerge:
        spec:
          securityContext:
            runAsNonRoot: true
```

**OPA/Gatekeeper** provides similar functionality with more flexibility:
```rego
package kubernetes.admission

deny[msg] {
    input.request.kind.kind == "Pod"
    not input.request.object.spec.securityContext.runAsNonRoot
    msg := "Pods must run as non-root"
}
```

**VibeGuard** is not designed for Kubernetes admission control.

#### API Authorization

**OPA** excels here:
```rego
package httpapi.authz

default allow = false

allow {
    input.method == "GET"
    token_is_valid
    user_owns_resource
}

token_is_valid {
    io.jwt.verify_hs256(input.token, secret)
}

user_owns_resource {
    data.resources[input.resource_id].owner == input.user
}
```

**VibeGuard** and **Kyverno** are not designed for runtime API authorization.

### Learning Curve

| Tool | Time to First Policy | Mastery Time | Prerequisites |
|------|---------------------|--------------|---------------|
| **VibeGuard** | 5 minutes | 1-2 hours | Shell commands, YAML |
| **OPA** | 30 minutes | Days-weeks | Rego language, logic programming |
| **Kyverno** | 15 minutes | Hours | Kubernetes, YAML |

### Integration Points

#### VibeGuard Integrations
- Git pre-commit hooks
- GitHub Actions, GitLab CI, Jenkins, CircleCI
- Claude Code hooks
- Any CI/CD system via exit codes
- Local development workflows

#### OPA Integrations
- Kubernetes (via Gatekeeper)
- Terraform (via Conftest)
- Docker, Envoy, Kafka
- REST APIs (via SDK/library)
- Custom applications

#### Kyverno Integrations
- Kubernetes admission controller
- CLI for testing policies
- Policy Reporter for audit trails
- GitOps workflows (ArgoCD, Flux)

### Performance Characteristics

| Metric | VibeGuard | OPA | Kyverno |
|--------|-----------|-----|---------|
| Cold start | ~10ms | ~100-200ms | N/A (always running) |
| Memory footprint | ~10MB | ~50-100MB | ~100-200MB |
| Policy evaluation | Shell command time | <1ms per query | <10ms per admission |
| Binary size | ~15MB | ~50MB | N/A (container) |

### CNCF Status and Community

| Tool | CNCF Status | GitHub Stars | First Release |
|------|-------------|--------------|---------------|
| **VibeGuard** | N/A | New project | 2024 |
| **OPA** | Graduated | 9k+ | 2016 |
| **Kyverno** | Incubating | 5k+ | 2019 |

## Feature Matrix

| Feature | VibeGuard | OPA | Kyverno |
|---------|-----------|-----|---------|
| YAML-based policies | Yes | No (Rego) | Yes |
| Shell command execution | Yes | No | No |
| Output parsing (grok) | Yes | No | No |
| Arithmetic assertions | Yes | Yes | Limited |
| Dependency ordering | Yes | No | No |
| Parallel execution | Yes | Yes | Yes |
| Fail-fast mode | Yes | N/A | N/A |
| Resource mutation | No | Yes (Gatekeeper) | Yes |
| Resource generation | No | No | Yes |
| Cross-platform | Yes | Yes | Kubernetes only |
| Single binary | Yes | Yes | No |
| AI-assisted setup | Yes | No | No |
| Pre-commit integration | Yes | Via Conftest | No |

## Migration Paths

### From Shell Scripts to VibeGuard

If you have shell scripts enforcing policies:

```bash
# Before: Custom script
#!/bin/bash
go fmt ./... || exit 1
go vet ./... || exit 1
go test ./... || exit 1
```

```yaml
# After: VibeGuard
version: "1"
checks:
  - id: fmt
    run: test -z "$(gofmt -l .)"
  - id: vet
    run: go vet ./...
  - id: test
    run: go test ./...
    requires: [fmt, vet]
```

### Using VibeGuard Alongside OPA

VibeGuard and OPA complement each other:

```yaml
# vibeguard.yaml - CI/CD enforcement
checks:
  - id: opa-policy-test
    run: opa test ./policies -v

  - id: conftest-validate
    run: conftest test deployment.yaml
```

Use VibeGuard to orchestrate OPA policy testing in CI/CD while OPA handles runtime decisions.

### Using VibeGuard Alongside Kyverno

```yaml
# vibeguard.yaml - Pre-deployment checks
checks:
  - id: kyverno-test
    run: kyverno test ./policies

  - id: kubectl-validate
    run: kubectl apply --dry-run=server -f manifests/
```

Use VibeGuard for CI/CD pipeline gates while Kyverno handles Kubernetes admission control.

## Recommendations by Team Size

### Small Teams / Solo Developers
**Recommendation: VibeGuard**
- Minimal setup overhead
- No infrastructure requirements
- Easy to understand and maintain

### Medium Teams (5-20 developers)
**Recommendation: VibeGuard + Kyverno (if using Kubernetes)**
- VibeGuard for CI/CD and local development
- Kyverno for Kubernetes guardrails
- Low operational overhead

### Large Organizations
**Recommendation: VibeGuard + OPA + Kyverno**
- VibeGuard for CI/CD orchestration
- OPA for cross-platform policy decisions
- Kyverno for Kubernetes-specific needs
- Unified policy testing via VibeGuard

## Summary

| When you need... | Use... |
|-----------------|--------|
| CI/CD policy enforcement | VibeGuard |
| Local development checks | VibeGuard |
| Pre-commit hooks | VibeGuard |
| Kubernetes admission control | Kyverno |
| Kubernetes resource mutation | Kyverno |
| API authorization | OPA |
| Terraform policy checks | OPA (Conftest) |
| Cross-platform policy engine | OPA |
| Simple, fast, no dependencies | VibeGuard |

Each tool has its strengths. VibeGuard fills the gap for lightweight, CI/CD-focused policy enforcement that complements rather than competes with OPA and Kyverno.

## Further Reading

- [VibeGuard Documentation](../README.md)
- [OPA Documentation](https://www.openpolicyagent.org/docs/)
- [Kyverno Documentation](https://kyverno.io/docs/)
- [Conftest (OPA for CI/CD)](https://www.conftest.dev/)
- [Gatekeeper (OPA for Kubernetes)](https://open-policy-agent.github.io/gatekeeper/)
