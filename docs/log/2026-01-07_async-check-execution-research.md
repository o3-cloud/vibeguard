---
summary: Research async execution patterns for long-running checks (5-10 minutes) to prevent blocking the orchestrator
event_type: research
sources:
  - golang.org/x/sync/errgroup (current implementation)
  - context package (cancellation and deadline handling)
  - Go channels and goroutines documentation
  - Internal orchestrator implementation analysis
tags:
  - async
  - architecture
  - long-running-checks
  - orchestrator
  - design-decision
  - performance
  - concurrency
  - scalability
---

# Async Check Execution for Long-Running Tasks

## Context

VibeGuard currently executes checks synchronously within a level-by-level dependency graph. The current architecture uses:
- **errgroup.WithContext()** for goroutine management
- **Semaphore pattern** to limit parallel execution (`--parallel` flag)
- **Per-check timeouts** (default 30s)
- **Fail-fast mechanism** to cancel remaining checks on error

**Problem Statement:** Long-running checks (5-10 minutes, e.g., integration tests, E2E test suites, security scans, performance benchmarks) can:
1. Block the entire orchestrator process
2. Block dependent checks unnecessarily
3. Waste resources if not all checks are needed (e.g., pre-commit scenarios)
4. Make it difficult to integrate VibeGuard into LLM-based tools that need responsive feedback

## Current Execution Model

```
Synchronous Level-by-Level:
Level 1: [check-a, check-b, check-c] (parallel within level)
  ↓ (wait for ALL to complete)
Level 2: [check-d] (depends on level 1)
  ↓ (wait)
Level 3: [check-e]
```

**Current behavior:** Orchestrator blocks until check completes (up to timeout duration)

## Long-Running Check Use Cases

1. **Integration Tests** (5-10 min): Full integration test suite with multiple services
2. **E2E Test Suites** (5-10 min): Browser-based end-to-end tests
3. **Security Scans** (5-10 min): SAST/DAST tools, vulnerability scanning
4. **Performance Benchmarks** (5-10 min): Load testing, stress testing
5. **Database Migrations** (2-5 min): Large dataset migrations/schema changes
6. **Docker Image Building** (3-10 min): Building multi-stage Docker images
7. **Comprehensive Code Analysis** (5-10 min): SonarQube, code climate analysis
8. **Remote API Tests** (5-10 min): API testing against staging environments

## Async Execution Patterns

### Pattern 1: Fire-and-Forget with Status Polling

**How it works:**
- Checks spawn as background goroutines
- Orchestrator returns immediately with "in-progress" status
- Separate polling mechanism queries results
- Dependent checks wait for results

**Pros:**
- Orchestrator returns immediately (non-blocking)
- Simple to implement
- Familiar pattern

**Cons:**
- Complex state management
- Polling overhead and latency
- Dependent checks need retry logic
- Difficult to handle cancellation uniformly
- No clean exit semantics (process could exit with running checks)

**Use case:** When you want responsive LLM feedback but don't need final results immediately

---

### Pattern 2: Async Checks with Explicit Await

**How it works:**
- Long-running checks marked with `async: true` in config
- Orchestrator marks as "pending" and continues to dependent checks
- Dependent checks either:
  - **Option A**: Block until async check completes (defeats purpose)
  - **Option B**: Execute conditionally (skip if async not done)
  - **Option C**: Explicitly await async checks when needed

**Example config:**
```yaml
- id: integration-tests
  run: go test ./tests/integration -timeout 10m
  async: true
  timeout: 12m
  tags: [slow, async]

- id: deploy-check
  requires: [integration-tests]  # Explicit dependency
  run: ./scripts/validate-deploy.sh
```

**Pros:**
- Explicit control over which checks are async
- Clear dependency semantics
- Can mix sync and async checks
- Orchestrator can report progress

**Cons:**
- More configuration complexity
- Need to rethink dependency graph (different semantics)
- Result availability timing creates edge cases
- Dependent check behavior ambiguous

**Use case:** When some checks should be async but others need to block

---

### Pattern 3: Fire-and-Forget Queue System

**How it works:**
- Long-running checks written to a job queue (file, database, or in-memory)
- Separate worker processes consume queue
- Orchestrator returns with "queued" status
- Results fetched on demand or via webhook

**Implementation:**
```go
type QueuedCheck struct {
    CheckID   string
    Config    Check
    QueuedAt  time.Time
    JobID     string
}

// Orchestrator writes to queue instead of executing
if check.Async {
    queuedCheck := QueuedCheck{CheckID: check.ID, ...}
    queue.Enqueue(queuedCheck)
    result.Status = "queued"
} else {
    // Execute synchronously
}
```

**Pros:**
- True non-blocking execution
- Orchestrator never waits
- Scalable to distributed systems
- Can replay checks
- Clean separation of concerns

**Cons:**
- Much higher complexity
- Requires external storage/worker infrastructure
- Result polling or webhook complexity
- Fail-fast semantics broken (can't cascade failures)
- Dependency graph becomes much more complex

**Use case:** Enterprise scenarios with dedicated CI/CD infrastructure

---

### Pattern 4: Hybrid Sync/Async with Context Multiplexing

**How it works:**
- Per-check timeout determines if async needed
- Checks with `timeout > 60s` or `async: true` spawn as background tasks
- Main orchestrator continues with dependent checks
- At critical points (before exit, before dependent checks), wait for results

**Implementation:**
```go
type CheckExecution struct {
    Check    *Check
    Result   chan *executor.Result
    Cancel   context.CancelFunc
    Task     *async.Task
}

// Execute check
if check.Timeout > 60*time.Second {
    // Async: spawn and continue
    execution := executeAsync(ctx, check)
    pendingChecks[check.ID] = execution
} else {
    // Sync: wait for result
    result := execute(ctx, check)
}

// Before returning, ensure all pending checks complete (with timeout)
waitForPendingChecks(ctx, pendingChecks, time.Second*30)
```

**Pros:**
- Automatic based on timeout config (no new YAML keywords)
- Backward compatible
- Scales naturally (short checks stay responsive)
- Handles dependencies elegantly
- Can still enforce final results before exit

**Cons:**
- More goroutine management complexity
- Race conditions possible if not careful
- Harder to debug
- Still blocks on critical paths (dependencies)

**Use case:** General-purpose solution that works for most scenarios

---

### Pattern 5: Streaming Results with Cancellable Background Workers

**How it works:**
- Checks run in background worker pool
- Results streamed via channel as they complete
- Consumer can process results immediately
- Dependent checks handled via result stream filtering
- Early exit stops workers gracefully

**Implementation:**
```go
type CheckResult struct {
    Check   *Check
    Result  *executor.Result
    Err     error
}

// Returns channel that streams results as they complete
results := orchestrator.RunAsync(ctx, checks)

for result := range results {
    if result.Err != nil {
        handleError(result)
        continue
    }

    if shouldFailFast(result) {
        cancel()  // Stop all remaining workers
    }

    processResult(result)
}
```

**Pros:**
- True streaming/reactive approach
- Responsive feedback (results available immediately)
- Natural backpressure handling
- Elegant cancellation (context-based)
- Composable with other async operations

**Cons:**
- Requires consumer to handle streaming
- API changes significantly
- Dependency handling complex (must track pending)
- Not good for final exit code computation (need to wait for all)

**Use case:** Integration with LLM tools that want streaming feedback

---

## Recommendation for VibeGuard

### Short-term: Pattern 4 (Hybrid Sync/Async with Context Multiplexing)

**Rationale:**
- Minimal YAML config changes (users don't need to learn new keywords)
- Automatic: long timeouts = async behavior
- Backward compatible: existing configs work unchanged
- Solves the LLM blocking problem
- Maintains dependency semantics
- Handles fail-fast correctly

**Implementation sketch:**
1. Add `async.ExecutionPool` to orchestrator (manages background tasks)
2. Modify `orchestrator.executeCheck()` to spawn async if timeout > threshold
3. At critical points (dependency satisfaction, exit), wait for pending async results
4. Maintain compatibility with `--fail-fast` flag
5. Extend `executor.Result` with `Status` field: "pending", "completed", "timeout", "cancelled"

**New behavior:**
```
Checks with timeout <= 60s:  Synchronous (current behavior)
Checks with timeout > 60s:   Async (spawns, returns with Status=pending)
Dependent checks:            Block until dependency Status != pending
Process exit:                Waits up to timeout for pending checks
Fail-fast:                   Cancels all pending + running checks
```

---

### Medium-term: Pattern 5 (Streaming Results)

**Rationale:**
- For better LLM integration
- CLI can display live progress
- More reactive and responsive

**Implementation:**
1. Add `vibeguard check --stream` flag
2. Returns channel instead of waiting for all results
3. Consumer processes results as they arrive
4. Still respects dependencies (filters stream)

---

## Next Steps

1. **Decide on pattern** - Recommend Pattern 4 for initial implementation
2. **Design config** - How to mark checks as async? (timeout-based vs explicit flag?)
3. **Implement** - Add async execution to orchestrator
4. **Test** - Add tests for:
   - Async check completion
   - Dependent check blocking on async
   - Fail-fast with async checks
   - Timeout/cancellation of async checks
   - Exit behavior (waiting for pending)
5. **Document** - Update README with async capabilities
6. **Consider ADR** - If significant architectural change, document decision

## Related Decisions

- See ADR-003 (Go implementation) for concurrency patterns available
- See ADR-005 (VibeGuard for policy enforcement) for use cases
- Current implementation uses `errgroup.WithContext()` and semaphore pattern

