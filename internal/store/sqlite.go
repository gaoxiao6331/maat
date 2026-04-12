package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"maat/internal/model"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	s := &SQLiteStore{db: db}
	if err := s.initSchema(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) initSchema(ctx context.Context) error {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			deploy_path TEXT NOT NULL,
			owner TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS env_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			env_name TEXT NOT NULL,
			html_content TEXT NOT NULL,
			build_id TEXT NOT NULL,
			updated_at INTEGER NOT NULL,
			UNIQUE(project_id, env_name),
			FOREIGN KEY(project_id) REFERENCES projects(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_projects_deploy_path ON projects(deploy_path);`,
		`CREATE INDEX IF NOT EXISTS idx_env_records_project_env ON env_records(project_id, env_name);`,
	}
	for _, q := range schema {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteStore) CreateProject(ctx context.Context, p model.Project) (bool, error) {
	now := time.Now().Unix()
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO projects(name, deploy_path, owner, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		p.Name, p.DeployPath, p.Owner, now, now,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SQLiteStore) ListProjects(ctx context.Context, query string) ([]model.Project, error) {
	query = strings.TrimSpace(query)
	base := `SELECT id, name, deploy_path, owner FROM projects`
	args := []any{}
	if query != "" {
		base += ` WHERE name LIKE ? OR owner LIKE ? OR deploy_path LIKE ?`
		like := fmt.Sprintf("%%%s%%", query)
		args = append(args, like, like, like)
	}
	base += ` ORDER BY updated_at DESC, id DESC`

	rows, err := s.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.Id, &p.Name, &p.DeployPath, &p.Owner); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (s *SQLiteStore) FindProjectByPath(ctx context.Context, path string) (*model.Project, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, sql.ErrNoRows
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	rows, err := s.db.QueryContext(ctx, `SELECT id, name, deploy_path, owner FROM projects ORDER BY LENGTH(deploy_path) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.Id, &p.Name, &p.DeployPath, &p.Owner); err != nil {
			return nil, err
		}
		if strings.HasPrefix(path, p.DeployPath) {
			return &p, nil
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, sql.ErrNoRows
}

func (s *SQLiteStore) PublishEnv(ctx context.Context, projectID int64, envName, htmlBody, buildID string) (bool, error) {
	if projectID <= 0 || envName == "" || htmlBody == "" || buildID == "" {
		return false, errors.New("invalid publish payload")
	}
	now := time.Now().Unix()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO env_records(project_id, env_name, html_content, build_id, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, env_name)
		 DO UPDATE SET html_content=excluded.html_content, build_id=excluded.build_id, updated_at=excluded.updated_at`,
		projectID, envName, htmlBody, buildID, now,
	)
	if err != nil {
		return false, err
	}

	_, err = tx.ExecContext(ctx, `UPDATE projects SET updated_at=? WHERE id=?`, now, projectID)
	if err != nil {
		return false, err
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (s *SQLiteStore) GetEnvDetail(ctx context.Context, projectID int64, envName string) (*model.EnvRecord, error) {
	var rec model.EnvRecord
	err := s.db.QueryRowContext(
		ctx,
		`SELECT project_id, env_name, html_content, build_id, updated_at
		 FROM env_records WHERE project_id=? AND env_name=?`,
		projectID, envName,
	).Scan(&rec.ProjectId, &rec.EnvName, &rec.HtmlContent, &rec.BuildId, &rec.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}
