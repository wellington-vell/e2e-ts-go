import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, useRouteContext } from '@tanstack/react-router';
import React from 'react';

import { orpc } from '@/lib/orpc';
import { cn } from '@/lib/utils';

export const Route = createFileRoute('/')({
  component: HomeComponent,
  loader: async () => ({
    healthCheck: orpc.getApiV1Health.queryOptions(),
  }),
});

function HomeComponent() {
  const healthCheck = useSuspenseQuery(orpc.getApiV1Health.queryOptions());
  const { session } = useRouteContext({ from: '__root__' });

  return (
    <div className="container mx-auto px-4 py-2 my-auto">
      {JSON.stringify(session, null, 2)}
      <div className="grid gap-6">
        <React.Suspense fallback={<div>Loading...</div>}>
          <section className="rounded-lg border p-4">
            <h2 className="mb-2 font-medium">API Status</h2>
            <div className="flex items-center gap-2">
              <div
                className={cn(
                  'h-2 w-2 rounded-full',
                  healthCheck.data ? 'bg-success' : 'bg-destructive',
                )}
              />
              <span className="text-sm text-muted-foreground">
                {healthCheck.isLoading
                  ? 'Checking...'
                  : healthCheck.data
                    ? 'Connected'
                    : 'Disconnected'}
              </span>
              <pre className="overflow-x-auto font-mono text-sm">
                {JSON.stringify(healthCheck.data, null, 2)}
              </pre>
            </div>
          </section>
        </React.Suspense>
      </div>
    </div>
  );
}
