import { useQuery } from '@tanstack/react-query';
import React from 'react';

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Time } from '@/components/ui/time';
import { orpc } from '@/lib/orpc';

export function RevokedSessionsTable() {
  const { data, isLoading } = useQuery(
    orpc.getAuthAdminSessionsStatesRevoked.queryOptions(),
  );

  const states = data?.body ?? [];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">Loading...</div>
    );
  }

  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Session ID</TableHead>
            <TableHead>Revoked at</TableHead>
            <TableHead>Revoked by</TableHead>
            <TableHead>Reason</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <React.Suspense fallback={<>Loading...</>}>
            {states.map((state) => (
              <TableRow key={state.session_id}>
                <TableCell>
                  <code className="text-xs">{state.session_id}</code>
                </TableCell>
                <TableCell>
                  <Time
                    date={state.revoked_at ? new Date(state.revoked_at) : null}
                  />
                </TableCell>
                <TableCell>
                  {state.revoked_by_user_id ? (
                    <code className="text-xs">{state.revoked_by_user_id}</code>
                  ) : (
                    <span className="text-muted-foreground">—</span>
                  )}
                </TableCell>
                <TableCell className="max-w-xs truncate">
                  {state.revoked_reason || (
                    <span className="text-muted-foreground">—</span>
                  )}
                </TableCell>
              </TableRow>
            ))}
            {states.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={4}
                  className="h-24 text-center text-muted-foreground"
                >
                  No revoked sessions.
                </TableCell>
              </TableRow>
            )}
          </React.Suspense>
        </TableBody>
      </Table>
    </div>
  );
}
