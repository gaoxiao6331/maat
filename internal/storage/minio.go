package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type MinioAssetStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioAssetStorage(cfg MinioConfig) (*MinioAssetStorage, error) {
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("invalid minio config")
	}

	c, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioAssetStorage{client: c, bucket: cfg.Bucket}, nil
}

func (s *MinioAssetStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
}

func (s *MinioAssetStorage) GenerateUploadURL(ctx context.Context, objectKey string, expireSeconds int64) (string, error) {
	expire := time.Duration(expireSeconds) * time.Second
	if expire <= 0 {
		expire = 10 * time.Minute
	}
	u, err := s.client.PresignedPutObject(ctx, s.bucket, objectKey, expire)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *MinioAssetStorage) Upload(ctx context.Context, objectKey string, body []byte, contentType string) error {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err := s.client.PutObject(ctx, s.bucket, objectKey, bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

func (s *MinioAssetStorage) UploadReader(ctx context.Context, objectKey string, r io.Reader, size int64, contentType string) error {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err := s.client.PutObject(ctx, s.bucket, objectKey, r, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}
