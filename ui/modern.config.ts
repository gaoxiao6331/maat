import { appTools, defineConfig } from '@modern-js/app-tools';

export default defineConfig({
  plugins: [appTools({ bundler: 'rspack' })],
  tools: {
    devServer: {
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:8080',
          changeOrigin: true,
        },
      },
    },
    rspack: (config) => {
      config.resolve = config.resolve || {};
      config.resolve.conditionNames = ['require', 'import', 'browser', 'default'];
    },
  },
  source: {
    disableDefaultEntries: true,
    entries: {
      index: './src/main.tsx',
    },
  },
  dev: {
    port: 9000,
  },
});
