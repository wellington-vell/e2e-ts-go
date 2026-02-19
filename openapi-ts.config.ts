import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "http://localhost:8080/openapi.json",
  output: "apps/web/src/lib/api",
});
