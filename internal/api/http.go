package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"maat/internal/gateway"
	"maat/internal/model"
	"maat/internal/service"
)

type HTTPHandler struct {
	svc *service.ReleaseService
}

func NewHTTPHandler(svc *service.ReleaseService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) Register(r *route.Engine) {
	r.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, map[string]any{"ok": true})
	})

	r.GET("/api/projects", h.ListProjects)
	r.POST("/api/projects", h.CreateProject)
	r.POST("/api/publish", h.PublishEnv)
	r.GET("/api/env", h.GetEnvDetail)
	r.POST("/api/assets/presign", h.PresignAssetUpload)

	r.GET("/*path", h.GatewayHTML)
}

func (h *HTTPHandler) ListProjects(ctx context.Context, c *app.RequestContext) {
	query := string(c.Query("query"))
	projects, err := h.svc.ListProjects(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]any{"data": projects})
}

func (h *HTTPHandler) CreateProject(ctx context.Context, c *app.RequestContext) {
	var req model.Project
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	ok, err := h.svc.CreateProject(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]any{"ok": ok})
}

func (h *HTTPHandler) PublishEnv(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ProjectId int64  `json:"project_id"`
		EnvName   string `json:"env_name"`
		HTMLBody  string `json:"html_body"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	ok, err := h.svc.PublishEnv(ctx, req.ProjectId, req.EnvName, req.HTMLBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]any{"ok": ok})
}

func (h *HTTPHandler) GetEnvDetail(ctx context.Context, c *app.RequestContext) {
	projectID, err := strconv.ParseInt(string(c.Query("project_id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid project_id"})
		return
	}
	envName := string(c.Query("env_name"))
	rec, err := h.svc.GetEnvDetail(ctx, projectID, envName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, map[string]any{"error": "env not found"})
			return
		}
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]any{"data": rec})
}

func (h *HTTPHandler) GatewayHTML(ctx context.Context, c *app.RequestContext) {
	path := string(c.Param("path"))
	queryEnv := string(c.Query("x_m_env"))
	cookieEnv := string(c.Cookie("x_m_env"))

	html, envName, err := h.svc.ResolveHTMLByRequest(ctx, path, queryEnv, cookieEnv)
	if err != nil {
		c.String(http.StatusNotFound, "maat gateway: html not found")
		return
	}

	out := gateway.InjectEnv(html, envName)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, out)
}

func (h *HTTPHandler) PresignAssetUpload(ctx context.Context, c *app.RequestContext) {
	var req model.PresignAssetReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	resp, err := h.svc.PresignAssetUpload(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]any{"data": resp})
}
