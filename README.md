# redis-mini

A Redis-compatible in-memory key-value store written in Go.

redis-mini implements the [RESP (REdis Serialization Protocol)](https://redis.io/docs/reference/protocol-spec/) wire format and can be used with `redis-cli` as a drop-in replacement for a real Redis server for basic operations.

## Architecture

```
redis-cli (or any RESP client)
        │
        │  TCP :8080
        ▼
┌─────────────────────────────────┐
│           main.go               │
│  AOF replay on startup          │
│  net.Listen → Accept loop       │
│  goroutine per connection       │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│       HandleConnection          │
│  bufio.Reader (one per conn)    │
│  persistent command loop        │
│  RESP response writes           │
└──────┬──────────────┬───────────┘
       │              │
       ▼              ▼
┌─────────────┐  ┌──────────────────┐
│ HandleResp  │  │ HandleAOFWrite   │
│ ParseArray  │  │ Appends RESP     │
│ ParseBulk   │  │ buffer to        │
│ String      │  │ database.aof     │
└──────┬──────┘  └──────────────────┘
       │
       ▼
┌─────────────────────────────────┐
│         data.Store              │
│  map[string]interface{}         │
│  sync.RWMutex                   │
└─────────────────────────────────┘
```

## Features

- RESP protocol parser supporting array and inline command formats
- Persistent TCP connections — multiple commands per session
- Goroutine-per-connection concurrency with `sync.RWMutex`-guarded shared state
- AOF (Append-Only File) persistence — writes are logged on every `SET` and replayed from `database.aof` on startup
- Compatible with `redis-cli`

## Supported Commands

| Command | Description |
|---------|-------------|
| `PING` | Returns `PONG` |
| `SET key value` | Stores a key-value pair |
| `GET key` | Returns the value for a key, or null if not found |
| `COMMAND` | Returns an empty array (satisfies `redis-cli` startup handshake) |

## In Progress

- `DEL key` — key deletion (store support exists, command not yet wired)
- `EXISTS key` — key existence check
- TTL / key expiry (`SET key value EX seconds`)
- AOF missing-file handling on fresh startup

## Getting Started

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/pierredyl/redis-mini.git
cd redis-mini
go run ./cmd/main.go
```

The server listens on `:8080`. Connect with `redis-cli`:

```bash
redis-cli -p 8080
```

```
127.0.0.1:8080> SET foo bar
OK
127.0.0.1:8080> GET foo
"bar"
127.0.0.1:8080> PING
PONG
```

Data is persisted to `database.aof` in the working directory and replayed automatically on next startup.

## Running Tests

```bash
go test ./internal/test/ -race
```

The test suite covers store unit tests, full connection integration tests (via `net.Pipe`), AOF persistence round trips, and concurrency correctness under the race detector.

## Project Structure

```
redis-mini/
├── cmd/
│   └── main.go                  # Entry point, TCP listener, AOF startup replay
├── internal/
│   ├── data/
│   │   └── store.go             # In-memory store with RWMutex
│   ├── handlers/
│   │   ├── handle_resp.go       # RESP parser (array + inline)
│   │   ├── handle_connection.go # Per-connection command loop
│   │   ├── handle_operations.go # SET / GET handlers
│   │   └── handle_AOF.go        # AOF write and read/replay
│   └── test/
│       ├── handle_resp_test.go  # RESP parser unit tests
│       └── persistence_concurrency_test.go  # Integration + concurrency tests
```
