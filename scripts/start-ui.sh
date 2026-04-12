#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR/ui"

if ! command -v pnpm >/dev/null 2>&1; then
  echo "pnpm is required but not found in PATH"
  exit 1
fi

echo "[maat-ui] installing dependencies"
pnpm install

echo "[maat-ui] starting dev server on :5173"
pnpm dev
