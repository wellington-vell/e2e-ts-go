import { experimental_toORPCClient } from '@orpc/hey-api';
import { createTanstackQueryUtils } from '@orpc/tanstack-query';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

import { client } from '@/lib/api/client.gen';
import * as sdk from '@/lib/api/sdk.gen';
import { env } from '@/lib/env';

function formatErrorMessage(error: Error): string {
  let message = error.message;
  try {
    const parsed = JSON.parse(message);
    if (parsed?.message) {
      message = parsed.message;
    }
  } catch {
    // Not JSON, keep original
  }
  return message;
}

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      if (
        error.message.includes('401') ||
        error.message.includes('Unauthorized')
      ) {
        return;
      }
      const message = formatErrorMessage(error);
      toast.error(`Error: ${message}`, {
        action: {
          label: 'retry',
          onClick: () => query.invalidate(),
        },
      });
    },
  }),
  mutationCache: new MutationCache({
    onError: (error, _variables, _context, mutation) => {
      const message = formatErrorMessage(error);
      toast.error(`Error: ${message}`, {
        action: {
          label: 'retry',
          onClick: () => mutation.execute(_variables),
        },
      });
    },
  }),
});

client.setConfig({
  baseUrl: env.VITE_SERVER_URL,
  credentials: 'include',
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
