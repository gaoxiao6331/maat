import React from 'react';
import { Button, Empty, Form, Input, Space, Table, Toast } from '@douyinfe/semi-ui';
import { api } from '../api/client';
import type { Project } from '../types';

export function ProjectDashboard({ onOpenProject }: { onOpenProject: (p: Project) => void }) {
  const [loading, setLoading] = React.useState(false);
  const [projects, setProjects] = React.useState<Project[]>([]);
  const [query, setQuery] = React.useState('');

  async function loadProjects() {
    setLoading(true);
    try {
      const rows = await api.listProjects(query);
      setProjects(rows);
    } catch (err) {
      Toast.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }

  React.useEffect(() => {
    void loadProjects();
  }, []);

  return (
    <Space vertical align="start" style={{ width: '100%' }} spacing="medium">
      <Form layout="horizontal" onSubmit={loadProjects} style={{ width: '100%' }}>
        <Space>
          <Input value={query} onChange={(v) => setQuery(v)} placeholder="搜索项目/owner/路径" />
          <Button theme="solid" htmlType="submit" loading={loading}>
            搜索
          </Button>
        </Space>
      </Form>

      <Form
        labelPosition="left"
        style={{ background: '#fff', padding: 16, borderRadius: 10, width: '100%' }}
        onSubmit={async (values) => {
          try {
            await api.createProject(values as { name: string; owner: string; deploy_path: string });
            Toast.success('创建成功');
            await loadProjects();
          } catch (err) {
            Toast.error((err as Error).message);
          }
        }}
      >
        <Form.Input field="name" label="项目名称" rules={[{ required: true }]} />
        <Form.Input field="owner" label="Owner" rules={[{ required: true }]} />
        <Form.Input field="deploy_path" label="部署路径" initValue="/" rules={[{ required: true }]} />
        <Button theme="solid" htmlType="submit">
          创建项目
        </Button>
      </Form>

      {projects.length === 0 ? (
        <Empty title="暂无项目" description="先创建一个项目" />
      ) : (
        <Table
          columns={[
            { title: 'ID', dataIndex: 'id', width: 90 },
            { title: '项目名称', dataIndex: 'name' },
            { title: 'Owner', dataIndex: 'owner' },
            { title: '部署路径', dataIndex: 'deployPath' },
            {
              title: '操作',
              render: (_: unknown, record: Project) => <Button onClick={() => onOpenProject(record)}>进入环境</Button>,
            },
          ]}
          dataSource={projects}
          pagination={false}
        />
      )}
    </Space>
  );
}
