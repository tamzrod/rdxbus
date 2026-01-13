# RDXBus Project Guidelines

This document defines the **authoritative rules** for the RDXBus project.
Its purpose is to prevent architectural drift, duplicated logic, and uncontrolled growth.

If code conflicts with this document, **the document wins**.

---

## 0. Absolute Rule (Non‑Negotiable)

**Every code snippet, discussion, or proposal MUST explicitly state the file path and package name.**

Required format:

```go
// internal/worker/worker.go
package worker
```

Rules:
- File path must be exact and relative to repo root
- `package` name must match the directory
- Snippets without filenames are considered **invalid context**
- Missing filenames will result in rejected or ignored changes

This rule exists to:
- prevent reference loss
- stop accidental duplication
- avoid ChatGPT or human guesswork
- enforce architectural discipline

---

## 1. Project Identity

RDXBus is a **Modbus client and protocol workbench**.

- Stress testing is a *mode*, not the project identity
- Scanning, polling, and writing are *consumers* of the same engine
- There must be exactly **one implementation** of Modbus read/write execution

---

## 2. Folder Ownership Rules

### `cmd/rdxbus/`
**Purpose:** CLI entry point and orchestration only

Allowed:
- flag parsing
- mode dispatch
- wiring components together

Not allowed:
- Modbus protocol logic
- packet building or parsing
- retry logic
- scheduling logic

If `main.go` grows large, a rule was broken.

---

### `internal/config/`
**Purpose:** Configuration *data*, not I/O

Rules:
- defines structs and validation only
- must not depend on `flag`, `os.Exit`, or CLI concerns
- engine-facing configs must be pure data

CLI parsing belongs in `cmd/rdxbus`, not here.

---

### `internal/client/`
**Purpose:** Modbus protocol implementation

Rules:
- build Modbus request frames
- parse Modbus responses
- manage TCP connections

Must not:
- know about CLI flags
- know about workers, schedulers, or benchmarks
- interpret register meaning or data types

---

### `internal/worker/`
**Purpose:** Execution and concurrency

Rules:
- execute exactly one task per invocation
- manage goroutines and lifecycles
- call client logic

Must not:
- build Modbus packets
- parse protocol responses
- decide scheduling or rates

---

### `internal/scheduler/`
**Purpose:** Time and rate control

Rules:
- pacing, rate limiting, ramp logic
- protocol-agnostic by design

Must not:
- contain Modbus logic
- parse or build packets

---

### `internal/stats/`
**Purpose:** Measurement only

Rules:
- counters, histograms, latency tracking
- error observation

Must not:
- decide pass/fail
- print CLI output
- interpret protocol semantics

---

## 3. Engine Rules

- There must be **one read engine**
- There must be **one write engine**
- Benchmarking reuses the same engines under load
- No mode may reimplement protocol logic

If similar logic appears twice, it is a bug.

---

## 4. Configuration Rules

- Configuration is data, not behavior
- Parsing, defaults, and I/O belong to the CLI layer
- Engine configs must be reusable by:
  - CLI
  - tests
  - future UI / REST

---

## 5. Non-Goals (By Design)

RDXBus will not:
- scale values
- guess data types
- hide Modbus errors
- act as a PLC or RTU
- embed a GUI
- implement SCADA logic

Those belong upstream.

---

## 6. Reference Files

The following files are authoritative references:

- `TREE.md` — file and directory structure
- `GUIDELINES.md` — architectural rules (this file)

Any discussion, refactor, or feature proposal must reference them.

---

## 7. Change Discipline

When adding new code:
1. Include exact filename and package declaration
2. Update `TREE.md` if structure changes
3. Verify rules in this document are not violated
4. Prefer reuse over extension
5. Avoid premature abstraction

---

## Final Rule

> Missing filenames or missing package declarations
> invalidate architectural discussion.

