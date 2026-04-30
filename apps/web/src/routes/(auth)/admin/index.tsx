import { useMutation, useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import React from 'react';
import { toast } from 'sonner';

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
import type { AdminUserState, BanUserRequest } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

export const Route = createFileRoute('/(auth)/admin/')({
  component: AdminPage,
});

function AdminPage() {
  const { data: usersData, isLoading: usersLoading } = useQuery(
    orpc.getAuthAdminUsers.queryOptions(),
  );
  const { data: bannedData, isLoading: bannedLoading } = useQuery(
    orpc.getAuthAdminUsersStatesBanned.queryOptions(),
  );

  const impersonate = useMutation(
    orpc.postAuthAdminImpersonations.mutationOptions({
      onSuccess: () => {
        toast.success('Impersonation started');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
      },
    }),
  );

  const banUser = useMutation(
    orpc.postAuthAdminUsersByUserIdBan.mutationOptions({
      onSuccess: () => {
        toast.success('User banned successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsersStatesBanned.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
      },
    }),
  );

  const unbanUser = useMutation(
    orpc.postAuthAdminUsersByUserIdUnban.mutationOptions({
      onSuccess: () => {
        toast.success('User unbanned successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsersStatesBanned.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
      },
    }),
  );

  const deleteUser = useMutation(
    orpc.deleteAuthAdminUsersByUserId.mutationOptions({
      onSuccess: () => {
        toast.success('User deleted successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsers.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
      },
    }),
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
                  {bannedMap.has(user.id ?? '') ? (
                    <span className="text-destructive">Banned</span>
                  ) : (
                    <span className="text-green-600">Active</span>
                  )}
                </TableCell>
                <TableCell className="flex gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={impersonate.isPending}
                    onClick={() =>
                      void impersonate.mutateAsync({
                        body: {
                          target_user_id: user.id ?? '',
                          reason: 'Admin action',
                        },
                      })
                    }
                  >
                    Impersonate
                  </Button>
                  {bannedMap.has(user.id ?? '') ? (
                    <Button
                      size="sm"
                      variant="secondary"
                      disabled={bannedLoading}
                      onClick={() =>
                        void unbanUser.mutateAsync({
                          path: { user_id: user.id ?? '' },
                        })
                      }
                    >
                      Unban
                    </Button>
                  ) : (
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={bannedLoading}
                      onClick={() => {
                        const params: BanUserRequest = {
                          reason: 'Admin ban',
                        };
                        void banUser.mutateAsync({
                          path: { user_id: user.id ?? '' },
                          body: params,
                        });
                      }}
                    >
                      Ban
                    </Button>
                  )}
                  <Button
                    size="sm"
                    variant="destructive"
                    disabled={deleteUser.isPending}
                    onClick={() =>
                      void deleteUser.mutateAsync({
                        path: { user_id: user.id ?? '' },
                      })
                    }
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

export function AdminPageV2() {
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
