import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: { path: "http://localhost:8080/openapi.json" },
  output: {
    path: "apps/web/src/lib/api",
    entryFile: false,
    postProcess: ["oxfmt"],
  },
  plugins: [
    {
      name: "zod",
      requests: true,
      responses: false,
    },
    {
      name: "@hey-api/sdk",
      validator: true,
    },
  ],
});
