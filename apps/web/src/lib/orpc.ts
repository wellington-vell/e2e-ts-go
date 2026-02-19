import { experimental_toORPCClient } from "@orpc/hey-api";
import { createTanstackQueryUtils } from "@orpc/tanstack-query";
import { QueryCache, QueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { client } from "@/lib/api/client.gen";
import * as sdk from "@/lib/api/sdk.gen";
import { env } from "@/lib/env";

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      toast.error(`Error: ${error.message}`, {
        action: {
          label: "retry",
          onClick: query.invalidate,
        },
      });
    },
  }),
});

client.setConfig({
  baseUrl: env.VITE_SERVER_URL,
  credentials: "include",
});

client.interceptors.error.use((error) => {
  return error;
});

client.interceptors.response.use(async (response, _request, _opts) => {
  return response;
});

client.interceptors.request.use((request, _opts) => {
  return request;
});

const apiClient = experimental_toORPCClient(sdk);

export const orpc = createTanstackQueryUtils(apiClient);
