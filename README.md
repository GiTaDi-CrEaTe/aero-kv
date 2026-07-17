# 🚀 Aero-KV

Aero-KV is a high-performance, concurrent, in-memory key-value store I built from scratch in Go. 

I didn't want to just build another simple wrapper around Go's native map. Instead, I wanted to really get my hands dirty with low-level systems engineering. Aero-KV is my playground for exploring how production-grade databases handle massive concurrency, crash recovery, and memory constraints under heavy loads.

---

## ⚙️ What's Under the Hood?

*   **Massive Concurrency:** I used `sync.RWMutex` to allow parallel reads to fire at full speed while safely isolating writes. This keeps the state thread-safe without bottlenecking performance.
*   **Crash Resilience (WAL):** To ensure zero data loss, I built a custom Write-Ahead Log. Every single mutation is appended to a disk-backed log before it ever touches memory. If the server goes down, the data stays safe.
*   **Smart Memory Safeguards:** To stop out-of-memory (OOM) crashes during massive data ingestion, I implemented a custom LRU (Least Recently Used) cache from scratch using a doubly-linked list and a hash map.
*   **Raw Socket Speed:** I bypassed HTTP entirely and exposed a custom, lightweight TCP server. This strips away web overhead and relies on raw socket performance for client connections.

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

## 🧠 Interactive Deep Dive: My Engineering Choices

<details>
<summary><b>🔍 Click to expand: Why build a custom LRU from scratch?</b></summary>
<br>
Go's built-in map is fantastic, but it will greedily consume memory until the OS kills the process (OOM panic) under heavy ingestion. By pairing a standard Go map with a doubly-linked list, I can track item recency. When memory hits its limit, the oldest keys are automatically evicted in $O(1)$ time, keeping the footprint highly predictable.
</details>

<details>
<summary><b>⚡ Click to expand: Why raw TCP instead of a REST API?</b></summary>
<br>
HTTP brings a lot of baggage: heavy headers, cookie parsing, and extra text processing. For a key-value store where microseconds matter, raw TCP sockets let me define a minimal byte protocol. The client sends only what is necessary, and the server parses it instantly.
</details>

---

## 🚦 Project Status

*Current Status: **Prototyping Phase** 🛠️*  
I am currently working unde the project.

---
