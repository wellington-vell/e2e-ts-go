import { useMutation, useQuery } from '@tanstack/react-query';
import { Check, Plus, ShieldAlert, Trash2, X } from 'lucide-react';
import React from 'react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { ScrollArea } from '@/components/ui/scroll-area';
import { SelectTrigger, SelectValue } from '@/components/ui/select';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Time } from '@/components/ui/time';
import type {
  Role,
  User,
  UserPermissionInfo,
  UserRoleInfo,
} from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

function invalidateUserAccess(userId: string) {
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAccessControlUsersByUserIdRoles.key(),
  });
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAccessControlUsersByUserIdPermissions.key(),
  });
  // Self-impersonation or self-edit may change session perms
  if (userId) {
    void queryClient.invalidateQueries({ queryKey: orpc.getAuthMe.key() });
  }
}

export function ManageUserRolesSheet({
  user,
  open,
  onOpenChange,
}: {
  user: User;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const userId = user.id ?? '';

  const { data: currentRolesData, isLoading: rolesLoading } = useQuery(
    orpc.getAuthAccessControlUsersByUserIdRoles.queryOptions({
      input: { path: { user_id: userId } },
      enabled: open,
    }),
  );

  const { data: allRolesData } = useQuery(
    orpc.getAuthAccessControlRoles.queryOptions({ enabled: open }),
  );

  const currentRoles: UserRoleInfo[] = React.useMemo(
    () => currentRolesData?.body ?? [],
    [currentRolesData],
  );
  const allRoles: Role[] = React.useMemo(
    () => allRolesData?.body ?? [],
    [allRolesData],
  );

  const removeRole = useMutation(
    orpc.deleteAuthAccessControlUsersByUserIdRolesByRoleId.mutationOptions({
      onSuccess: () => {
        toast.success('Role removed');
        invalidateUserAccess(userId);
      },
    }),
  );

  const [bulkOpen, setBulkOpen] = React.useState(false);

  const assignedRoleIds = React.useMemo(() => {
    const ids = new Set<string>();
    for (const r of currentRoles) {
      if (r.role_id) ids.add(r.role_id);
    }
    return ids;
  }, [currentRoles]);

  const availableForAssign = React.useMemo(
    () => allRoles.filter((r) => r.id && !assignedRoleIds.has(r.id)),
    [allRoles, assignedRoleIds],
  );

  return (
    <>
      <Sheet open={open} onOpenChange={onOpenChange}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Manage Roles</SheetTitle>
            <SheetDescription>
              Manage roles and inspect permissions for{' '}
              <strong>{user.name || user.email}</strong>
            </SheetDescription>
          </SheetHeader>
          <div className="grid flex-1 auto-rows-min gap-6 px-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold">Current roles</h3>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setBulkOpen(true)}
                  disabled={allRoles.length === 0}
                >
                  Bulk replace
                </Button>
              </div>
              {rolesLoading ? (
                <div className="text-sm text-muted-foreground">Loading...</div>
              ) : currentRoles.length === 0 ? (
                <div className="rounded-md border border-dashed p-4 text-center text-sm text-muted-foreground">
                  No roles assigned.
                </div>
              ) : (
                <div className="space-y-2">
                  {currentRoles.map((role) => (
                    <div
                      key={role.role_id}
                      className="flex items-center justify-between rounded-md border p-3"
                    >
                      <div className="space-y-0.5">
                        <p className="text-sm font-medium">{role.role_name}</p>
                        <p className="text-xs text-muted-foreground">
                          {role.role_description || (
                            <span className="italic">No description</span>
                          )}
                        </p>
                        <div className="flex gap-3 pt-1 text-xs text-muted-foreground">
                          <span>
                            Assigned{' '}
                            <Time
                              date={
                                role.assigned_at
                                  ? new Date(role.assigned_at)
                                  : null
                              }
                            />
                          </span>
                          {role.expires_at && (
                            <span>
                              Expires <Time date={new Date(role.expires_at)} />
                            </span>
                          )}
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="size-8 text-destructive"
                        onClick={() =>
                          void removeRole.mutateAsync({
                            path: {
                              user_id: userId,
                              role_id: role.role_id ?? '',
                            },
                          })
                        }
                        disabled={removeRole.isPending}
                      >
                        <Trash2 className="size-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <AssignRoleForm
              userId={userId}
              availableRoles={availableForAssign}
            />

            <EffectivePermissionsSection userId={userId} open={open} />

            <CheckPermissionsForm userId={userId} />
          </div>
        </SheetContent>
      </Sheet>
      {bulkOpen && (
        <BulkReplaceRolesDialog
          userId={userId}
          currentRoleIds={Array.from(assignedRoleIds)}
          allRoles={allRoles}
          onOpenChange={setBulkOpen}
        />
      )}
    </>
  );
}

const assignSchema = z.object({
  role_id: z.string().min(1, 'Role is required'),
  expires_at: z.union([z.date(), z.undefined()]),
});

function AssignRoleForm({
  userId,
  availableRoles,
}: {
  userId: string;
  availableRoles: Role[];
}) {
  const mutation = useMutation(
    orpc.postAuthAccessControlUsersByUserIdRoles.mutationOptions({
      onSuccess: () => {
        toast.success('Role assigned');
        invalidateUserAccess(userId);
        form.reset();
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      role_id: '',
      expires_at: undefined as Date | undefined,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { user_id: userId },
        body: {
          role_id: value.role_id,
          ...(value.expires_at
            ? { expires_at: value.expires_at.toISOString() }
            : {}),
        },
      });
    },
    validators: { onSubmit: assignSchema, onChange: assignSchema },
  });

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-semibold">Assign new role</h3>
      {availableRoles.length === 0 ? (
        <div className="text-sm text-muted-foreground">
          All roles are already assigned.
        </div>
      ) : (
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="space-y-3 rounded-md border p-3"
        >
          <form.AppField name="role_id">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Role</field.Label>
                <field.Select>
                  <SelectTrigger id="role_id">
                    <SelectValue placeholder="Select a role" />
                  </SelectTrigger>
                  <field.SelectContent>
                    {availableRoles.map((role) => (
                      <field.SelectItem key={role.id} value={role.id ?? ''}>
                        {role.name}
                        {role.description ? ` — ${role.description}` : ''}
                      </field.SelectItem>
                    ))}
                  </field.SelectContent>
                </field.Select>
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="expires_at">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Expires at (optional)</field.Label>
                <field.Calendar
                  disabled={{ before: new Date() }}
                  className="w-full"
                />
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
              <Plus className="size-4 mr-2" />
              Assign role
            </form.Button>
          </form.AppForm>
        </form>
      )}
    </div>
  );
}

function BulkReplaceRolesDialog({
  userId,
  currentRoleIds,
  allRoles,
  onOpenChange,
}: {
  userId: string;
  currentRoleIds: string[];
  allRoles: Role[];
  onOpenChange: (open: boolean) => void;
}) {
  const [selected, setSelected] = React.useState<Set<string>>(
    () => new Set(currentRoleIds),
  );

  const mutation = useMutation(
    orpc.putAuthAccessControlUsersByUserIdRoles.mutationOptions({
      onSuccess: () => {
        toast.success('Roles replaced');
        invalidateUserAccess(userId);
        onOpenChange(false);
      },
    }),
  );

  const toggle = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleConfirm = () => {
    void mutation.mutateAsync({
      path: { user_id: userId },
      body: { role_ids: Array.from(selected) },
    });
  };

  return (
    <AlertDialog open onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2">
            <ShieldAlert className="size-5 text-destructive" />
            Bulk replace roles
          </AlertDialogTitle>
          <AlertDialogDescription>
            This will overwrite the user's entire role set. Roles not selected
            here will be removed.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <div className="max-h-72 overflow-y-auto space-y-2 rounded-md border p-3">
          {allRoles.length === 0 ? (
            <div className="text-sm text-muted-foreground">
              No roles available.
            </div>
          ) : (
            allRoles.map((role) => {
              const id = role.id ?? '';
              const checked = selected.has(id);
              return (
                <label
                  key={id}
                  className="flex items-start gap-3 cursor-pointer"
                >
                  <Checkbox
                    checked={checked}
                    onCheckedChange={() => toggle(id)}
                  />
                  <div className="space-y-0.5">
                    <p className="text-sm font-medium">{role.name}</p>
                    {role.description && (
                      <p className="text-xs text-muted-foreground">
                        {role.description}
                      </p>
                    )}
                  </div>
                </label>
              );
            })
          )}
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={mutation.isPending}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={handleConfirm}
            disabled={mutation.isPending}
          >
            {mutation.isPending ? 'Replacing...' : 'Replace roles'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function EffectivePermissionsSection({
  userId,
  open,
}: {
  userId: string;
  open: boolean;
}) {
  const { data, isLoading } = useQuery(
    orpc.getAuthAccessControlUsersByUserIdPermissions.queryOptions({
      input: { path: { user_id: userId } },
      enabled: open,
    }),
  );

  const permissions: UserPermissionInfo[] = data?.body?.permissions ?? [];

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-semibold">Effective permissions</h3>
      {isLoading ? (
        <div className="text-sm text-muted-foreground">Loading...</div>
      ) : permissions.length === 0 ? (
        <div className="rounded-md border border-dashed p-4 text-center text-sm text-muted-foreground">
          No effective permissions.
        </div>
      ) : (
        <ScrollArea className="h-72 rounded-md border">
          <div className="space-y-2 p-2">
            {permissions.map((perm) => (
              <div
                key={perm.permission_id}
                className="rounded-md border p-3 space-y-1"
              >
                <div className="flex items-center justify-between gap-2">
                  <code className="text-sm font-medium">
                    {perm.permission_key}
                  </code>
                  {perm.sources && perm.sources.length > 0 && (
                    <div className="flex flex-wrap gap-1">
                      {perm.sources.map((src) => (
                        <Badge
                          key={`${perm.permission_id}-${src.role_id ?? src.granted_at ?? ''}`}
                          variant="secondary"
                          className="text-xs"
                        >
                          {src.role_name || src.role_id}
                        </Badge>
                      ))}
                    </div>
                  )}
                </div>
                {perm.permission_description && (
                  <p className="text-xs text-muted-foreground">
                    {perm.permission_description}
                  </p>
                )}
              </div>
            ))}
          </div>
        </ScrollArea>
      )}
    </div>
  );
}

const checkSchema = z.object({
  keys: z.string().min(1, 'Enter at least one permission key'),
});

function CheckPermissionsForm({ userId }: { userId: string }) {
  const [result, setResult] = React.useState<boolean | null>(null);

  const mutation = useMutation(
    orpc.postAuthAccessControlUsersByUserIdPermissionsCheck.mutationOptions({
      onSuccess: (resp) => {
        const has = resp?.body?.has_permissions ?? false;
        setResult(has);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: { keys: '' },
    onSubmit: async ({ value }) => {
      const parsed = value.keys
        .split(/[\n,]/)
        .map((s) => s.trim())
        .filter(Boolean);
      if (parsed.length === 0) {
        toast.error('Enter at least one permission key');
        return;
      }
      await mutation.mutateAsync({
        path: { user_id: userId },
        body: { permission_keys: parsed },
      });
    },
    validators: { onSubmit: checkSchema, onChange: checkSchema },
  });

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-semibold">Check permissions</h3>
      <form
        onSubmit={async (e) => {
          e.preventDefault();
          e.stopPropagation();
          setResult(null);
          await form.handleSubmit();
        }}
        className="space-y-3 rounded-md border p-3"
      >
        <form.AppField name="keys">
          {(field) => (
            <div className="space-y-2">
              <field.Label>
                Permission keys (one per line, or comma-separated)
              </field.Label>
              <field.Textarea
                placeholder={'users.read\nusers.write'}
                rows={3}
              />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <div className="flex items-center gap-3">
          <form.AppForm>
            <form.Button
              loading={mutation.isPending}
              loadingText="Checking..."
              size="sm"
            >
              Check
            </form.Button>
          </form.AppForm>
          {result !== null && !mutation.isPending && (
            <Badge variant={result ? 'default' : 'destructive'}>
              {result ? (
                <>
                  <Check className="size-3 mr-1" />
                  Has all permissions
                </>
              ) : (
                <>
                  <X className="size-3 mr-1" />
                  Missing one or more
                </>
              )}
            </Badge>
          )}
        </div>
      </form>
    </div>
  );
}
