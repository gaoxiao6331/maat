import React from 'react';
import { Button, Card, Form, Space, Tag, Toast, Typography } from '@douyinfe/semi-ui';
import { api } from '../api/client';
import type { EnvRecord, Project } from '../types';

const ENV_NAMES = ['production', 'test', 'feature-a', 'feature-b'];

export function EnvironmentDetail({ project, onBack }: { project: Project; onBack: () => void }) {
  const [envData, setEnvData] = React.useState<Record<string, EnvRecord | null>>({});

  async function loadEnv(env: string) {
    try {
      const rec = await api.getEnv(project.id, env);
      setEnvData((prev) => ({ ...prev, [env]: rec }));
    } catch {
      setEnvData((prev) => ({ ...prev, [env]: null }));
    }
  }

  React.useEffect(() => {
    ENV_NAMES.forEach((env) => {
      void loadEnv(env);
    });
  }, [project.id]);

  return (
    <Space vertical align="start" style={{ width: '100%' }} spacing="medium">
      <Space>
        <Button onClick={onBack}>返回</Button>
        <Typography.Title heading={4} style={{ margin: 0 }}>
          {project.name}
        </Typography.Title>
        <Tag color="cyan">{project.deployPath}</Tag>
      </Space>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill,minmax(280px,1fr))', gap: 12, width: '100%' }}>
        {ENV_NAMES.map((env) => {
          const rec = envData[env];
          return (
            <Card key={env} title={env} shadows="hover" bordered>
              <Space vertical align="start">
                <Typography.Text>构建ID: {rec?.buildId || '-'}</Typography.Text>
                <Typography.Text>更新时间: {rec?.updatedAt || '-'}</Typography.Text>
                <Button onClick={() => window.open(`${project.deployPath}/index.html?x_m_env=${env}`, '_blank')}>预览</Button>
              </Space>
            </Card>
          );
        })}
      </div>

      <Card title="发布 HTML" style={{ width: '100%' }}>
        <Form
          onSubmit={async (v: unknown) => {
            try {
              const values = v as { env_name: string; html_body: string };
              await api.publishEnv({ project_id: project.id, env_name: values.env_name, html_body: values.html_body });
              Toast.success('发布成功');
              await loadEnv(values.env_name);
            } catch (err) {
              Toast.error((err as Error).message);
            }
          }}
        >
          <Form.Select
            field="env_name"
            label="环境"
            rules={[{ required: true }]}
            optionList={ENV_NAMES.map((x) => ({ label: x, value: x }))}
            initValue="test"
          />
          <Form.TextArea field="html_body" label="HTML 内容" rows={10} rules={[{ required: true }]} />
          <Button htmlType="submit" theme="solid">
            发布
          </Button>
        </Form>
      </Card>

      <Card title="获取 MinIO 预签名上传 URL" style={{ width: '100%' }}>
        <Form
          onSubmit={async (v: unknown) => {
            try {
              const values = v as { build_id: string; file_key: string; expire_seconds?: string };
              const res = await api.presignAsset({
                project: project.name,
                build_id: values.build_id,
                file_key: values.file_key,
                expire_seconds: Number(values.expire_seconds) || 600,
              });
              Toast.success(`已生成: ${res.uploadUrl}`);
            } catch (err) {
              Toast.error((err as Error).message);
            }
          }}
        >
          <Form.Input field="build_id" label="Build ID" rules={[{ required: true }]} />
          <Form.Input field="file_key" label="资源名" placeholder="main.js" rules={[{ required: true }]} />
          <Form.Input field="expire_seconds" label="过期秒数" initValue="600" />
          <Button htmlType="submit">生成上传URL</Button>
        </Form>
      </Card>
    </Space>
  );
}
