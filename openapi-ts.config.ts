import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "http://localhost:8080/openapi.json",
  output: "apps/web/src/lib/api",
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
