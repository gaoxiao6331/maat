#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v protoc >/dev/null 2>&1; then
  echo "protoc not found. please install protoc first."
  exit 1
fi

if [ ! -x "$ROOT_DIR/ui/node_modules/.bin/protoc-gen-ts_proto" ]; then
  echo "ts-proto plugin not found. run: cd ui && pnpm install"
  exit 1
fi

mkdir -p "$ROOT_DIR/ui/src/gen"

protoc \
  -I "$ROOT_DIR/idl" \
  --plugin=protoc-gen-ts_proto="$ROOT_DIR/ui/node_modules/.bin/protoc-gen-ts_proto" \
  --ts_proto_out="$ROOT_DIR/ui/src/gen" \
  --ts_proto_opt=esModuleInterop=true,outputClientImpl=false,useOptionals=messages,env=browser \
  "$ROOT_DIR/idl/release_service.proto"

echo "generated ts-proto files in ui/src/gen"
