package service

import (
	"context"
	"database/sql"
	"errors"
	"path"
	"strings"
	"time"

	"maat/internal/model"
	"maat/internal/storage"
	"maat/internal/store"
)

type ReleaseService struct {
	store  store.MetadataStore
	assets storage.AssetStorage
}

func NewReleaseService(s store.MetadataStore, assets storage.AssetStorage) *ReleaseService {
	if assets == nil {
		assets = storage.NewNoopAssetStorage()
	}
	return &ReleaseService{store: s, assets: assets}
}

func (s *ReleaseService) ListProjects(ctx context.Context, query string) ([]model.Project, error) {
	return s.store.ListProjects(ctx, query)
}

func (s *ReleaseService) CreateProject(ctx context.Context, req model.Project) (bool, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.DeployPath = strings.TrimSpace(req.DeployPath)
	req.Owner = strings.TrimSpace(req.Owner)
	if req.Name == "" || req.DeployPath == "" || req.Owner == "" {
		return false, errors.New("name/deploy_path/owner are required")
	}
	if !strings.HasPrefix(req.DeployPath, "/") {
		req.DeployPath = "/" + req.DeployPath
	}
	return s.store.CreateProject(ctx, req)
}

func (s *ReleaseService) PublishEnv(ctx context.Context, projectID int64, envName, htmlBody string) (bool, error) {
	envName = strings.TrimSpace(envName)
	htmlBody = strings.TrimSpace(htmlBody)
	if projectID <= 0 || envName == "" || htmlBody == "" {
		return false, errors.New("project_id/env_name/html_body are required")
	}
	buildID := time.Now().Format("20060102T150405")
	return s.store.PublishEnv(ctx, projectID, envName, htmlBody, buildID)
}

func (s *ReleaseService) GetEnvDetail(ctx context.Context, projectID int64, envName string) (*model.EnvRecord, error) {
	if projectID <= 0 || strings.TrimSpace(envName) == "" {
		return nil, errors.New("project_id/env_name are required")
	}
	return s.store.GetEnvDetail(ctx, projectID, envName)
}

func (s *ReleaseService) ResolveHTMLByRequest(ctx context.Context, path, queryEnv, cookieEnv string) (string, string, error) {
	project, err := s.store.FindProjectByPath(ctx, path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", sql.ErrNoRows
		}
		return "", "", err
	}

	envName := strings.TrimSpace(queryEnv)
	if envName == "" {
		envName = strings.TrimSpace(cookieEnv)
	}
	if envName == "" {
		envName = "production"
	}

	rec, err := s.store.GetEnvDetail(ctx, project.Id, envName)
	if err != nil {
		if envName != "production" && errors.Is(err, sql.ErrNoRows) {
			rec, err = s.store.GetEnvDetail(ctx, project.Id, "production")
			envName = "production"
		}
		if err != nil {
			return "", "", err
		}
	}
	return rec.HtmlContent, envName, nil
}

func (s *ReleaseService) PresignAssetUpload(ctx context.Context, req model.PresignAssetReq) (*model.PresignAssetResp, error) {
	req.Project = strings.TrimSpace(req.Project)
	req.BuildId = strings.TrimSpace(req.BuildId)
	req.FileKey = strings.TrimSpace(req.FileKey)
	if req.Project == "" || req.BuildId == "" || req.FileKey == "" {
		return nil, errors.New("project/build_id/file_key are required")
	}
	objectKey := path.Join(req.Project, req.BuildId, "assets", req.FileKey)
	uploadURL, err := s.assets.GenerateUploadURL(ctx, objectKey, req.ExpireSeconds)
	if err != nil {
		return nil, err
	}
	return &model.PresignAssetResp{
		ObjectKey: objectKey,
		UploadUrl: uploadURL,
	}, nil
}
