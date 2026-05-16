import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: {
    path: 'apps/server/docs/swagger.json',
  },
  output: {
    path: 'apps/web/src/lib/api',
    entryFile: false,
    postProcess: ['oxfmt'],
  },
  parser: {
    transforms: {
      schemaName: (name) => {
        return name
          .replace(/^server_internal_models\./i, '')
          .replace(/^types\./i, '')
          .replace(/^models\./i, '');
      },
    },
  },
  plugins: [
    {
      name: 'zod',
      requests: true,
      responses: false,
      case: 'snake_case',
      comments: false,
    },
    {
      name: '@hey-api/sdk',
      validator: true,
      comments: false,
    },
  ],
});
