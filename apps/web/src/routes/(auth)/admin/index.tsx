import { createFileRoute } from '@tanstack/react-router';
import type { AdminUserState, BanUserRequest } from 'authula/plugins';
import React from 'react';

import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useAuth } from '@/context/auth';
import { useAdminUsers, useBannedUsers } from '@/context/auth';

export const Route = createFileRoute('/(auth)/admin/')({
  component: AdminPage,
});

function AdminPage() {
  const { banUser, unbanUser, deleteUser, impersonate } = useAuth();

  const { data: usersData, isLoading: usersLoading } = useAdminUsers();
  const { data: bannedData, isLoading: bannedLoading } = useBannedUsers();

  const users = usersData?.users ?? [];
  const bannedMap = React.useMemo(() => {
    const map = new Map<string, AdminUserState>();
    bannedData?.forEach((state) => {
      map.set(state.userId, state);
    });
    return map;
  }, [bannedData]);

  const isLoading = usersLoading || bannedLoading;

  const isActionLoading = (userId: string) =>
    (impersonate.isPending && impersonate.variables?.targetUserId === userId) ||
    (banUser.isPending && banUser.variables?.userId === userId) ||
    (unbanUser.isPending && unbanUser.variables === userId) ||
    (deleteUser.isPending && deleteUser.variables === userId);

  return (
    <div className="p-6">
      <div className="mb-4">
        <h1 className="text-2xl font-bold">Admin Panel</h1>
        <p className="text-muted-foreground">Manage users and permissions</p>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-8">Loading...</div>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  {bannedMap.has(user.id) ? (
                    <span className="text-destructive">Banned</span>
                  ) : (
                    <span className="text-green-600">Active</span>
                  )}
                </TableCell>
                <TableCell className="flex gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={isActionLoading(user.id)}
                    onClick={() =>
                      void impersonate.mutateAsync({
                        targetUserId: user.id,
                        reason: 'Admin action',
                      })
                    }
                  >
                    Impersonate
                  </Button>
                  {bannedMap.has(user.id) ? (
                    <Button
                      size="sm"
                      variant="secondary"
                      disabled={isActionLoading(user.id)}
                      onClick={() => void unbanUser.mutateAsync(user.id)}
                    >
                      Unban
                    </Button>
                  ) : (
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={isActionLoading(user.id)}
                      onClick={() => {
                        const params: BanUserRequest = { reason: 'Admin ban' };
                        void banUser.mutateAsync({ userId: user.id, params });
                      }}
                    >
                      Ban
                    </Button>
                  )}
                  <Button
                    size="sm"
                    variant="destructive"
                    disabled={isActionLoading(user.id)}
                    onClick={() => void deleteUser.mutateAsync(user.id)}
                  >
                    Delete
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  );
}

function AdminPageV2() {
  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex flex-col gap-1.5">
        <h1 className="text-3xl font-semibold tracking-tight">
          User Management
        </h1>
        <p className="text-muted-foreground text-lg">
          Manage user roles, bans, sessions, and impersonation.
        </p>
      </div>
      <Card className="shadow-sm rounded-xl p-4">
        <CardContent className="p-0">{/* <UsersTable /> */}</CardContent>
      </Card>
    </div>
  );
}
