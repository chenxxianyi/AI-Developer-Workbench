# Go Backend Performance Skill Design

## Goal

Create a project-local Codex skill at `.codex/skills/go-backend-performance`.
The skill analyzes Go backend performance using a fixed evidence chain:

`Benchmark → pprof → slow SQL → connection pool → concurrency/locks → cache → before/after comparison`

The first invocation is analysis-only. It must produce a reviewable report and stop before changing code.

## Trigger Scope

Use the skill for Go backend performance reviews involving:

- slow APIs, high latency, low throughput, or high CPU/memory use;
- Go benchmarks, `pprof`, execution traces, allocations, goroutines, mutexes, or blocking;
- `database/sql`, GORM, MySQL, slow queries, N+1 queries, indexes, or connection pools;
- contention, excessive serialization, duplicate work, or cache opportunities;
- validating whether a proposed optimization produced a measurable improvement.

Do not use it as a generic correctness review, security audit, or architecture redesign unless those concerns directly affect measured performance.

## Structure

The skill contains:

- `SKILL.md`: mandatory workflow, safety gates, evidence rules, and reporting requirements;
- `references/go-profiling.md`: benchmark, `pprof`, trace, race, mutex, and block profiling guidance;
- `references/database-performance.md`: slow SQL, query plans, GORM, MySQL, and connection-pool checks;
- `assets/performance-report-template.md`: standardized report copied into `docs/performance/`;
- `agents/openai.yaml`: UI metadata and default invocation prompt.

No executable scripts are required initially. Go projects vary in routing, workload generation, database access, and deployment topology; the skill should discover and use the project's native commands instead of imposing a fragile wrapper.

## Two-Phase Authorization Model

### Phase A: Analysis and Report

Phase A is the default and is strictly read-only with respect to application source, configuration, dependencies, schemas, and persistent data.

Allowed actions:

- inspect code, configuration examples, tests, migrations, logs, and documentation;
- run existing tests and benchmarks that do not mutate persistent data;
- collect local CPU, heap, allocation, goroutine, mutex, block, and trace profiles when a safe existing target is available;
- inspect existing slow-query evidence and read-only query plans;
- read connection-pool configuration and metrics;
- create temporary profiling artifacts under the system temporary directory;
- write only the final report under `docs/performance/`.

Forbidden actions:

- modify application code, tests, benchmarks, configuration, environment files, migrations, or dependencies;
- enable production profiling endpoints or expose debug endpoints;
- run write queries, migrations, destructive load tests, or tests against an unverified database;
- install profiling or load-testing tools without explicit approval;
- claim an optimization result before a repeatable before/after measurement exists.

If a required benchmark, profile endpoint, workload, database permission, or metric is missing, record the evidence gap and the minimum proposed instrumentation. Do not add it during Phase A.

At the end of Phase A, write the report and stop with an explicit approval request. Suggestions such as “go ahead,” “fix it,” or “implement the approved items” in a later user message authorize Phase B only for the named or clearly selected findings.

### Phase B: Approved Implementation and Validation

Enter Phase B only after explicit user approval following report review.

For each approved finding:

1. preserve or create a repeatable baseline;
2. use test-driven development for behavior changes and regression coverage;
3. make the smallest targeted change;
4. run correctness, race, and relevant performance checks;
5. repeat measurements under comparable conditions;
6. report raw before/after values, variance, trade-offs, and remaining uncertainty.

Approval is scoped. Do not implement unapproved recommendations merely because they appear in the same report.

## Fixed Analysis Workflow

The skill must execute the stages in order. A stage may be marked “blocked” or “not applicable,” but it may not be silently skipped.

### 1. Benchmark

- Identify existing unit, integration, endpoint, and benchmark commands.
- Establish reproducible workload, environment, sample count, and key metrics.
- Prefer `go test -bench` with allocation reporting when suitable.
- Record latency distributions or request throughput for endpoint workloads.
- Separate cold-start, external API, disk, and database effects when possible.

### 2. pprof

- Use CPU and heap/allocation profiles first.
- Inspect goroutine, mutex, and block profiles when symptoms justify them.
- Use execution trace only when scheduler, GC, blocking, or goroutine behavior requires it.
- Correlate hot paths with the benchmark workload; do not optimize code merely because it appears in a static scan.

### 3. Slow SQL

- Locate query construction, ORM preload patterns, loops issuing queries, pagination, sorting, and aggregate queries.
- Prefer actual slow-query logs or measured query timing.
- Use read-only `EXPLAIN` evidence when access and safety are verified.
- Flag suspected N+1 queries or missing indexes as hypotheses until measured.

### 4. Connection Pool

- Inspect `MaxOpenConns`, `MaxIdleConns`, connection lifetime, idle lifetime, timeouts, and database capacity.
- Use `sql.DBStats` where available.
- Relate wait count and wait duration to workload concurrency before recommending pool changes.
- Avoid universal pool-size formulas.

### 5. Concurrency and Locks

- Examine goroutine ownership, bounded concurrency, channels, mutexes, shared state, blocking I/O, and cancellation.
- Use race detection for correctness evidence, not as a performance benchmark.
- Require mutex, block, trace, or workload evidence before recommending synchronization rewrites.

### 6. Cache

- Recommend caching only for measured repeated expensive work.
- Define key, value, scope, TTL, invalidation, consistency, memory bound, and observability.
- Consider request coalescing and duplicate suppression before introducing a persistent cache.
- State correctness and staleness risks explicitly.

### 7. Before/After Comparison

In Phase A, this stage records the baseline and defines the future comparison protocol.
In Phase B, it compares equivalent workloads and reports:

- command and environment;
- sample count and variability;
- time/op, allocs/op, bytes/op, latency percentiles, throughput, SQL count/time, and pool waits as applicable;
- absolute and percentage differences;
- regressions, trade-offs, and confidence limits.

## Report Contract

Write reports to:

`docs/performance/YYYY-MM-DD-<scope>-performance-review.md`

The report must include:

1. executive summary;
2. scope, environment, and safety assumptions;
3. commands run and evidence collected;
4. findings for every workflow stage;
5. prioritized recommendation table;
6. proposed implementation and validation plan;
7. evidence gaps and blocked checks;
8. explicit approval gate.

Each finding includes:

- identifier and severity;
- observed symptom;
- evidence and reproduction command;
- likely cause, clearly labeled as fact or hypothesis;
- proposed change;
- expected metric impact without invented percentages;
- correctness and operational risks;
- validation method.

The report must distinguish:

- measured fact;
- code-based inference;
- unverified hypothesis.

The final line must state that no source code was changed and request explicit approval before Phase B.

## Error and Safety Handling

- Redact credentials, DSNs, tokens, and sensitive query parameters from reports and command output.
- Never print or copy values from `.env`; inspect variable names only when necessary.
- Stop database profiling if the target environment or write safety cannot be verified.
- Avoid benchmarking against production unless the user explicitly authorizes a bounded procedure.
- If measurements are noisy or incomparable, report uncertainty rather than choosing the favorable run.
- Preserve unrelated user changes and avoid generated artifacts outside the approved report and temporary directory.

## Validation Strategy for the Skill

Validate the skill itself using documentation TDD:

1. run baseline agent scenarios without the skill and record where the fixed sequence, evidence standard, or approval gate is violated;
2. create the minimal skill addressing those failures;
3. run equivalent scenarios with the skill;
4. verify the agent covers all seven stages, labels missing evidence, writes a report, and refuses to modify code before approval;
5. add explicit counters for any newly observed shortcut or rationalization;
6. run the official skill validator and inspect `agents/openai.yaml`.

Forward-test at least these cases:

- a Go service with no benchmarks or profiling endpoint;
- a GORM/MySQL service with suspected slow queries and connection waits;
- a user who pressures the agent to apply an “obvious one-line optimization” during Phase A;
- an approved Phase B change requiring a comparable before/after result.

## Success Criteria

The skill is complete when:

- Codex discovers it for relevant Go performance requests;
- it follows all seven stages in order;
- Phase A changes no application source or configuration;
- missing instrumentation becomes a documented proposal, not an unauthorized edit;
- reports are evidence-based, redact secrets, and use the standard template;
- Phase B begins only after explicit user approval;
- performance claims include comparable before/after measurements;
- the skill passes structural validation and forward tests.
