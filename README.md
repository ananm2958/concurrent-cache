# Concurrent In-Memory Cache Server

### Overview
This project implements a production-style cache system designed to handle high-concurrency workloads while maintaining simplicity, durability, observability, and correctness
It combines:
- Efficient in-memory data structures (map + doubly linked list for LRU)
- Thread-safe access using sync.RWMutex
- HTTP-based API for interaction
- Persistence via append-only logging (AOF) and periodic snapshots
- Metrics for monitoring latency, hit rate, and system performance

## Key Features
### Concurrent Cache
- Thread-safe cache using sync.RWMutex
- Supports high read/write throughput
- Maintains correctness while supporting concurrent access

### LRU Eviction
- Least Recently Used (LRU) policy
- Implemented using:
    - Hashmap for average O(1) lookup
    - Doubly-linked list for O(1) updates
- Automatically evicts older updates when cache capacity is exceeded

### TTL Expiration
- Each key has an expiration time
- Expired keys are
  - Removed lazily on access
  - Cleaned periodically in the background

### HTTP Server
- Implements REST-style endpoints:
  - GET /cache?key=...
  - POST /cache
  - DELETE /cache?key=...
- JSON-based request/response format
- Handles concurrent requests as well

### Persistence
#### Append Only File (AOF)
 - Logs every write operation (SET, DELETE)
 - Ensures durability across crashes
 - Replayed on startup to restore recent state
#### Snapshotting
 - Periodically writes full cache state to disk
 - Reduces recovery time
 - Stored as JSON entries
#### Recovery
 - On startup:
     1. Load Snapshot
     2. Replay AOF
- Restores cache state after unclean shutdown


### Metrics and Observability
Exposes Prometheus-style metrics via /metrics:
- Total requests
- Cache hits / misses
- Evictions
- Request latency (histogram)

Example:
cache_requests_total 1200
cache_hits_total 900
cache_misses_total 300

cache_request_duration_seconds_bucket{le="0.001"} 400
cache_request_duration_seconds_bucket{le="0.01"} 1000

cache_request_duration_seconds_sum 2.35
cache_request_duration_seconds_count 1200


### Architecture
<img width="1343" height="302" alt="image" src="https://github.com/user-attachments/assets/48d66e17-c699-4851-b542-22a460c294ea" />


Components
- cache/
  - Core data structure (LRU + TTL)
- server/
  - HTTP server and request handlers
- metrics/
  - Tracks performance and system stats
- persistence/
  -  AOF logging, snapshotting, recovery

### Request Flow
#### GET
<img width="910" height="647" alt="image" src="https://github.com/user-attachments/assets/ea9d4919-dabd-431a-adf8-6c187598854c" />

#### SET / DELETE
<img width="575" height="631" alt="image" src="https://github.com/user-attachments/assets/3bfa9b0f-57b2-4da5-b25d-348ede658513" />


### Getting Started
#### Prerequisites
Go 1.20+

#### Installation
git clone <repo-url>
cd concurrent-cache
go mod tidy
Run the server
go run src/main.go

#### Server starts on:
http://localhost:8080

#### Set a value
curl -X POST http://localhost:8080/cache \
  -d '{"key":"foo","value":"bar"}'
  
#### Get a value
- curl "http://localhost:8080/cache?key=foo"
- Delete a value
- curl -X DELETE "http://localhost:8080/cache?key=foo"

#### View metrics
- curl http://localhost:8080/metrics

#### Background Processes
- The system runs background goroutines for:

#### Periodic snapshotting
- Expired key cleanup
- (Optional) AOF batching / flushing

#### Testing
Run tests:

go test ./...

Includes:
- Unit tests for cache correctness
- Concurrency tests
- Integration tests for HTTP endpoints

#### Performance Goals
- High throughput under concurrent load
- Low latency (millisecond-level p99)
- Efficient memory usage with bounded capacity
