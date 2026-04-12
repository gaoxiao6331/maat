package storage

import "context"

// AssetStorage 抽象静态资源托管能力，便于替换为 MinIO / CDN API / S3 等实现。
type AssetStorage interface {
	GenerateUploadURL(ctx context.Context, objectKey string, expireSeconds int64) (string, error)
	Upload(ctx context.Context, objectKey string, body []byte, contentType string) error
}

// NoopAssetStorage 是 MVP 的占位实现，后续可替换为 MinIO、CDN API 等。
type NoopAssetStorage struct{}

func NewNoopAssetStorage() *NoopAssetStorage {
	return &NoopAssetStorage{}
}

func (s *NoopAssetStorage) GenerateUploadURL(_ context.Context, _ string, _ int64) (string, error) {
	return "", nil
}

func (s *NoopAssetStorage) Upload(_ context.Context, _ string, _ []byte, _ string) error {
	return nil
}
