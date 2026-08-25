#!/usr/bin/env bash
set -euo pipefail

base_url="${1:-http://localhost:8080}"
workers="${WORKERS:-50}"
requests="${REQUESTS:-50000}"

curl -fsS -X POST \
  -H 'Content-Type: application/json' \
  -d '{"key":"load-test","value":"ready"}' \
  "$base_url/cache" >/dev/null

seq 1 "$requests" | xargs -P "$workers" -n 1 \
  curl -fsS "$base_url/cache?key=load-test" >/dev/null &
load_pid=$!

for sample in $(seq 1 5); do
  sleep 1
  echo "sample $sample"
  curl -fsS "$base_url/metrics" |
    grep -E 'cache_(requests_per_second|concurrent_connections|requests_total)'
done

wait "$load_pid"
