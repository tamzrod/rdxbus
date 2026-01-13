# RDXBus Project File Tree

This file captures the **authoritative file and directory structure** of the RDXBus repository.
It exists to prevent context loss during reviews, discussions, and refactors.

Always keep this file updated when files are added, removed, or renamed.

---

## File Tree

```text
rdxbus/
├── .gitignore
├── go.mod
│
├── cmd/
│   └── rdxbus/
│       └── main.go
│
└── internal/
    ├── client/
    │   ├── connection.go
    │   ├── parser.go
    │   └── request.go
    │
    ├── config/
    │   └── config.go
    │
    ├── scheduler/
    │   └── rate.go
    │
    ├── stats/
    │   ├── counters.go
    │   ├── histogram.go
    │   └── report.go
    │
    └── worker/
        └── worker.go
```

---

## Notes

- `cmd/rdxbus/main.go` is the CLI entry point only.
- All implementation logic must live under `internal/`.
- Stress testing is a **mode**, not a separate project.
- This tree is the reference for any architectural discussion.

