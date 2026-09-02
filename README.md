# Concurrent Cache

An HTTP in-memory cache with LRU eviction, TTL, crash recovery, and Prometheus-style metrics. Callers GET, SET, and DELETE keys over HTTP without depending on Redis.

## Features

- **Concurrent cache** — thread-safe access via `sync.RWMutex` around a hashmap and doubly linked list
- **LRU eviction** — O(1) lookup and update; evicts the least recently used key when capacity is exceeded
- **TTL expiration** — lazy removal on access plus a background sweep of expired keys
- **HTTP API** — REST-style `GET` / `POST` / `PUT` / `DELETE` on `/cache` with JSON responses
- **Persistence** — append-only log for writes, periodic JSON snapshots, snapshot-then-AOF recovery on startup
- **Observability** — Prometheus text metrics on `/metrics` (hits, misses, latency, connections)

## Tech Stack

| Layer | Technology |
| --- | --- |
| Language | Go 1.20+ |
| HTTP | stdlib `net/http` |
| Cache | Hashmap + doubly linked list (LRU) |
| Concurrency | `sync.RWMutex` |
| Persistence | JSON AOF + JSON snapshots |
| Metrics | Prometheus exposition format |

## Project Structure

```
concurrent-cache/
├── src/
│   ├── main.go
│   ├── cache/          # LRU + TTL
│   ├── server/         # HTTP handlers + load tests
│   ├── metrics/        # counters, histogram, connection stats
│   └── persistence/    # AOF, snapshot, recovery
├── load_test.sh
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.20+

## Getting Started

### 1. Clone the repository

```bash
git clone <repo-url>
cd concurrent-cache
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

```bash
go run src/main.go
```

The server starts on [http://localhost:8080](http://localhost:8080).

### Runtime defaults

These values are set in `src/main.go`:

| Setting | Default |
| --- | --- |
| Listen port | `8080` |
| Capacity | `10000` keys |
| Default TTL (in-process `Cache.Set`) | `60s` |
| Expired-key sweep | every `10s` |
| Snapshot interval | every `1m` |
| Snapshot file | `snapshot.json` |
| AOF file | `appendonly.aof` |

Keys set through the HTTP API with omitted or `0` `ttl_seconds` do not expire (zero expiry).

## API Reference

All cache routes live under `/cache`. Errors return JSON `{"error": "..."}`.

| Status | Meaning |
| --- | --- |
| `200` | Success |
| `400` | Missing key or invalid JSON |
| `404` | Key not found |
| `405` | Method not allowed |
| `500` | AOF write failed |

### Get a value

`GET /cache?key=...`

```bash
curl "http://localhost:8080/cache?key=foo"
```

```json
{"key":"foo","value":"bar"}
```

### Set a value

`POST` or `PUT /cache`

JSON body:

```bash
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{"key":"foo","value":"bar","ttl_seconds":60}'
```

If `Content-Type` is not `application/json`, query parameters work instead: `key`, `value`, and optional `ttl`.

```bash
curl -X POST "http://localhost:8080/cache?key=foo&value=bar&ttl=60"
```

### Delete a value

`DELETE /cache?key=...`

```bash
curl -X DELETE "http://localhost:8080/cache?key=foo"
```

```json
{"message":"deleted"}
```

### Metrics

`GET /metrics`

```bash
curl http://localhost:8080/metrics
```

Exposes Prometheus text format. No auth required.

## Persistence

Writes go through two durability paths:

- **AOF** — every SET and DELETE is appended as a JSON line to `appendonly.aof`. A failed append returns HTTP 500 so the caller knows the write was not logged.
- **Snapshots** — once a minute the process dumps live, non-expired entries to `snapshot.json`.

On startup the server:

1. Loads `snapshot.json` (missing file is ignored)
2. Replays `appendonly.aof` on top
3. Skips entries whose expiry is already in the past

That restores state after an unclean shutdown without replaying the full history from an empty cache.

## Observability

`GET /metrics` emits:

| Metric | Description |
| --- | --- |
| `cache_requests_total` | Total `/cache` requests |
| `cache_hits_total` | Successful GET and DELETE |
| `cache_misses_total` | GET/DELETE for missing keys |
| `cache_evictions_total` | LRU evictions (counter exists; not incremented from the cache today, so this stays `0`) |
| `cache_requests_per_second` | Requests in the last 1 second |
| `cache_connections_total` | Accepted HTTP connections |
| `cache_concurrent_connections` | Open HTTP connections |
| `cache_request_duration_seconds_bucket` | Latency histogram (`le` of 0.001, 0.01, 0.1, 1, `+Inf`) |
| `cache_request_duration_seconds_sum` | Sum of request latency in seconds |
| `cache_request_duration_seconds_count` | Requests included in the histogram |

## Testing

Unit tests:

```bash
go test ./src/cache ./src/metrics ./src/server
```

Connection load test (starts a test server in-process). Set any positive connection count with `CACHE_LOAD_CONNECTIONS`:

```bash
CACHE_LOAD_CONNECTIONS=250 go test -v ./src/server -run '^TestMaxConcurrentConnections$' -count=1
```

Against a running server, `./load_test.sh` hammers `GET /cache` and samples throughput metrics:

```bash
./load_test.sh http://localhost:8080
```
