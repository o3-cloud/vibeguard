# Integration Guides

This guide shows how to integrate VibeGuard into various CI/CD platforms and development workflows.

## Table of Contents

1. [GitHub Actions](#github-actions)
2. [GitLab CI](#gitlab-ci)
3. [Jenkins](#jenkins)
4. [CircleCI](#circleci)
5. [Travis CI](#travis-ci)
6. [Git Pre-commit Hook](#git-pre-commit-hook)
7. [Local Development](#local-development)

## GitHub Actions

### Basic Workflow

Create `.github/workflows/vibeguard.yml`:

```yaml
name: VibeGuard Quality Checks

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  vibeguard:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Install VibeGuard
        run: |
          curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard
          sudo mv vibeguard /usr/local/bin/

      - name: Run VibeGuard
        run: vibeguard check
```

### With Setup for Go Projects

For projects using Go:

```yaml
name: VibeGuard Quality Checks

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  vibeguard:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Install VibeGuard
        run: |
          curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard
          sudo mv vibeguard /usr/local/bin/

      - name: Run VibeGuard
        run: vibeguard check
```

### With Setup for Node.js Projects

For projects using Node.js:

```yaml
name: VibeGuard Quality Checks

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  vibeguard:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci

      - name: Install VibeGuard
        run: |
          curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard
          sudo mv vibeguard /usr/local/bin/

      - name: Run VibeGuard
        run: vibeguard check
```

### With JSON Output for Reporting

```yaml
name: VibeGuard Quality Checks

on: [push, pull_request]

jobs:
  vibeguard:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Install VibeGuard
        run: |
          curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard
          sudo mv vibeguard /usr/local/bin/

      - name: Run VibeGuard
        run: |
          vibeguard check -v --json | tee vibeguard-results.json

      - name: Upload results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: vibeguard-results
          path: vibeguard-results.json
```

### Matrix Strategy (Multiple Versions)

```yaml
name: VibeGuard Quality Checks

on: [push, pull_request]

jobs:
  vibeguard:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.21', '1.24']

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}

      - name: Install VibeGuard
        run: |
          # Installation varies by OS - this is simplified
          curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
          chmod +x vibeguard
          sudo mv vibeguard /usr/local/bin/

      - name: Run VibeGuard
        run: vibeguard check -v
```

## GitLab CI

### Basic Pipeline

Create `.gitlab-ci.yml`:

```yaml
vibeguard:
  stage: test
  image: golang:1.24-alpine

  before_script:
    - apk add --no-cache curl bash
    - |
      curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
      chmod +x vibeguard
      sudo mv vibeguard /usr/local/bin/

  script:
    - vibeguard check

  only:
    - merge_requests
    - main
    - develop
```

### With Artifacts

```yaml
vibeguard:
  stage: test
  image: golang:1.24-alpine

  before_script:
    - apk add --no-cache curl bash
    - |
      curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
      chmod +x vibeguard
      sudo mv vibeguard /usr/local/bin/

  script:
    - vibeguard check -v --json | tee vibeguard-results.json

  artifacts:
    name: vibeguard-results
    paths:
      - vibeguard-results.json
      - .vibeguard/log/
    when: always
    expire_in: 30 days

  allow_failure: false
```

### Multi-stage Pipeline

```yaml
stages:
  - lint
  - test
  - check

lint:
  stage: lint
  image: golang:1.24-alpine
  script:
    - go fmt ./...
    - go vet ./...

test:
  stage: test
  image: golang:1.24-alpine
  script:
    - go test ./...

vibeguard:
  stage: check
  image: golang:1.24-alpine

  before_script:
    - apk add --no-cache curl bash
    - |
      curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
      chmod +x vibeguard
      sudo mv vibeguard /usr/local/bin/

  script:
    - vibeguard check -v

  needs:
    - lint
    - test
```

## Jenkins

### Declarative Pipeline

Create `Jenkinsfile`:

```groovy
pipeline {
    agent any

    stages {
        stage('Install VibeGuard') {
            steps {
                sh '''
                    curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
                    chmod +x vibeguard
                    sudo mv vibeguard /usr/local/bin/
                '''
            }
        }

        stage('Run Quality Checks') {
            steps {
                sh 'vibeguard check -v'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: '.vibeguard/log/**', allowEmptyArchive: true
        }
        failure {
            echo 'VibeGuard quality checks failed'
        }
    }
}
```

### Scripted Pipeline

```groovy
node {
    try {
        stage('Install') {
            sh '''
                curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
                chmod +x vibeguard
                sudo mv vibeguard /usr/local/bin/
            '''
        }

        stage('Checkout') {
            checkout scm
        }

        stage('Quality Check') {
            sh 'vibeguard check --fail-fast'
        }

    } catch (Exception e) {
        currentBuild.result = 'FAILURE'
        throw e
    } finally {
        archiveArtifacts artifacts: '.vibeguard/log/**', allowEmptyArchive: true
    }
}
```

## CircleCI

### Configuration

Create `.circleci/config.yml`:

```yaml
version: 2.1

jobs:
  vibeguard:
    docker:
      - image: golang:1.24-alpine

    steps:
      - checkout

      - run:
          name: Install VibeGuard
          command: |
            apk add --no-cache curl bash
            curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
            chmod +x vibeguard
            sudo mv vibeguard /usr/local/bin/

      - run:
          name: Run Quality Checks
          command: vibeguard check -v

      - store_artifacts:
          path: .vibeguard/log/
          destination: vibeguard-logs

workflows:
  version: 2
  test:
    jobs:
      - vibeguard:
          filters:
            branches:
              only:
                - main
                - develop
```

## Travis CI

### Configuration

Create `.travis.yml`:

```yaml
language: go

go:
  - '1.24'

before_script:
  - |
    curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
    chmod +x vibeguard
    sudo mv vibeguard /usr/local/bin/

script:
  - vibeguard check

after_success:
  - echo "Quality checks passed!"

after_failure:
  - echo "Quality checks failed!"
```

## Git Pre-commit Hook

Integrate VibeGuard into your local development workflow with git hooks.

### Manual Setup

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash

# Run VibeGuard before committing
vibeguard check

# Capture exit code
exit_code=$?

# Exit with VibeGuard's exit code
exit $exit_code
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

### Using Pre-commit Framework

If you use [pre-commit](https://pre-commit.com/), create `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/yourusername/vibeguard
    rev: v1.0.0
    hooks:
      - id: vibeguard
        name: VibeGuard Quality Checks
        entry: vibeguard check
        language: system
        types: [text]
        stages: [commit]
```

Install the hook:
```bash
pre-commit install
```

### Husky Setup (Node.js Projects)

For Node.js projects using Husky:

```bash
npm install husky --save-dev
npx husky install
```

Create `.husky/pre-commit`:

```bash
#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

vibeguard check --fail-fast
```

Make it executable:
```bash
chmod +x .husky/pre-commit
```

## Local Development

### Development Makefile

Add to your `Makefile`:

```makefile
.PHONY: check
check:
	vibeguard check

.PHONY: check-verbose
check-verbose:
	vibeguard check -v

.PHONY: check-list
check-list:
	vibeguard list

.PHONY: check-validate
check-validate:
	vibeguard validate

.PHONY: check-quick
check-quick:
	vibeguard check --fail-fast

.PHONY: check-logs
check-logs:
	tail -f .vibeguard/log/*.log
```

Usage:
```bash
make check          # Run all checks
make check-verbose  # Run with verbose output
make check-quick    # Stop on first failure
make check-logs     # View logs
```

### Shell Alias

Add to your `.bashrc` or `.zshrc`:

```bash
alias vg="vibeguard check"
alias vg-v="vibeguard check -v"
alias vg-list="vibeguard list"
alias vg-validate="vibeguard validate"
```

Usage:
```bash
vg          # Quick check
vg-v        # Verbose check
vg-list     # List checks
vg-validate # Validate config
```

### IDE Integration

#### VS Code

Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "VibeGuard: Run All Checks",
      "type": "shell",
      "command": "vibeguard",
      "args": ["check"],
      "problemMatcher": []
    },
    {
      "label": "VibeGuard: Run with Verbose Output",
      "type": "shell",
      "command": "vibeguard",
      "args": ["check", "-v"],
      "problemMatcher": []
    },
    {
      "label": "VibeGuard: Run Specific Check",
      "type": "shell",
      "command": "vibeguard",
      "args": ["check", "${input:checkId}"],
      "problemMatcher": [],
      "inputs": [
        {
          "id": "checkId",
          "type": "promptString",
          "description": "Enter check ID",
          "default": "test"
        }
      ]
    }
  ]
}
```

Then use `Cmd+Shift+P` → "Run Task" to run VibeGuard from VS Code.

#### IntelliJ/GoLand

Create `.idea/runConfigurations/VibeGuard.xml`:

```xml
<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="VibeGuard" type="GoApplicationRunConfiguration" factoryName="Go Application">
    <option name="CMD_LINE_ARGS" value="check" />
    <option name="WORKING_DIRECTORY" value="$PROJECT_DIR$" />
    <method v="2" />
  </configuration>
</component>
```

Then select "VibeGuard" in the run dropdown and click Run.

### Docker Integration

#### Dockerfile with VibeGuard

```dockerfile
FROM golang:1.24-alpine as builder

# Build stage
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/main.go

# Install VibeGuard
RUN apk add --no-cache curl bash && \
    curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard && \
    chmod +x vibeguard

# Test stage - Run VibeGuard
RUN ./vibeguard check

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/app /usr/local/bin/
ENTRYPOINT ["app"]
```

#### Docker Compose

```yaml
version: '3'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    volumes:
      - .:/app
    command: vibeguard check -v
```

Run with:
```bash
docker-compose run app
```

## Best Practices

### CI/CD Integration

1. **Run on every push** - Catch issues early
2. **Use `--fail-fast`** - Get quick feedback
3. **Run in parallel** - Default parallelism is fine for most projects
4. **Save artifacts** - Keep logs for debugging
5. **Block PRs** - Set CI to fail if checks don't pass

### Local Development

1. **Use pre-commit hook** - Prevent commits with violations
2. **Run with `-v` flag** - See all results, not just failures
3. **Run specific checks** - Speed up iteration
4. **Check logs** - Look in `.vibeguard/log/` for details

### Configuration Management

1. **Keep one config file** - `vibeguard.yaml` in repo root
2. **Version your config** - Track changes in git
3. **Use variables** - Make config reusable
4. **Document custom checks** - Help team understand requirements

### Troubleshooting CI Failures

1. **Check exit codes** - VibeGuard returns specific codes
2. **Review logs** - Check `.vibeguard/log/` for details
3. **Run locally** - Reproduce before debugging in CI
4. **Increase timeout** - Some checks may need more time in CI
5. **Check dependencies** - Ensure tools are installed in CI environment
