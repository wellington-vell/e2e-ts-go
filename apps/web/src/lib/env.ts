import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const env = createEnv({
  clientPrefix: "VITE_",
  client: {
    VITE_WEB_PORT: z.coerce.number().int().positive(),
  },
  runtimeEnv: import.meta.env,
  emptyStringAsUndefined: true,
});
