package store

import (
	"context"

	"maat/internal/model"
)

// MetadataStore 抽象项目和环境元数据存储，便于替换为 PostgreSQL / Redis 等实现。
type MetadataStore interface {
	CreateProject(ctx context.Context, p model.Project) (bool, error)
	ListProjects(ctx context.Context, query string) ([]model.Project, error)
	FindProjectByPath(ctx context.Context, path string) (*model.Project, error)
	PublishEnv(ctx context.Context, projectID int64, envName, htmlBody, buildID string) (bool, error)
	GetEnvDetail(ctx context.Context, projectID int64, envName string) (*model.EnvRecord, error)
	Close() error
}
