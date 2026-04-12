package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/hertz/pkg/app/server"
	"maat/internal/api"
	"maat/internal/service"
	"maat/internal/storage"
	"maat/internal/store"
)

func main() {
	dsn := os.Getenv("MAAT_SQLITE_DSN")
	if dsn == "" {
		dsn = "file:maat.db?_pragma=busy_timeout(5000)"
	}

	s, err := store.NewSQLiteStore(dsn)
	if err != nil {
		log.Fatalf("init store failed: %v", err)
	}
	defer s.Close()

	svc := service.NewReleaseService(s, storage.NewNoopAssetStorage())
	if assetStorage := initAssetStorage(); assetStorage != nil {
		svc = service.NewReleaseService(s, assetStorage)
	}
	h := api.NewHTTPHandler(svc)

	r := server.New(server.WithHostPorts(":8080"))
	h.Register(r)

	log.Println("maat mvp server listening on :8080")
	r.Spin()
}

func initAssetStorage() storage.AssetStorage {
	if strings.TrimSpace(os.Getenv("MAAT_ASSET_DRIVER")) != "minio" {
		return nil
	}
	cfg := storage.MinioConfig{
		Endpoint:  os.Getenv("MAAT_MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MAAT_MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MAAT_MINIO_SECRET_KEY"),
		UseSSL:    os.Getenv("MAAT_MINIO_USE_SSL") == "true",
		Bucket:    os.Getenv("MAAT_MINIO_BUCKET"),
	}
	m, err := storage.NewMinioAssetStorage(cfg)
	if err != nil {
		log.Printf("minio disabled: %v", err)
		return nil
	}
	if err := m.EnsureBucket(context.Background()); err != nil {
		log.Printf("minio ensure bucket failed: %v", err)
		return nil
	}
	log.Printf("minio asset storage enabled, bucket=%s", cfg.Bucket)
	return m
}
