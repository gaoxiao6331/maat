#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v go >/dev/null 2>&1; then
  echo "go is required but not found in PATH"
  exit 1
fi

export MAAT_SQLITE_DSN="${MAAT_SQLITE_DSN:-file:maat.db?_pragma=busy_timeout(5000)}"

echo "[maat] using sqlite dsn: ${MAAT_SQLITE_DSN}"
echo "[maat] starting server on :8080"

go mod tidy
go run ./cmd/maat-server
