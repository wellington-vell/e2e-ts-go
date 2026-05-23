import { useMutation, useQuery } from '@tanstack/react-query';
import { Plus } from 'lucide-react';
import { MoreHorizontal, Pencil, Trash2 } from 'lucide-react';
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
import type { Permission } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

export function PermissionsTable() {
  const { data, isLoading } = useQuery(
    orpc.getAuthAccessControlPermissions.queryOptions(),
  );

  const deletePermission = useMutation(
    orpc.deleteAuthAccessControlPermissionsByPermissionId.mutationOptions({
      onSuccess: () => {
        toast.success('Permission deleted successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlPermissions.key(),
        });
      },
    }),
  );

  const permissions = data?.body ?? [];
  const [editingPermission, setEditingPermission] =
    React.useState<Permission | null>(null);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">Loading...</div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <CreatePermissionSheet />
      </div>
      <div className="overflow-hidden rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Key</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>System</TableHead>
              <TableHead className="w-15" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {permissions.map((permission) => (
              <TableRow key={permission.id}>
                <TableCell className="font-medium">{permission.key}</TableCell>
                <TableCell>{permission.description}</TableCell>
                <TableCell>
                  {permission.is_system ? (
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
                    <DropdownMenuContent align="end" className="w-40">
                      <DropdownMenuItem
                        onClick={() => {
                          setEditingPermission(permission);
                        }}
                      >
                        <Pencil className="size-4 mr-2" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() => {
                          void deletePermission.mutateAsync({
                            path: { permission_id: permission.id ?? '' },
                          });
                        }}
                        disabled={
                          deletePermission.isPending || permission.is_system
                        }
                      >
                        <Trash2 className="size-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
            {permissions.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={4}
                  className="h-24 text-center text-muted-foreground"
                >
                  No permissions found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      {editingPermission && (
        <EditPermissionSheet
          permission={editingPermission}
          open={!!editingPermission}
          onOpenChange={(open) => !open && setEditingPermission(null)}
        />
      )}
    </div>
  );
}

export function CreatePermissionSheet() {
  const [open, setOpen] = React.useState(false);

  const schema = z.object({
    key: z.string().min(1, 'Key is required'),
    description: z.string(),
    is_system: z.boolean(),
  });

  const mutation = useMutation(
    orpc.postAuthAccessControlPermissions.mutationOptions({
      onSuccess: () => {
        toast.success('Permission created successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlPermissions.key(),
        });
        setOpen(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      key: '',
      description: '',
      is_system: false,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        body: {
          key: value.key,
          description: value.description,
          is_system: value.is_system,
        },
      });
    },
    validators: {
      onSubmit: schema,
    },
  });

  return (
    <>
      <Button onClick={() => setOpen(true)} size="sm">
        <Plus className="size-4 mr-2" />
        Create Permission
      </Button>
      <Sheet open={open} onOpenChange={setOpen}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Create Permission</SheetTitle>
            <SheetDescription>Create a new permission.</SheetDescription>
          </SheetHeader>
          <form
            onSubmit={async (e) => {
              e.preventDefault();
              e.stopPropagation();
              await form.handleSubmit();
            }}
            className="space-y-4 py-4"
          >
            <form.AppField name="key">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Key</field.Label>
                  <field.Input placeholder="e.g. users.read" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="description">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Description</field.Label>
                  <field.Input placeholder="Permission description" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="is_system">
              {(field) => (
                <div className="flex items-center gap-2">
                  <field.Checkbox />
                  <field.Label>System Permission</field.Label>
                </div>
              )}
            </form.AppField>
            <form.AppForm>
              <form.Button
                className="w-full"
                loading={mutation.isPending}
                loadingText="Creating..."
              >
                Create Permission
              </form.Button>
            </form.AppForm>
          </form>
        </SheetContent>
      </Sheet>
    </>
  );
}

export function EditPermissionSheet({
  permission,
  open,
  onOpenChange,
}: {
  permission: Permission;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const schema = z.object({
    description: z.string(),
  });

  const mutation = useMutation(
    orpc.patchAuthAccessControlPermissionsByPermissionId.mutationOptions({
      onSuccess: () => {
        toast.success('Permission updated successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAccessControlPermissions.key(),
        });
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      description: permission.description ?? '',
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { permission_id: permission.id ?? '' },
        body: {
          description: value.description,
        },
      });
    },
    validators: {
      onSubmit: schema,
    },
  });

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>Edit Permission</SheetTitle>
          <SheetDescription>Update permission details.</SheetDescription>
        </SheetHeader>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="space-y-4 py-4"
        >
          <div className="space-y-2">
            <label className="text-sm font-medium">Key</label>
            <p className="text-sm text-muted-foreground">{permission.key}</p>
          </div>
          <form.AppField name="description">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Description</field.Label>
                <field.Input />
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
              Update Permission
            </form.Button>
          </form.AppForm>
        </form>
      </SheetContent>
    </Sheet>
  );
}
