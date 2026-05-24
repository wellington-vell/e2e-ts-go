import { keepPreviousData, useQuery } from '@tanstack/react-query';
import React from 'react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Time } from '@/components/ui/time';
import type { AdminUserState } from '@/lib/api/types.gen';
import { orpc } from '@/lib/orpc';
import { AdminUsersProvider } from '@/routes/(auth)/admin/-components/context';
import { CreateUserSheet } from '@/routes/(auth)/admin/-components/create-user-sheet';
import { RowActions } from '@/routes/(auth)/admin/-components/row-actions';

const PAGE_SIZE = 10;

export function UsersTable() {
  const [cursor, setCursor] = React.useState<string | undefined>(undefined);
  const [cursorHistory, setCursorHistory] = React.useState<
    (string | undefined)[]
  >([]);

  const {
    data: usersData,
    isLoading: usersLoading,
    isFetching: usersFetching,
  } = useQuery(
    orpc.getAuthAdminUsers.queryOptions({
      input: { query: { cursor, limit: PAGE_SIZE } },
      placeholderData: keepPreviousData,
    }),
  );
  const { data: bannedData, isLoading: bannedLoading } = useQuery(
    orpc.getAuthAdminUsersStatesBanned.queryOptions(),
  );

  const users = usersData?.body?.users ?? [];
  const bannedMap = React.useMemo(() => {
    const map = new Map<string, AdminUserState>();
    bannedData?.body?.forEach((state) => {
      if (state.user_id) {
        map.set(state.user_id, state);
      }
    });
    return map;
  }, [bannedData]);

  const isLoading = usersLoading || bannedLoading;

  const handleNext = React.useCallback(() => {
    const nextCursor = usersData?.body?.next_cursor;
    if (!nextCursor) return;
    setCursor(nextCursor);
    setCursorHistory((prev) => [...prev, cursor]);
  }, [usersData, cursor]);

  const handlePrevious = React.useCallback(() => {
    if (cursorHistory.length === 0) return;
    const previousCursor = cursorHistory[cursorHistory.length - 1];
    setCursor(previousCursor);
    setCursorHistory((prev) => prev.slice(0, -1));
  }, [cursorHistory]);

  return (
    <>
      {isLoading ? (
        <div className="flex items-center justify-center py-8">Loading...</div>
      ) : (
        <AdminUsersProvider bannedMap={bannedMap}>
          <div className="mb-4 flex justify-end">
            <CreateUserSheet />
          </div>
          <div className="overflow-hidden rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Verified</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead>Updated At</TableHead>
                  <TableHead className="w-15" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium">{user.name}</TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>
                      <Badge
                        variant={user.email_verified ? 'default' : 'secondary'}
                      >
                        {user.email_verified ? 'Yes' : 'No'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {bannedMap.has(user.id ?? '') ? (
                        <span className="text-destructive font-medium">
                          Banned
                        </span>
                      ) : (
                        <span className="text-green-600 font-medium">
                          Active
                        </span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Time
                        date={
                          user.created_at ? new Date(user.created_at) : null
                        }
                      />
                    </TableCell>
                    <TableCell>
                      <Time
                        date={
                          user.updated_at ? new Date(user.updated_at) : null
                        }
                      />
                    </TableCell>
                    <TableCell>
                      <RowActions user={user} />
                    </TableCell>
                  </TableRow>
                ))}
                {users.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={7}
                      className="h-24 text-center text-muted-foreground"
                    >
                      No users found.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
          <div className="flex items-center justify-end gap-2 mt-4">
            <Button
              variant="outline"
              size="sm"
              onClick={handlePrevious}
              disabled={cursorHistory.length === 0 || usersFetching}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={handleNext}
              disabled={!usersData?.body?.next_cursor || usersFetching}
            >
              Next
            </Button>
          </div>
        </AdminUsersProvider>
      )}
    </>
  );
}
