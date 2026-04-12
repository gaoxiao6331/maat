#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v hz >/dev/null 2>&1; then
  echo "hz not found. Install with: go install github.com/cloudwego/hertz/cmd/hz@latest"
  exit 1
fi

IDL_FILE="$ROOT_DIR/idl/release_service.proto"
OUT_DIR="$ROOT_DIR/idl/gen/go"

mkdir -p "$OUT_DIR"

# Hertz protobuf generation. Generated files are used as the backend API contract types.
hz pb \
  -I "$ROOT_DIR/idl" \
  -proto "$IDL_FILE" \
  -module maat \
  -out_dir "$OUT_DIR"

echo "generated go protobuf files in $OUT_DIR"
