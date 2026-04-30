import { TanStackDevtools } from '@tanstack/react-devtools';
import { formDevtoolsPlugin } from '@tanstack/react-form-devtools';
import type { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools';
import {
  HeadContent,
  Outlet,
  createRootRouteWithContext,
} from '@tanstack/react-router';
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';

import { Header } from '@/components/header';
import { Toaster } from '@/components/ui/sonner';
import { ThemeProvider } from '@/context/theme';
import type { GetMeResponse } from '@/lib/api/types.gen';
import { env } from '@/lib/env';
import { orpc } from '@/lib/orpc';

import '@/index.css';

interface RouterContext {
  queryClient: QueryClient;
  orpc: typeof orpc;
  session: GetMeResponse | null;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  beforeLoad: async ({ context }) => {
    try {
      const response = await context.queryClient.fetchQuery(
        orpc.getAuthMe.queryOptions(),
      );
      return { session: response.body ?? null };
    } catch {
      return { session: null };
    }
  },
  head: () => ({
    meta: [
      {
        title: 'E2E GO/TS',
      },
      {
        name: 'description',
        content: 'E2E GO/TS is a web application for testing GO/TS integration',
      },
    ],
    links: [
      {
        rel: 'icon',
        href: '/favicon.ico',
      },
    ],
  }),
});

function RootComponent() {
  return (
    <>
      <HeadContent />
      <ThemeProvider
        attribute="class"
        defaultTheme="dark"
        disableTransitionOnChange
        storageKey="vite-ui-theme"
      >
        <main className="grid grid-rows-[auto_1fr] h-svh">
          <Header />
          <Outlet />
        </main>
        <Toaster richColors />
      </ThemeProvider>

      {env.VITE_NODE_ENV === 'development' && (
        <TanStackDevtools
          config={{ position: 'bottom-left', panelLocation: 'bottom' }}
          plugins={[
            formDevtoolsPlugin(),
            {
              name: 'TanStack Query',
              render: <ReactQueryDevtoolsPanel />,
            },
            {
              name: 'TanStack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
          ]}
        />
      )}
    </>
  );
}
