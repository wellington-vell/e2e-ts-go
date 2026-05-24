import { createFileRoute } from '@tanstack/react-router';

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PermissionsTable } from '@/routes/(auth)/admin/-components/permissions-table';
import { RevokedSessionsTable } from '@/routes/(auth)/admin/-components/revoked-sessions-table';
import { RolesTable } from '@/routes/(auth)/admin/-components/roles-table';
import { UsersTable } from '@/routes/(auth)/admin/-components/users-table';

export const Route = createFileRoute('/(auth)/admin/')({
  component: AdminPage,
});

function AdminPage() {
  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-4">
        <h1 className="text-2xl font-bold">Admin Panel</h1>
        <p className="text-muted-foreground">Manage users and permissions</p>
      </div>

      <Tabs defaultValue="users">
        <TabsList>
          <TabsTrigger value="users">Users</TabsTrigger>
          <TabsTrigger value="roles">Roles</TabsTrigger>
          <TabsTrigger value="permissions">Permissions</TabsTrigger>
          <TabsTrigger value="sessions">Sessions</TabsTrigger>
        </TabsList>
        <TabsContent value="users">
          <UsersTable />
        </TabsContent>
        <TabsContent value="roles">
          <RolesTable />
        </TabsContent>
        <TabsContent value="permissions">
          <PermissionsTable />
        </TabsContent>
        <TabsContent value="sessions">
          <RevokedSessionsTable />
        </TabsContent>
      </Tabs>
    </div>
  );
}
