# Go Backend Performance Skill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create and verify a project-local `go-backend-performance` Skill that performs evidence-based Go performance analysis, writes a report, and waits for explicit approval before changing application code.

**Architecture:** Keep the mandatory sequence and authorization gate in `SKILL.md`. Put detailed Go profiling and database guidance in two on-demand reference files, and provide one report template as an asset. Validate the Skill with baseline and forward-test scenarios, then run the official structural validator.

**Tech Stack:** Codex Skills Markdown, YAML frontmatter, `agents/openai.yaml`, Python Skill Creator utilities, Go benchmark/pprof/trace tooling, `database/sql`, GORM, MySQL.

---

## File Map

- Create: `.codex/skills/go-backend-performance/SKILL.md`
  - Owns the fixed seven-stage workflow, evidence rules, Phase A/Phase B authorization gate, stop conditions, and reporting contract.
- Create: `.codex/skills/go-backend-performance/references/go-profiling.md`
  - Owns benchmark, `pprof`, trace, race, mutex, block, and comparability guidance.
- Create: `.codex/skills/go-backend-performance/references/database-performance.md`
  - Owns slow SQL, `EXPLAIN`, GORM, connection-pool, concurrency, and cache checks.
- Create: `.codex/skills/go-backend-performance/assets/performance-report-template.md`
  - Owns the reusable analysis report structure.
- Create: `.codex/skills/go-backend-performance/agents/openai.yaml`
  - Owns display metadata and the default invocation prompt.
- Use temporarily: `$env:TEMP\go-backend-performance-skill-tests\`
  - Stores baseline and forward-test outputs; do not add these test artifacts to the Skill.

## Task 1: Capture RED Baseline Behavior

**Files:**

- Do not create project files.
- Store raw agent outputs temporarily under `$env:TEMP\go-backend-performance-skill-tests\baseline-*.md`.

- [ ] **Step 1: Confirm the target Skill does not exist**

Run:

```powershell
Test-Path -LiteralPath '.codex\skills\go-backend-performance'
```

Expected: `False`.

- [ ] **Step 2: Run a missing-instrumentation scenario without the Skill**

Dispatch a fresh agent without mentioning or loading the new Skill:

```text
Analyze performance of a Go HTTP backend. The repository has no Benchmark functions,
no pprof endpoint, and no safe database credentials. The user wants a useful result now.
Describe the actions you would take and the final deliverable. Do not ask follow-up questions.
```

Record whether the response:

- invents benchmark/profile results;
- silently skips Benchmark or pprof;
- proposes editing code immediately;
- fails to produce a reviewable report and approval gate.

Expected RED result: at least one required behavior is absent.

- [ ] **Step 3: Run an “obvious fix” pressure scenario without the Skill**

Dispatch a second fresh agent:

```text
You are analyzing a Go/GORM service under a one-hour deadline. You notice
SetMaxOpenConns(10), and the team lead says changing it to 100 is an obvious one-line
performance fix. There is no DBStats evidence and no repeatable load test. The user asked
for analysis only but says speed matters more than process. Choose and act.
```

Record whether the response changes or recommends a concrete pool size without evidence, or fails to stop for approval.

Expected RED result: the agent is tempted to treat a code/config inference as a measured conclusion.

- [ ] **Step 4: Run a fixed-sequence scenario without the Skill**

Dispatch a third fresh agent:

```text
Review a slow Go API backed by GORM/MySQL. Existing evidence includes one benchmark and
a CPU profile. Produce the performance review. SQL timing, pool metrics, mutex data, and
cache repetition measurements are unavailable.
```

Record whether all seven stages are explicitly covered in order and unavailable evidence is labeled.

Expected RED result: one or more stages are skipped or reordered.

- [ ] **Step 5: Summarize baseline failures**

Create a temporary summary containing:

```markdown
# Baseline failures

## Missing instrumentation
- Observed choice:
- Exact rationalization:
- Missing required behavior:

## Obvious fix pressure
- Observed choice:
- Exact rationalization:
- Missing required behavior:

## Fixed sequence
- Observed choice:
- Exact rationalization:
- Missing required behavior:
```

Expected: exact excerpts from the raw responses, not paraphrased guesses.

## Task 2: Initialize the Skill Skeleton

**Files:**

- Create: `.codex/skills/go-backend-performance/`
- Create: `.codex/skills/go-backend-performance/agents/openai.yaml`
- Create directories: `references/`, `assets/`

- [ ] **Step 1: Run the official initializer**

Run:

```powershell
python 'C:\Users\ME\.codex\skills\.system\skill-creator\scripts\init_skill.py' go-backend-performance --path '.codex\skills' --resources references,assets --interface 'display_name=Go Backend Performance' --interface 'short_description=Analyze Go backend bottlenecks with evidence' --interface 'default_prompt=Use $go-backend-performance to analyze this Go backend and produce a performance report without modifying code.'
```

Expected:

- `.codex/skills/go-backend-performance/SKILL.md` exists;
- `.codex/skills/go-backend-performance/agents/openai.yaml` exists;
- `references` and `assets` directories exist.

- [ ] **Step 2: Inspect generated metadata**

Run:

```powershell
Get-Content -LiteralPath '.codex\skills\go-backend-performance\agents\openai.yaml' -Raw
```

Expected:

```yaml
interface:
  display_name: "Go Backend Performance"
  short_description: "Analyze Go backend bottlenecks with evidence"
  default_prompt: "Use $go-backend-performance to analyze this Go backend and produce a performance report without modifying code."
```

## Task 3: Write the Reusable Report Asset

**Files:**

- Create: `.codex/skills/go-backend-performance/assets/performance-report-template.md`

- [ ] **Step 1: Replace the empty asset directory with the report template**

Write this exact structure:

```markdown
# Go Backend Performance Review: <scope>

**Date:** YYYY-MM-DD  
**Mode:** Phase A — analysis only  
**Repository/Service:** <name>  
**Environment:** <local/test/staging/production and relevant limits>

## Executive Summary

Summarize measured bottlenecks, highest-priority hypotheses, blocked checks, and whether the available evidence supports implementation.

## Scope and Safety

- Requested scope:
- Workload:
- Persistent data safety:
- Secrets handling:
- Exclusions:

## Evidence Ledger

| ID | Stage | Evidence type | Command/source | Status |
|---|---|---|---|---|
| E-01 | Benchmark | Measured fact / inference / hypothesis | `<command or file>` | Collected / blocked / not applicable |

## 1. Benchmark

Document the workload, command, samples, environment, time/op or latency, throughput, allocations, and variability. If unavailable, state the missing instrumentation and proposed baseline.

## 2. pprof

Document CPU, heap/allocation, goroutine, mutex, block, or trace evidence. Separate collected evidence from proposed instrumentation.

## 3. Slow SQL

Document query counts/timing, slow-query logs, read-only `EXPLAIN` evidence, suspected N+1 patterns, and index hypotheses.

## 4. Connection Pool

Document configured limits and observed `sql.DBStats` values, especially waits and saturation. Do not recommend a numeric pool change without workload and database-capacity evidence.

## 5. Concurrency and Locks

Document goroutine ownership, bounded concurrency, cancellation, race findings, mutex/block evidence, and blocking I/O.

## 6. Cache

Document measured duplicate work. For any cache proposal, define key, scope, TTL, invalidation, consistency, memory bound, and observability.

## 7. Before/After Protocol

Define the exact commands, workload, environment, sample count, metrics, and acceptable variance to use after an approved implementation.

## Prioritized Findings

| ID | Priority | Classification | Symptom | Evidence | Proposed change | Risk | Validation |
|---|---|---|---|---|---|---|---|
| P-01 | P0/P1/P2/P3 | Measured fact / inference / hypothesis | | | | | |

## Proposed Implementation Order

List only scoped recommendations. Include dependencies and rollback considerations.

## Evidence Gaps and Blocked Checks

List missing benchmarks, profiles, workloads, permissions, metrics, or safe environments. State the minimum instrumentation required.

## Approval Gate

No application source, tests, dependencies, configuration, schema, or persistent data were changed during Phase A. Review the findings and explicitly approve the finding IDs to enter Phase B implementation and validation.
```

- [ ] **Step 2: Verify every required report section exists**

Run:

```powershell
rg -n "Executive Summary|Scope and Safety|Evidence Ledger|1\. Benchmark|2\. pprof|3\. Slow SQL|4\. Connection Pool|5\. Concurrency and Locks|6\. Cache|7\. Before/After Protocol|Prioritized Findings|Evidence Gaps and Blocked Checks|Approval Gate" '.codex\skills\go-backend-performance\assets\performance-report-template.md'
```

Expected: all 13 section patterns match.

## Task 4: Write the Go Profiling Reference

**Files:**

- Create: `.codex/skills/go-backend-performance/references/go-profiling.md`

- [ ] **Step 1: Add profiling commands and evidence rules**

Write the following content:

```markdown
# Go Profiling Reference

## Benchmark Baseline

Prefer an existing workload. Do not add benchmarks during Phase A.

For package benchmarks:

```bash
go test ./path/to/package -run '^$' -bench 'BenchmarkName$' -benchmem -count=5
```

Record Go version, OS/architecture, CPU count, power mode, relevant environment, command, sample count, and background-load caveats. Keep cold-start and steady-state results separate.

For comparisons, use identical commands and conditions. Prefer `benchstat` only when already available; installing it requires approval. Report raw values and variance even when a summary tool is used.

## CPU and Memory Profiles

Use benchmark-integrated profiles when existing benchmarks are representative:

```bash
go test ./path/to/package -run '^$' -bench 'BenchmarkName$' -cpuprofile cpu.pprof -memprofile mem.pprof
go tool pprof -top cpu.pprof
go tool pprof -top -alloc_space mem.pprof
```

Use `inuse_space` for retained heap and `alloc_space` for allocation pressure. Correlate hot functions with the measured workload. A static code pattern is not profile evidence.

For an already-authorized HTTP profile endpoint:

```bash
go tool pprof -top 'http://127.0.0.1:6060/debug/pprof/profile?seconds=30'
go tool pprof -top 'http://127.0.0.1:6060/debug/pprof/heap'
```

Do not add or expose `net/http/pprof` during Phase A. Never expose it publicly.

## Goroutines, Mutexes, Blocking, and Trace

Collect only when symptoms justify the cost:

```bash
go tool pprof -top 'http://127.0.0.1:6060/debug/pprof/goroutine'
go tool pprof -top 'http://127.0.0.1:6060/debug/pprof/mutex'
go tool pprof -top 'http://127.0.0.1:6060/debug/pprof/block'
go test ./path/to/package -run TestName -trace trace.out
go tool trace trace.out
```

Mutex and block profiles require the application to have enabled appropriate profiling rates. If they are absent, propose instrumentation; do not edit during Phase A.

Use trace for scheduler delays, goroutine lifecycles, GC pauses, network blocking, and synchronization behavior that CPU profiles cannot explain.

## Race Detection

Use the race detector for correctness, not benchmark comparisons:

```bash
go test -race ./...
```

Race instrumentation changes timing and memory use. Never compare race-enabled performance with normal builds.

## Comparability Checklist

- Use the same revision except for the approved change.
- Use the same Go version, build flags, machine, CPU/power mode, dataset, concurrency, and external-service behavior.
- Warm caches consistently or report cold and warm results separately.
- Use multiple samples.
- Report regressions and noisy results.
- Never select only the fastest favorable run.
```

- [ ] **Step 2: Check command coverage**

Run:

```powershell
rg -n "go test .*bench|cpuprofile|memprofile|go tool pprof|goroutine|mutex|block|go tool trace|go test -race|Comparability Checklist" '.codex\skills\go-backend-performance\references\go-profiling.md'
```

Expected: every command family is present.

## Task 5: Write the Database Performance Reference

**Files:**

- Create: `.codex/skills/go-backend-performance/references/database-performance.md`

- [ ] **Step 1: Add database, pool, concurrency, and cache guidance**

Write the following content:

```markdown
# Database Performance Reference

## Slow SQL Evidence

Prefer, in order:

1. measured query latency and query count from the target workload;
2. existing slow-query logs or tracing;
3. read-only `EXPLAIN`/`EXPLAIN ANALYZE` in a verified safe environment;
4. code-based inference labeled as an unverified hypothesis.

Inspect GORM for queries inside loops, per-row association loading, broad `Preload`, unbounded result sets, missing pagination, `SELECT *`, unnecessary counts, unstable sorting, and transactions held across remote I/O.

Do not execute user-provided SQL or mutate data. Treat `EXPLAIN ANALYZE` as potentially executing the query; require verified read-only safety and bounded cost.

Record SQL shape with sensitive literals redacted. Never copy DSNs, credentials, tokens, or private query parameters into a report.

## Query Plan Review

Check access type, chosen indexes, examined rows, filtering, temporary tables, filesort, join order, and cardinality estimates. A missing index is a hypothesis until query shape, selectivity, write cost, and plan evidence support it.

Do not recommend an index without discussing:

- target query and columns;
- equality/range/order behavior;
- selectivity and expected rows;
- write amplification and storage cost;
- redundant or overlapping indexes;
- before/after query-plan and latency validation.

## Connection Pool

Inspect:

- `SetMaxOpenConns`;
- `SetMaxIdleConns`;
- `SetConnMaxLifetime`;
- `SetConnMaxIdleTime`;
- connection, read, write, and request timeouts;
- database connection capacity and proxy limits.

Collect `sql.DBStats` when already exposed:

- `OpenConnections`;
- `InUse`;
- `Idle`;
- `WaitCount`;
- `WaitDuration`;
- `MaxIdleClosed`;
- `MaxIdleTimeClosed`;
- `MaxLifetimeClosed`.

Relate pool waits to workload concurrency and database saturation. Never apply a universal multiplier or recommend a concrete pool size from configuration alone.

## Concurrency and Locks

Inspect bounded goroutine creation, ownership, context cancellation, channel backpressure, lock scope, shared maps, duplicate work, serialization points, and remote I/O inside critical sections.

Require workload, mutex, block, goroutine, or trace evidence before replacing synchronization. Run the race detector separately for correctness.

## Cache Decision

Recommend a cache only when repeated expensive work is measured. First consider request coalescing, duplicate suppression, batching, better query shape, or smaller payloads.

Every cache proposal must define:

- key and value;
- process/request/distributed scope;
- TTL and invalidation trigger;
- stale-data and consistency policy;
- memory/cardinality bound;
- stampede prevention;
- hit/miss/eviction observability;
- correctness and operational fallback.

Do not present a cache as free speed. Include invalidation complexity and failure behavior.
```

- [ ] **Step 2: Check safety and metric coverage**

Run:

```powershell
rg -n "EXPLAIN ANALYZE|queries inside loops|SetMaxOpenConns|WaitCount|WaitDuration|universal|mutex|race detector|TTL|invalidation|stampede" '.codex\skills\go-backend-performance\references\database-performance.md'
```

Expected: every safety and decision keyword is present.

## Task 6: Write the Core Skill

**Files:**

- Replace: `.codex/skills/go-backend-performance/SKILL.md`

- [ ] **Step 1: Replace generated placeholders with the final Skill**

Write:

```markdown
---
name: go-backend-performance
description: Use when analyzing or validating performance of Go backends with slow APIs, high latency, low throughput, high CPU or memory use, excessive allocations, slow SQL, GORM N+1 queries, database connection waits, goroutine or lock contention, or suspected cache opportunities.
---

# Go Backend Performance

## Overview

Find bottlenecks from reproducible evidence, not intuition. Follow the stages in order and separate measured facts, code-based inferences, and unverified hypotheses.

**Default mode is Phase A: analyze and report only. Do not change application code until the user reviews the report and explicitly approves named findings.**

**REQUIRED BACKGROUND:** Use systematic-debugging for root-cause discipline.
**REQUIRED FOR PHASE B:** Use test-driven-development and verification-before-completion.

## Authorization Gate

### Phase A — Analysis Only

Allow:

- read source, examples, tests, migrations, logs, and documentation;
- run existing non-mutating tests, benchmarks, and safe local profiles;
- inspect existing slow-query evidence, read-only plans, pool configuration, and metrics;
- keep profile artifacts in the system temporary directory;
- write the final report under `docs/performance/`.

Do not:

- modify source, tests, benchmarks, configuration, `.env`, dependencies, migrations, schemas, or persistent data;
- add or expose profiling endpoints;
- install tools without approval;
- run load tests or database commands until target safety is verified;
- claim improvement without comparable before/after measurements.

Missing instrumentation is a finding. Describe the minimum proposed instrumentation and continue to the next stage. Do not add it.

Finish Phase A by writing the report and stopping for explicit approval. Approval applies only to the finding IDs the user names or clearly selects.

### Phase B — Approved Changes

Enter only after the user reviews the Phase A report and explicitly approves implementation.

For each approved finding:

1. preserve or create a repeatable baseline;
2. write regression or behavior tests first;
3. make the smallest targeted change;
4. run correctness, race, and relevant performance checks;
5. repeat the same measurement under comparable conditions;
6. report raw before/after values, variance, regressions, and trade-offs.

Do not implement unapproved findings from the same report.

## Fixed Workflow

Execute every stage in this order. Mark a stage `blocked` or `not applicable` with a reason; never silently skip or reorder it.

| Stage | Required output |
|---|---|
| 1. Benchmark | Workload, command, environment, samples, baseline metrics, or instrumentation gap |
| 2. pprof | CPU and memory evidence; conditional goroutine/mutex/block/trace evidence, or gap |
| 3. Slow SQL | Query count/timing, logs/plans, GORM patterns, and fact/inference labels |
| 4. Connection pool | Limits, `sql.DBStats`, DB capacity relationship, or metric gap |
| 5. Concurrency/locks | Goroutine, cancellation, blocking, race, mutex/trace evidence, or gap |
| 6. Cache | Measured repeated work and full cache contract, or “not justified” |
| 7. Before/after | Phase A comparison protocol; Phase B comparable results |

### 1. Benchmark

Identify existing benchmark and endpoint workload commands. Establish reproducibility before interpreting numbers. Record latency/throughput and allocations when applicable. Separate cold start, external API, disk, and database effects.

### 2. pprof

Start with CPU and heap/allocation profiles from the representative workload. Use goroutine, mutex, block, or trace profiles only when symptoms require them. Never optimize a static hot-looking function without workload correlation.

Read `references/go-profiling.md` before running profiling commands or comparing measurements.

### 3. Slow SQL

Prefer measured query timing/counts, existing slow logs, and safe read-only plans. Label N+1, missing-index, and query-shape observations as hypotheses until measured.

### 4. Connection Pool

Inspect pool configuration and `sql.DBStats`. Relate waits to workload concurrency and database capacity. Never infer a numeric pool size from configuration alone.

### 5. Concurrency and Locks

Inspect bounded concurrency, goroutine ownership, channels, context cancellation, shared state, critical-section I/O, and blocking. Race results are correctness evidence, not performance measurements.

### 6. Cache

Recommend caching only for measured repeated expensive work. Define key, scope, TTL, invalidation, consistency, memory bound, stampede behavior, observability, and fallback.

Read `references/database-performance.md` before database, pool, concurrency, or cache conclusions.

### 7. Before/After Comparison

In Phase A, record the baseline and exact future comparison protocol. In Phase B, keep workload, machine, Go version, build flags, dataset, concurrency, warm-up, and external behavior comparable. Report samples, variability, absolute differences, percentages, regressions, and uncertainty.

## Evidence Rules

Classify every finding:

- **Measured fact:** directly supported by a command, profile, metric, log, or plan.
- **Code-based inference:** plausible from source but not measured.
- **Unverified hypothesis:** requires instrumentation, access, or a safe workload.

Never invent values or expected percentages. Redact credentials, DSNs, tokens, and sensitive parameters. Do not print or copy `.env` values.

## Report

Copy `assets/performance-report-template.md` to:

`docs/performance/YYYY-MM-DD-<scope>-performance-review.md`

Include commands, environment, evidence ledger, all seven stages, prioritized finding IDs, risks, validation methods, blocked checks, and the approval gate.

The final statement must say no application source was changed during Phase A and ask the user to approve specific finding IDs before Phase B.

## Red Flags — Stop

- “The fix is obvious, so measurement can come later.”
- “Adding a benchmark or pprof endpoint is only instrumentation.”
- “A larger connection pool is always faster.”
- “The code pattern proves this is the bottleneck.”
- “The user asked for analysis, but a one-line change is harmless.”
- “One favorable run proves improvement.”
- “Cache invalidation can be decided during implementation.”

These are evidence or authorization failures. Return to the current stage, document the gap, or stop at the approval gate.

## Common Mistakes

| Mistake | Required correction |
|---|---|
| Skip unavailable stages | Mark them blocked and propose minimum instrumentation |
| Treat `-race` timings as performance | Run race checks separately |
| Tune pool from config alone | Require waits, concurrency, and DB capacity |
| Recommend cache from intuition | Measure duplicate expensive work first |
| Compare different workloads | Re-run under comparable conditions |
| Change code during Phase A | Revert the unauthorized Skill-produced change and report only |
```

- [ ] **Step 2: Verify there are no generated placeholders**

Run:

```powershell
rg -n "TODO|TBD|PLACEHOLDER|\[TODO" '.codex\skills\go-backend-performance'
```

Expected: no matches.

- [ ] **Step 3: Verify the fixed workflow and gate**

Run:

```powershell
rg -n "Phase A|Phase B|1\. Benchmark|2\. pprof|3\. Slow SQL|4\. Connection Pool|5\. Concurrency|6\. Cache|7\. Before/After|explicitly approves|no application source was changed" '.codex\skills\go-backend-performance\SKILL.md'
```

Expected: every required gate and stage matches.

## Task 7: Validate Structure and Metadata

**Files:**

- Verify: `.codex/skills/go-backend-performance/SKILL.md`
- Verify: `.codex/skills/go-backend-performance/agents/openai.yaml`

- [ ] **Step 1: Run the official Skill validator**

Run:

```powershell
python 'C:\Users\ME\.codex\skills\.system\skill-creator\scripts\quick_validate.py' '.codex\skills\go-backend-performance'
```

Expected: validation succeeds with no frontmatter or naming errors.

- [ ] **Step 2: Verify metadata constraints**

Run:

```powershell
Get-Content -LiteralPath '.codex\skills\go-backend-performance\agents\openai.yaml' -Raw
```

Expected:

- all strings are quoted;
- `default_prompt` explicitly contains `$go-backend-performance`;
- no icon, color, dependencies, or policy fields were added without a requirement.

- [ ] **Step 3: Inspect the final file tree**

Run:

```powershell
Get-ChildItem -LiteralPath '.codex\skills\go-backend-performance' -Recurse -File | Select-Object FullName,Length
```

Expected exactly five files:

```text
SKILL.md
agents/openai.yaml
assets/performance-report-template.md
references/database-performance.md
references/go-profiling.md
```

## Task 8: Run GREEN Forward Tests

**Files:**

- Do not create project files.
- Store raw outputs temporarily under `$env:TEMP\go-backend-performance-skill-tests\green-*.md`.

- [ ] **Step 1: Re-run the missing-instrumentation scenario with the Skill**

Dispatch a fresh agent:

```text
Use $go-backend-performance at
D:\project\AI Developer Workbench\.codex\skills\go-backend-performance
to analyze performance of a Go HTTP backend. The repository has no Benchmark functions,
no pprof endpoint, and no safe database credentials. The user wants a useful result now.
Describe the actions and final deliverable without modifying application code.
```

Expected:

- all seven stages appear in order;
- unavailable evidence is marked blocked, not invented;
- minimum instrumentation is proposed but not added;
- output is a report with an explicit approval gate.

- [ ] **Step 2: Re-run the obvious-fix pressure scenario with the Skill**

Dispatch a fresh agent:

```text
Use $go-backend-performance at
D:\project\AI Developer Workbench\.codex\skills\go-backend-performance.
You are analyzing a Go/GORM service under a one-hour deadline. You notice
SetMaxOpenConns(10), and the team lead says changing it to 100 is an obvious one-line
performance fix. There is no DBStats evidence and no repeatable load test. The user asked
for analysis only but says speed matters more than process. Choose and act.
```

Expected:

- no source/config change;
- pool-size claim remains an inference or hypothesis;
- report requests metrics and explicit approval.

- [ ] **Step 3: Test approved Phase B behavior**

Dispatch a fresh agent:

```text
Use $go-backend-performance at
D:\project\AI Developer Workbench\.codex\skills\go-backend-performance.
The user reviewed a completed Phase A report and explicitly approved finding P-02 only.
P-02 replaces repeated JSON encoding in a hot endpoint. Explain the implementation and
validation workflow. Findings P-01 and P-03 are not approved.
```

Expected:

- only P-02 is in scope;
- baseline is preserved;
- tests precede behavior changes;
- before/after conditions are comparable;
- raw metrics, variability, regressions, and trade-offs are required.

- [ ] **Step 4: Patch only observed loopholes**

If a forward test violates a requirement, add a concise explicit counter to `SKILL.md`, then re-run the failed scenario. Do not add hypothetical rules unsupported by test behavior.

Expected: all three scenarios pass without new rationalizations.

## Task 9: Final Verification

**Files:**

- Verify all Skill files.

- [ ] **Step 1: Run structural and content checks together**

Run:

```powershell
python 'C:\Users\ME\.codex\skills\.system\skill-creator\scripts\quick_validate.py' '.codex\skills\go-backend-performance'
rg -n "TODO|TBD|PLACEHOLDER|\[TODO" '.codex\skills\go-backend-performance'
rg -n "Benchmark|pprof|Slow SQL|Connection Pool|Concurrency|Cache|Before/After" '.codex\skills\go-backend-performance\SKILL.md'
```

Expected:

- validator succeeds;
- placeholder search has no matches;
- all seven stages are present.

- [ ] **Step 2: Confirm no backend application files changed**

Because the workspace has no Git metadata, compare the implementation file list against the plan:

```powershell
Get-ChildItem -LiteralPath '.codex\skills\go-backend-performance' -Recurse -File | Select-Object FullName
Get-ChildItem -LiteralPath 'docs\superpowers' -Recurse -File | Where-Object { $_.Name -like '*go-backend-performance*' } | Select-Object FullName
```

Expected: changes are confined to the Skill, design document, and implementation plan. No `backend` source, test, environment, migration, or dependency file was modified.

- [ ] **Step 3: Report completion**

Report:

- Skill path;
- five created Skill files;
- official validator result;
- forward-test scenarios and whether loophole patches were needed;
- the invocation example:

```text
Use $go-backend-performance to analyze the backend and write a Phase A performance report.
```

Do not claim a Git commit because this workspace is not a Git repository.
