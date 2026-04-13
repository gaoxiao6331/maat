import { appTools, defineConfig } from '@modern-js/app-tools';

export default defineConfig({
  plugins: [appTools({ bundler: 'rspack' })],
  tools: {
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
  server: {
    port: 5173,
  },
});
