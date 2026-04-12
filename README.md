# maat mvp (protobuf + idl generation)

当前实现已切换为：

- IDL：`protobuf`
- 后端契约类型：由 Hertz 生成（`hz`）
- 前端契约类型：由 `ts-proto` 生成
- 前端工程：`TypeScript + Modern.js + Rspack`

## IDL

- `idl/release_service.proto`

## 代码生成

1. 生成 Go（Hertz）

```bash
cd /Users/go/Code/maat/code
./scripts/gen-go-from-proto.sh
```

2. 生成前端 TS 类型（ts-proto）

```bash
cd /Users/go/Code/maat/code/ui
pnpm install
pnpm run gen:proto
```

说明：
- 若无 `hz`：`go install github.com/cloudwego/hertz/cmd/hz@latest`
- 若无 `protoc`：先安装 protobuf 编译器

## 启动后端

```bash
cd /Users/go/Code/maat/code
./scripts/start.sh
```

## 启动前端（Modern.js + Rspack）

```bash
cd /Users/go/Code/maat/code
./scripts/start-ui.sh
```

前端默认：`http://localhost:5173`

## MinIO 可选配置

```bash
export MAAT_ASSET_DRIVER=minio
export MAAT_MINIO_ENDPOINT=127.0.0.1:9000
export MAAT_MINIO_ACCESS_KEY=minioadmin
export MAAT_MINIO_SECRET_KEY=minioadmin
export MAAT_MINIO_BUCKET=maat-assets
export MAAT_MINIO_USE_SSL=false
./scripts/start.sh
```
