import React from 'react';
import { Layout, Nav } from '@douyinfe/semi-ui';
import { ProjectDashboard } from './pages/ProjectDashboard';
import { EnvironmentDetail } from './pages/EnvironmentDetail';
import type { Project } from './types';

const { Header, Content } = Layout;

export function App() {
  const [view, setView] = React.useState<'projects' | 'env'>('projects');
  const [project, setProject] = React.useState<Project | null>(null);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#0f172a' }}>
        <Nav
          mode="horizontal"
          style={{ background: '#0f172a' }}
          selectedKeys={[view]}
          onSelect={({ itemKey }) => {
            if (itemKey === 'projects') {
              setView('projects');
              setProject(null);
            }
          }}
          items={[{ itemKey: 'projects', text: '项目列表' }]}
          header={{ text: 'MAAT 发布管理台' }}
        />
      </Header>
      <Content style={{ padding: 24 }}>
        <div className="mx-auto w-full max-w-7xl">
          {view === 'projects' && (
            <ProjectDashboard
              onOpenProject={(p) => {
                setProject(p);
                setView('env');
              }}
            />
          )}
          {view === 'env' && project && (
            <EnvironmentDetail
              project={project}
              onBack={() => {
                setView('projects');
                setProject(null);
              }}
            />
          )}
        </div>
      </Content>
    </Layout>
  );
}
