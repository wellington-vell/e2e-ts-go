import { useMutation, useQuery } from '@tanstack/react-query';
import { MoreHorizontal, Pencil, Shield, Trash2 } from 'lucide-react';
import { Plus } from 'lucide-react';
import React from 'react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { SelectTrigger, SelectValue } from '@/components/ui/select';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Role } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

export function RolesTable() {
  const { data, isLoading } = useQuery(
    orpc.getAuthAccessControlRoles.queryOptions(),
  );

  const deleteRole = useMutation(
    orpc.deleteAuthAccessControlRolesByRoleId.mutationOptions({
      onSuccess: () => {
        toast.success('Role deleted successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlRoles.key(),
        });
      },
    }),
  );

  const roles = data?.body ?? [];

  const [editingRole, setEditingRole] = React.useState<Role | null>(null);
  const [permissionsRole, setPermissionsRole] = React.useState<Role | null>(
    null,
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">Loading...</div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <CreateRoleSheet />
      </div>
      <div className="overflow-hidden rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Weight</TableHead>
              <TableHead>System</TableHead>
              <TableHead className="w-15" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles.map((role) => (
              <TableRow key={role.id}>
                <TableCell className="font-medium">{role.name}</TableCell>
                <TableCell>{role.description}</TableCell>
                <TableCell>{role.weight}</TableCell>
                <TableCell>
                  {role.is_system ? (
                    <Badge variant="default">System</Badge>
                  ) : (
                    <Badge variant="secondary">Custom</Badge>
                  )}
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      render={(props) => (
                        <Button
                          variant="ghost"
                          size="icon"
                          className="size-8 data-popup-open:bg-muted"
                          {...props}
                        >
                          <MoreHorizontal />
                          <span className="sr-only">Open menu</span>
                        </Button>
                      )}
                    />
                    <DropdownMenuContent align="end" className="w-48">
                      <DropdownMenuItem
                        onClick={() => {
                          setEditingRole(role);
                        }}
                      >
                        <Pencil className="size-4 mr-2" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => {
                          setPermissionsRole(role);
                        }}
                      >
                        <Shield className="size-4 mr-2" />
                        Permissions
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() => {
                          void deleteRole.mutateAsync({
                            path: { role_id: role.id ?? '' },
                          });
                        }}
                        disabled={deleteRole.isPending || role.is_system}
                      >
                        <Trash2 className="size-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
            {roles.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="h-24 text-center text-muted-foreground"
                >
                  No roles found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      {editingRole && (
        <EditRoleSheet
          role={editingRole}
          open={!!editingRole}
          onOpenChange={(open) => !open && setEditingRole(null)}
        />
      )}
      {permissionsRole && (
        <RolePermissionsSheet
          role={permissionsRole}
          open={!!permissionsRole}
          onOpenChange={(open) => !open && setPermissionsRole(null)}
        />
      )}
    </div>
  );
}

export function CreateRoleSheet() {
  const [open, setOpen] = React.useState(false);

  const schema = z.object({
    name: z.string().min(1, 'Name is required'),
    description: z.string(),
    weight: z
      .number()
      .min(0, 'Weight must be a positive number')
      .max(100, 'Weight must be less than or equal to 100'),
    is_system: z.boolean(),
  });

  const mutation = useMutation(
    orpc.postAuthAccessControlRoles.mutationOptions({
      onSuccess: () => {
        toast.success('Role created successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlRoles.key(),
        });
        setOpen(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      name: '',
      description: '',
      weight: 0,
      is_system: false,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        body: {
          name: value.name,
          description: value.description,
          weight: value.weight,
          is_system: value.is_system,
        },
      });
    },
    validators: {
      onSubmit: schema,
      onChange: schema,
    },
  });

  return (
    <>
      <Button onClick={() => setOpen(true)} size="sm">
        <Plus className="size-4 mr-2" />
        Create Role
      </Button>
      <Sheet open={open} onOpenChange={setOpen}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Create Role</SheetTitle>
            <SheetDescription>
              Create a new access control role.
            </SheetDescription>
          </SheetHeader>
          <form
            onSubmit={async (e) => {
              e.preventDefault();
              e.stopPropagation();
              await form.handleSubmit();
            }}
            className="grid flex-1 auto-rows-min gap-6 px-4"
          >
            <form.AppField name="name">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Name</field.Label>
                  <field.Input placeholder="e.g. admin" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="description">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Description</field.Label>
                  <field.Input placeholder="Role description" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="weight">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Weight</field.Label>
                  <field.Input type="number" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="is_system">
              {(field) => (
                <div className="flex items-center gap-2">
                  <field.Checkbox />
                  <field.Label>System Role</field.Label>
                </div>
              )}
            </form.AppField>
            <form.AppForm>
              <form.Button
                className="w-full"
                loading={mutation.isPending}
                loadingText="Creating..."
              >
                Create Role
              </form.Button>
            </form.AppForm>
          </form>
        </SheetContent>
      </Sheet>
    </>
  );
}

export function EditRoleSheet({
  role,
  open,
  onOpenChange,
}: {
  role: Role;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const schema = z.object({
    name: z.string().min(1, 'Name is required'),
    description: z.string(),
    weight: z
      .number()
      .min(0, 'Weight must be a positive number')
      .max(100, 'Weight must be less than or equal to 100'),
  });

  const mutation = useMutation(
    orpc.patchAuthAccessControlRolesByRoleId.mutationOptions({
      onSuccess: () => {
        toast.success('Role updated successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlRoles.key(),
        });
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      name: role.name ?? '',
      description: role.description ?? '',
      weight: role.weight ?? 0,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { role_id: role.id ?? '' },
        body: {
          name: value.name,
          description: value.description,
          weight: value.weight,
        },
      });
    },
    validators: {
      onSubmit: schema,
      onChange: schema,
    },
  });

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>Edit Role</SheetTitle>
          <SheetDescription>Update role details.</SheetDescription>
        </SheetHeader>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="grid flex-1 auto-rows-min gap-6 px-4"
        >
          <form.AppField name="name">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Name</field.Label>
                <field.Input />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="description">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Description</field.Label>
                <field.Input />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="weight">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Weight</field.Label>
                <field.Input type="number" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppForm>
            <form.Button
              className="w-full"
              loading={mutation.isPending}
              loadingText="Updating..."
            >
              Update Role
            </form.Button>
          </form.AppForm>
        </form>
      </SheetContent>
    </Sheet>
  );
}

export function RolePermissionsSheet({
  role,
  open,
  onOpenChange,
}: {
  role: Role;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { data, isLoading } = useQuery(
    orpc.getAuthAccessControlRolesByRoleIdPermissions.queryOptions({
      input: { path: { role_id: role.id ?? '' } },
      enabled: open,
    }),
  );

  const removePermission = useMutation(
    orpc.deleteAuthAccessControlRolesByRoleIdPermissionsByPermissionId.mutationOptions(
      {
        onSuccess: () => {
          toast.success('Permission removed from role');
          void queryClient.invalidateQueries({
            queryKey: orpc.getAuthAccessControlRolesByRoleIdPermissions.key(),
          });
        },
      },
    ),
  );

  const permissions = data?.body ?? [];
  const [assignOpen, setAssignOpen] = React.useState(false);

  return (
    <>
      <Sheet open={open} onOpenChange={onOpenChange}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Role Permissions</SheetTitle>
            <SheetDescription>
              Manage permissions for <strong>{role.name}</strong>
            </SheetDescription>
          </SheetHeader>
          <div className="grid flex-1 auto-rows-min gap-6 px-4">
            <div className="flex justify-end">
              <Button size="sm" onClick={() => setAssignOpen(true)}>
                <Plus className="size-4 mr-2" />
                Add Permission
              </Button>
            </div>
            {isLoading ? (
              <div className="flex items-center justify-center py-8">
                Loading...
              </div>
            ) : (
              <div className="space-y-2">
                {permissions.map((permission) => (
                  <div
                    key={permission.permission_id}
                    className="flex items-center justify-between rounded-md border p-3"
                  >
                    <div className="space-y-0.5">
                      <p className="text-sm font-medium">
                        {permission.permission_key}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {permission.permission_description}
                      </p>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="size-8 text-destructive"
                      onClick={() =>
                        void removePermission.mutateAsync({
                          path: {
                            role_id: role.id ?? '',
                            permission_id: permission.permission_id ?? '',
                          },
                        })
                      }
                      disabled={removePermission.isPending}
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  </div>
                ))}
                {permissions.length === 0 && (
                  <div className="text-center text-muted-foreground py-8">
                    No permissions assigned to this role.
                  </div>
                )}
              </div>
            )}
          </div>
        </SheetContent>
      </Sheet>
      <AssignPermissionSheet
        role={role}
        open={assignOpen}
        onOpenChange={setAssignOpen}
      />
    </>
  );
}

const schema = z.object({
  permission_id: z.string().min(1, 'Permission is required'),
});

export function AssignPermissionSheet({
  role,
  open,
  onOpenChange,
}: {
  role: Role;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { data } = useQuery(
    orpc.getAuthAccessControlPermissions.queryOptions(),
  );

  const mutation = useMutation(
    orpc.postAuthAccessControlRolesByRoleIdPermissions.mutationOptions({
      onSuccess: () => {
        toast.success('Permission assigned successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlRolesByRoleIdPermissions.key(),
        });
        onOpenChange(false);
      },
    }),
  );

  const permissions = data?.body ?? [];

  const form = useAppForm({
    defaultValues: {
      permission_id: '',
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { role_id: role.id ?? '' },
        body: {
          permission_id: value.permission_id,
        },
      });
    },
    validators: {
      onSubmit: schema,
      onChange: schema,
    },
  });

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>Assign Permission</SheetTitle>
          <SheetDescription>
            Add a permission to <strong>{role.name}</strong>
          </SheetDescription>
        </SheetHeader>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="grid flex-1 auto-rows-min gap-6 px-4"
        >
          <form.AppField name="permission_id">
            {(field) => (
              <div className="space-y-2 w-full">
                <field.Label>Permission</field.Label>
                <field.Select>
                  <SelectTrigger id="permission_id">
                    <SelectValue placeholder="Select a permission" />
                  </SelectTrigger>
                  <field.SelectContent>
                    {permissions.map((permission) => (
                      <field.SelectItem
                        key={permission.id}
                        value={permission.id ?? ''}
                      >
                        {permission.description}
                      </field.SelectItem>
                    ))}
                  </field.SelectContent>
                </field.Select>
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppForm>
            <form.Button
              className="w-full"
              loading={mutation.isPending}
              loadingText="Assigning..."
            >
              Assign Permission
            </form.Button>
          </form.AppForm>
        </form>
      </SheetContent>
    </Sheet>
  );
}
