# VibeGuard CI/CD Integration Guide

This guide shows how to integrate VibeGuard into popular CI/CD platforms and development workflows. VibeGuard's lightweight, single-binary design makes it easy to add to any pipeline.

## Quick Start for CI/CD

Every integration follows the same basic pattern:

1. **Download or build VibeGuard** — Get the binary into your pipeline
2. **Add your config** — Include `vibeguard.yaml` in version control
3. **Run checks** — Execute `vibeguard check` as a pipeline step
4. **Handle exit codes** — Use exit codes to pass/fail the pipeline step

## GitHub Actions

### Basic Integration

Add a step to your workflow file:

```yaml
name: Quality Checks
on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download VibeGuard
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard

      - name: Run VibeGuard checks
        run: ./vibeguard check
        # Automatically fails the job if exit code >= 1
```

### With Go Setup

For Go projects, set up the Go environment first:

```yaml
name: Go Quality Gate
on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.21', '1.22']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}

      - name: Download VibeGuard
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard

      - name: Install dependencies
        run: go mod tidy

      - name: Run VibeGuard checks
        run: ./vibeguard check --verbose
```

### With Node.js Setup

For JavaScript/TypeScript projects:

```yaml
name: Node Quality Gate
on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci

      - name: Download VibeGuard
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard

      - name: Run VibeGuard checks
        run: ./vibeguard check --json > results.json

      - name: Upload results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: vibeguard-results
          path: results.json
```

### Conditional Checks Based on File Changes

Run specific checks only when relevant files change:

```yaml
name: Smart Quality Gate
on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for git diff

      - name: Check modified files
        id: files
        run: |
          git diff --name-only origin/main...HEAD > changed-files.txt
          if grep -q '\.go$' changed-files.txt; then echo "go-changes=true" >> $GITHUB_OUTPUT; fi
          if grep -q '\.ts\|\.js$' changed-files.txt; then echo "js-changes=true" >> $GITHUB_OUTPUT; fi

      - name: Download VibeGuard
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard

      - name: Run Go checks
        if: steps.files.outputs.go-changes == 'true'
        run: ./vibeguard check go-fmt go-vet go-test

      - name: Run JavaScript checks
        if: steps.files.outputs.js-changes == 'true'
        run: ./vibeguard check js-lint js-test
```

### Multi-Platform Matrix

Test on multiple operating systems:

```yaml
name: Multi-Platform QA
on: [push, pull_request]

jobs:
  quality:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.21', '1.22']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}

      - name: Download VibeGuard (Linux)
        if: runner.os == 'Linux'
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard

      - name: Download VibeGuard (macOS)
        if: runner.os == 'macOS'
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-darwin-amd64 -o vibeguard
          chmod +x vibeguard

      - name: Download VibeGuard (Windows)
        if: runner.os == 'Windows'
        run: |
          curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-windows-amd64.exe -o vibeguard.exe

      - name: Run VibeGuard checks
        run: ./vibeguard check --verbose
```

## GitLab CI

### Basic Integration

Add to your `.gitlab-ci.yml`:

```yaml
image: golang:1.22

quality:
  stage: test
  before_script:
    - curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
    - chmod +x vibeguard
  script:
    - ./vibeguard check
  artifacts:
    when: always
    reports:
      junit: results.xml
```

### With Multiple Stages

```yaml
image: golang:1.22

stages:
  - quality
  - test
  - build

vibeguard_check:
  stage: quality
  before_script:
    - apt-get update && apt-get install -y curl
    - curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
    - chmod +x vibeguard
  script:
    - ./vibeguard check --fail-fast
  allow_failure: false

unit_tests:
  stage: test
  needs: ["vibeguard_check"]
  script:
    - go test -v ./... -coverprofile=cover.out

build:
  stage: build
  needs: ["unit_tests"]
  script:
    - go build -o myapp ./cmd
```

### Conditional Execution

```yaml
image: golang:1.22

vibeguard_check:
  stage: quality
  before_script:
    - curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
    - chmod +x vibeguard
  script:
    - ./vibeguard check
  only:
    - merge_requests
    - main
    - develop
```

## Jenkins

### Declarative Pipeline

```groovy
pipeline {
    agent any

    environment {
        VIBEGUARD_VERSION = "1.0.0"
        VIBEGUARD_URL = "https://github.com/vibeguard/vibeguard/releases/download/v${VIBEGUARD_VERSION}"
    }

    stages {
        stage('Setup') {
            steps {
                script {
                    sh '''
                        if [ ! -f "vibeguard" ]; then
                            curl -L ${VIBEGUARD_URL}/vibeguard-linux-amd64 -o vibeguard
                            chmod +x vibeguard
                        fi
                    '''
                }
            }
        }

        stage('Quality Checks') {
            steps {
                sh './vibeguard check --verbose'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./... -coverprofile=cover.out'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o app ./cmd'
            }
        }
    }

    post {
        always {
            // Publish results
            junit 'test-results.xml'
        }
        failure {
            echo 'Quality checks failed!'
        }
    }
}
```

### Scripted Pipeline

```groovy
node {
    try {
        stage('Checkout') {
            checkout scm
        }

        stage('Download VibeGuard') {
            sh '''
                curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
                chmod +x vibeguard
            '''
        }

        stage('Quality') {
            sh './vibeguard check'
        }

        stage('Test') {
            sh 'go test ./...'
        }

        stage('Build') {
            sh 'go build ./cmd'
        }

        currentBuild.result = 'SUCCESS'
    } catch (e) {
        currentBuild.result = 'FAILURE'
        throw e
    }
}
```

## CircleCI

### Basic Configuration

```yaml
version: 2.1

jobs:
  quality:
    docker:
      - image: cimg/go:1.22
    steps:
      - checkout
      - run:
          name: Download VibeGuard
          command: |
            curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
            chmod +x vibeguard
      - run:
          name: Run VibeGuard checks
          command: ./vibeguard check --verbose

  test:
    docker:
      - image: cimg/go:1.22
    steps:
      - checkout
      - run:
          name: Run tests
          command: go test ./... -coverprofile=cover.out

  build:
    docker:
      - image: cimg/go:1.22
    steps:
      - checkout
      - run:
          name: Build
          command: go build -o app ./cmd

workflows:
  quality_and_test:
    jobs:
      - quality
      - test:
          requires:
            - quality
      - build:
          requires:
            - test
```

## Git Pre-Commit Hook

Integrate VibeGuard into your local development workflow using git hooks.

### Setup

```bash
# Create the hook directory if it doesn't exist
mkdir -p .git/hooks

# Create a pre-commit hook script
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
set -e

# Download VibeGuard if not present
if [ ! -f "vibeguard" ]; then
    curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
    chmod +x vibeguard
fi

# Run VibeGuard checks
if ! ./vibeguard check; then
    echo "VibeGuard checks failed. Commit blocked."
    exit 1
fi
EOF

# Make it executable
chmod +x .git/hooks/pre-commit
```

### Using pre-commit Framework

Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: vibeguard
        name: VibeGuard Quality Checks
        entry: vibeguard check
        language: system
        pass_filenames: false
        always_run: true
        stages: [commit]
```

Then install and run:

```bash
pip install pre-commit
pre-commit install
pre-commit run --all-files
```

## Bitbucket Pipelines

### Basic Pipeline

```yaml
image: golang:1.22

pipelines:
  default:
    - step:
        name: Quality Checks
        script:
          - curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
          - chmod +x vibeguard
          - ./vibeguard check --verbose

    - step:
        name: Test
        script:
          - go test ./... -coverprofile=cover.out

    - step:
        name: Build
        script:
          - go build -o app ./cmd

  pull-requests:
    '**':
      - step:
          name: Quality & Tests
          script:
            - curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
            - chmod +x vibeguard
            - ./vibeguard check --fail-fast
            - go test ./...
```

## Azure Pipelines

### Basic Pipeline

```yaml
trigger:
  - main
  - develop

pool:
  vmImage: 'ubuntu-latest'

variables:
  goVersion: '1.22'

stages:
  - stage: Quality
    jobs:
      - job: RunVibeGuard
        steps:
          - task: GoTool@0
            inputs:
              version: $(goVersion)

          - script: |
              curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard
              chmod +x vibeguard
            displayName: 'Download VibeGuard'

          - script: ./vibeguard check --verbose
            displayName: 'Run VibeGuard Checks'

  - stage: Test
    dependsOn: Quality
    condition: succeeded()
    jobs:
      - job: RunTests
        steps:
          - task: GoTool@0
            inputs:
              version: $(goVersion)

          - script: go test ./... -coverprofile=cover.out
            displayName: 'Run Tests'
```

## Docker

### Dockerfile Integration

Build and run VibeGuard in a Docker container:

```dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Copy project files
COPY . .

# Download VibeGuard
RUN curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o /usr/local/bin/vibeguard && \
    chmod +x /usr/local/bin/vibeguard

# Run VibeGuard checks
RUN vibeguard check

# Build application
RUN go build -o app ./cmd

ENTRYPOINT ["./app"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  quality:
    image: golang:1.22-alpine
    volumes:
      - .:/app
    working_dir: /app
    command: |
      sh -c "
        curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o vibeguard && \
        chmod +x vibeguard && \
        ./vibeguard check
      "

  test:
    image: golang:1.22-alpine
    volumes:
      - .:/app
    working_dir: /app
    depends_on:
      - quality
    command: go test ./...

  build:
    image: golang:1.22-alpine
    volumes:
      - .:/app
    working_dir: /app
    depends_on:
      - test
    command: go build -o app ./cmd
```

## Cloud Code Hooks

Integrate VibeGuard with Claude Code using git pre-commit hooks:

### Setup

1. Install the hook:

```bash
mkdir -p .git/hooks
curl -L https://your-repo/scripts/install-vibeguard-hook.sh | bash
```

2. Configure in `.claude/settings.json`:

```json
{
  "hooks": {
    "pre_commit": {
      "command": "./vibeguard check --fail-fast",
      "enabled": true,
      "block_on_failure": true
    }
  }
}
```

### Example Hook Script

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Ensure VibeGuard is available
if ! command -v vibeguard &> /dev/null; then
    echo "VibeGuard not found. Installing..."
    curl -L https://github.com/vibeguard/vibeguard/releases/download/v1.0.0/vibeguard-linux-amd64 -o /usr/local/bin/vibeguard
    chmod +x /usr/local/bin/vibeguard
fi

# Run checks
vibeguard check --fail-fast

exit $?
```

## Troubleshooting

### Exit Code Reference

- **0** — All checks passed
- **2** — Configuration error (invalid YAML, validation failure)
- **3** — One or more violations detected
- **4** — Timeout or execution error

### Common Issues

**Issue: "vibeguard: command not found"**

Ensure the binary is in your PATH or use the full path:

```bash
./vibeguard check
# or
/usr/local/bin/vibeguard check
```

**Issue: "Permission denied" on executable**

Make the binary executable:

```bash
chmod +x vibeguard
```

**Issue: Timeout during checks**

Increase the timeout value in `vibeguard.yaml`:

```yaml
checks:
  - id: test
    run: go test ./...
    timeout: 120s  # Increase from default 30s
```

**Issue: Different behavior locally vs. CI**

Ensure:
- Go/Node version matches CI environment
- All dependencies are installed
- Environment variables are set correctly

Run `vibeguard validate` to check configuration:

```bash
vibeguard validate
```

## Performance Tips

### Parallel Execution

By default, VibeGuard runs 4 checks in parallel. Adjust for your environment:

```bash
vibeguard check --parallel 8
```

Or in your config:

```yaml
# In GitHub Actions
- name: Run VibeGuard checks
  run: ./vibeguard check --parallel 16
```

### Fail-Fast Mode

Stop on first failure for faster feedback:

```bash
vibeguard check --fail-fast
```

### Selective Checks

Run only specific checks to save time:

```bash
vibeguard check fmt vet  # Only run fmt and vet checks
```

### JSON Output for Processing

Use JSON output for integration with other tools:

```bash
vibeguard check --json | jq '.violations[] | select(.severity == "error")'
```

## See Also

- [Configuration Schema](../README.md#configuration-schema)
- [CLI Reference](../README.md#cli-reference)
- [Exit Codes](../README.md#exit-codes)
- [Examples](../examples/)
