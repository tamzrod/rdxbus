# RDXBus: AI Copilot Instructions

**RDXBus** is a Modbus TCP client stress testing and protocol workbench. This document guides AI agents to make productive contributions while respecting architectural constraints.

## 1. Non-Negotiable Rule

Every code snippet, discussion, or proposal **MUST explicitly state the file path and package name**:

```go
// internal/worker/worker.go
package worker
```

Snippets without filenames are considered **invalid context** and will be rejected. This prevents duplication and architectural drift.

---

## 2. Architecture Overview

### Components
- **`internal/client/`** — Modbus TCP protocol engine (frames, parsing, connections)
- **`internal/config/`** — Configuration data and validation only
- **`internal/scheduler/`** — Rate limiting and pacing (protocol-agnostic)
- **`internal/stats/`** — Metrics collection (counters, histograms, latency)
- **`internal/worker/`** — Concurrent task execution (goroutine lifecycle)
- **`cmd/rdxbus/`** — CLI entry point and mode orchestration only

### Critical Principle: One Engine
There is **exactly one read execution path**. All consumers (stress testing, polling, scanning) reuse the same Modbus protocol logic. No reimplementation is allowed—if similar logic appears twice, it is a bug.

---

## 3. Folder Ownership Rules

### `cmd/rdxbus/main.go`
**Allowed:** Flag parsing, signal handling, mode dispatch, component wiring  
**Forbidden:** Protocol logic, packet building/parsing, scheduling, retry logic  
**Pattern:** Keep `main()` focused on orchestration. If it grows large, a rule was broken.

### `internal/config/config.go`
**Allowed:** Struct definitions, flag parsing, validation  
**Forbidden:** CLI state management, engine behavior, file I/O for business logic  
**Pattern:** Configuration is data, not behavior. CLI parsing belongs here; engine behavior does not.

### `internal/client/{connection,request,parser}.go`
**Allowed:** Build Modbus request frames (FC 1–4), parse responses, TCP connection lifecycle  
**Forbidden:** Worker awareness, scheduling logic, register interpretation  
**Pattern:** Zero knowledge of CLI flags or worker semantics. Protocol-only.

**Key classes:**
- `Connection`: TCP lifecycle, deadline management, Nagle disabling
- `Request`: Builds MBAP + PDU frames, manages Transaction IDs
- `ResponseParser`: Validates responses, detects Modbus exceptions

### `internal/worker/worker.go`
**Allowed:** Goroutine lifecycle, task execution loop, single request-response cycle  
**Forbidden:** Packet building, protocol parsing, rate control decisions  
**Pattern:** One invocation = one request cycle. Defers cleanup in `Run()`.

### `internal/scheduler/rate.go`
**Allowed:** Token-bucket rate limiting, ramp scheduling, pacing  
**Forbidden:** Modbus logic, request building, worker management  
**Pattern:** Protocol-agnostic. Works with any client.

### `internal/stats/{counters,histogram,report}.go`
**Allowed:** Atomic counters, latency histograms, result aggregation  
**Forbidden:** Pass/fail decisions, protocol interpretation, CLI formatting  
**Pattern:** Pure measurement. Main uses atomic ops for goroutine-safe updates.

---

## 4. Data Flow & Integration Points

```
main.go
  ├─ config.Parse() → Config struct
  ├─ create Scheduler (rate limiting)
  ├─ spawn Workers (for i < cfg.Workers)
  │   └─ Worker.Run() reads from scheduler token channel
  │       ├─ client.Request.BuildReadRequest() → frame
  │       ├─ Connection.Write(frame)
  │       ├─ Connection.ReadFull(header)
  │       └─ ResponseParser.Parse() → Result
  └─ collect Results (chan worker.Result)
      └─ Counters.IncOK() / IncExceptions() / IncOtherErrs()
```

**Key flow rule:** Worker pulls tokens from scheduler; scheduler knows nothing of workers.

---

## 5. Common Patterns

### Error Handling
Modbus protocol errors are distinct:
```go
// Identify Modbus exceptions from other I/O errors
if me, ok := client.IsModbusException(r.Err); ok {
    counters.IncExceptions()  // protocol-level failure
} else if r.Err != nil {
    counters.IncOtherErrs()   // network/timeout/parse
}
```

### Configuration Validation
- All validation is in `config.validate()` 
- Invalid configs call `os.Exit(1)` immediately
- Ramp rates are parsed comma-separated: `"100,500,1000"` → `[]int{100, 500, 1000}`

### Concurrency
- Workers are goroutines spawned per invocation; no worker pool
- All stats updates use atomic operations (`sync/atomic`)
- Stop signal is a closed channel; all workers select on it

### Rate Limiting
- Unlimited mode (rate=0): scheduler channel is pre-closed, always ready
- Limited mode (rate>0): token-bucket fills at `interval = 1s / rate`
- Dropped tokens: if workers are slower than rate, tokens are discarded

---

## 6. Configuration Flags

Core flags (from `internal/config/config.go`):
- `-target` (default: `127.0.0.1:502`) — Modbus TCP endpoint
- `-workers` (default: `10`) — concurrent worker count
- `-rate` (default: `0`, unlimited) — requests per second
- `-duration` (default: `10s`) — test duration
- `-ramp` — comma-separated rates for multi-step load test
- `-step-duration` (default: `5s`) — duration per ramp step
- `-unit`, `-fc`, `-address`, `-quantity` — Modbus parameters
- `-timeout` (default: `100ms`) — socket read/write deadline
- `-strict` (default: `false`) — strict MBAP framing validation
- `-quiet` (default: `false`) — minimal output

---

## 7. Key Files Reference

| File | Responsibility |
|------|-----------------|
| [cmd/rdxbus/main.go](../cmd/rdxbus/main.go) | CLI orchestration, mode dispatch |
| [internal/config/config.go](../internal/config/config.go) | Configuration parsing & validation |
| [internal/client/connection.go](../internal/client/connection.go) | TCP lifecycle, socket options |
| [internal/client/request.go](../internal/client/request.go) | MBAP + PDU frame building |
| [internal/client/parser.go](../internal/client/parser.go) | Response parsing, exception handling |
| [internal/worker/worker.go](../internal/worker/worker.go) | Task execution, goroutine lifecycle |
| [internal/scheduler/rate.go](../internal/scheduler/rate.go) | Token-bucket rate control |
| [internal/stats/counters.go](../internal/stats/counters.go) | Atomic result counters |
| [docs/GUIDELINES.md](../docs/GUIDELINES.md) | Architectural rules (authoritative) |
| [docs/TREE.md](../docs/TREE.md) | File structure reference |

---

## 8. When Changing Code

1. **Identify the owner.** Which module should this logic live in?
2. **Check the rules.** Verify against GUIDELINES.md that the change respects ownership.
3. **Use exact filenames.** Include `// path/to/file.go` and `package` name.
4. **Avoid reimplementation.** If similar code exists, reuse or refactor—don't duplicate.
5. **Update docs.** If structure changes, update `TREE.md`.

---

## 9. Non-Goals

RDXBus will **not**:
- Scale or interpret register values (upstream concern)
- Guess data types or endianness
- Hide Modbus errors
- Act as a PLC, RTU, or slave
- Embed a GUI
- Implement SCADA logic

These responsibilities belong upstream of the stress test engine.

---

## 10. Build & Run

```bash
cd rdxbus
go build -o rdxbus ./cmd/rdxbus
./rdxbus -target 127.0.0.1:502 -workers 50 -rate 1000 -duration 10s
```

For ramp testing:
```bash
./rdxbus -target 127.0.0.1:502 -ramp "100,500,1000" -step-duration 5s
```

---

**Reference:** See `docs/GUIDELINES.md` and `docs/TREE.md` for the authoritative architecture rules.
