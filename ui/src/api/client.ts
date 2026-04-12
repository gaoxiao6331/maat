import type { EnvRecord, PresignAssetUploadResponse, Project } from '../types';

type ApiResponse<T> = { data: T; error?: string };

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const resp = await fetch(path, {
    headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) },
    ...init,
  });
  const json = (await resp.json().catch(() => ({}))) as ApiResponse<T>;
  if (!resp.ok) {
    throw new Error(json.error || `request failed: ${resp.status}`);
  }
  return json.data;
}

function normalizeProject(p: any): Project {
  return {
    id: Number(p.id || 0),
    name: String(p.name || ''),
    deployPath: String(p.deploy_path || p.deployPath || ''),
    owner: String(p.owner || ''),
  };
}

function normalizeEnv(r: any): EnvRecord {
  return {
    projectId: Number(r.project_id || r.projectId || 0),
    envName: String(r.env_name || r.envName || ''),
    htmlContent: String(r.html_content || r.htmlContent || ''),
    buildId: String(r.build_id || r.buildId || ''),
    updatedAt: Number(r.updated_at || r.updatedAt || 0),
  };
}

export const api = {
  async listProjects(query = ''): Promise<Project[]> {
    const data = await request<any[]>(`/api/projects?query=${encodeURIComponent(query)}`);
    return (data || []).map(normalizeProject);
  },
  createProject(payload: { name: string; owner: string; deploy_path: string }) {
    return request<{ ok: boolean }>('/api/projects', { method: 'POST', body: JSON.stringify(payload) });
  },
  publishEnv(payload: { project_id: number; env_name: string; html_body: string }) {
    return request<{ ok: boolean }>('/api/publish', { method: 'POST', body: JSON.stringify(payload) });
  },
  async getEnv(projectId: number, envName: string): Promise<EnvRecord> {
    const data = await request<any>(`/api/env?project_id=${projectId}&env_name=${encodeURIComponent(envName)}`);
    return normalizeEnv(data);
  },
  presignAsset(payload: { project: string; build_id: string; file_key: string; expire_seconds: number }) {
    return request<{ object_key: string; upload_url: string }>('/api/assets/presign', {
      method: 'POST',
      body: JSON.stringify(payload),
    }).then((r) => ({
      objectKey: r.object_key,
      uploadUrl: r.upload_url,
    }) as PresignAssetUploadResponse);
  },
};
