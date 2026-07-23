#  Aero-KV

Aero-KV is a high-performance, concurrent, in-memory key-value store I built from scratch in Go. 

I didn't want to just build another simple wrapper around Go's native map. Instead, I wanted to really get my hands dirty with low-level systems engineering. Aero-KV is my playground for exploring how production-grade databases handle massive concurrency, crash recovery, and memory constraints under heavy loads.
# ⚡ Aero-KV

[![Aero-KV Core Build](https://github.com/GiTaDi-CrEaTe/aero-kv/actions/workflows/go.yml/badge.svg)](https://github.com/GiTaDi-CrEaTe/aero-kv/actions/workflows/go.yml)

Aero-KV is a high-performance, concurrent, in-memory key-value store I built from scratch in Go. 

I didn't want to just build another simple wrapper around Go's native map. Instead, I wanted to explore how production-grade databases handle massive concurrency, crash recovery, and memory constraints under heavy loads.

---

## 🚀 Performance
Aero-KV is built for extreme concurrency. In benchmark testing (Intel i3-6100), the thread-safe engine successfully handled **175,000+ concurrent writes** at a speed of **7,794 ns/op** with zero race conditions or memory leaks.

## ⚡ Quick Start (Embedded Mode)

Aero-KV is designed to be embedded directly into your Go binaries for zero-dependency, sub-microsecond memory access.

### Installation
```bash
go get [github.com/GiTaDi-CrEaTe/aero-kv](https://github.com/GiTaDi-CrEaTe/aero-kv)
```

### Usage
```go
package main

import (
	"fmt"
	"[github.com/GiTaDi-CrEaTe/aero-kv](https://github.com/GiTaDi-CrEaTe/aero-kv)"
)

func main() {
	// Initialize an OOM-safe embedded store with a capacity of 10,000 keys
	db := aerokv.Open(10000)

	// Thread-safe writes
	db.Set("engine", "aero-kv")

	// Lightning-fast reads
	if val, found := db.Get("engine"); found {
		fmt.Printf("Connected to: %s\n", val)
	}
}
```

---

## ⚙️ What's Under the Hood?

* **Concurrency by Design:** I intentionally use a strict `sync.Mutex` rather than an `RWMutex`. Because LRU reads mutate the underlying doubly-linked list (moving accessed nodes to the head), parallel reads would cause race conditions. Prioritizing memory safety and data integrity over theoretical read-speeds is a deliberate architectural choice.
* **Crash Resilience (WAL):** To ensure zero data loss, I built a custom Write-Ahead Log. Every single mutation is appended to a disk-backed log before it ever touches memory. If the server goes down, the data stays safe.
* **Smart Memory Safeguards:** To stop out-of-memory (OOM) crashes during massive data ingestion, I implemented a custom LRU (Least Recently Used) cache from scratch using a doubly-linked list and a hash map.
* **Raw Socket Speed:** I bypassed HTTP entirely and exposed a custom, lightweight TCP server. This strips away web overhead and relies on raw socket performance for client connections.

---

## 🛠️ Codebase Map

```text
aero-kv/
├── core/
│   ├── store.go  ── The central, thread-safe memory engine
│   └── lru.go    ── My custom eviction policy manager
├── wal/
│   └── log.go    ── Disk I/O management and durability sync
└── server/
    └── tcp.go    ── The socket listener and protocol parser
```

---

## 🚦 Project Status
*Current Status: **V1 Stable** 🚀* The core engine, WAL, and embedded APIs are actively tested via GitHub Actions.


---
