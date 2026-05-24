import { useMutation, useQuery } from '@tanstack/react-query';
import { KeyRound, MoreHorizontal, Pencil, Plus, Trash2 } from 'lucide-react';
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
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
import { Time } from '@/components/ui/time';
import type { Account, User } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

const KNOWN_PROVIDERS = ['email', 'google', 'github', 'discord'] as const;

function invalidateAccounts() {
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminUsersByUserIdAccounts.key(),
  });
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminAccountsById.key(),
  });
}

export function ManageUserAccountsSheet({
  user,
  open,
  onOpenChange,
}: {
  user: User;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const userId = user.id ?? '';

  const { data, isLoading } = useQuery(
    orpc.getAuthAdminUsersByUserIdAccounts.queryOptions({
      input: { path: { user_id: userId } },
      enabled: open,
    }),
  );

  const accounts: Account[] = React.useMemo(
    () => data?.body?.accounts ?? [],
    [data],
  );

  const [createOpen, setCreateOpen] = React.useState(false);
  const [editingId, setEditingId] = React.useState<string | null>(null);
  const [deleting, setDeleting] = React.useState<Account | null>(null);

  return (
    <>
      <Sheet open={open} onOpenChange={onOpenChange}>
        <SheetContent className="sm:max-w-xl overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Provider accounts</SheetTitle>
            <SheetDescription>
              Add, edit, or remove the credential sets linked to{' '}
              <strong>{user.name || user.email}</strong>. Existing tokens and
              passwords are never displayed.
            </SheetDescription>
          </SheetHeader>
          <div className="grid flex-1 auto-rows-min gap-6 px-4">
            <div className="flex justify-end">
              <Button size="sm" onClick={() => setCreateOpen(true)}>
                <Plus className="size-4 mr-2" />
                Add account
              </Button>
            </div>
            {isLoading ? (
              <div className="text-sm text-muted-foreground">Loading...</div>
            ) : accounts.length === 0 ? (
              <div className="rounded-md border border-dashed p-4 text-center text-sm text-muted-foreground">
                No accounts linked.
              </div>
            ) : (
              <div className="space-y-2">
                {accounts.map((account) => (
                  <AccountRow
                    key={account.id}
                    account={account}
                    onEdit={() => setEditingId(account.id ?? null)}
                    onDelete={() => setDeleting(account)}
                  />
                ))}
              </div>
            )}
          </div>
        </SheetContent>
      </Sheet>

      {createOpen && (
        <CreateAccountSheet
          userId={userId}
          onOpenChange={(o) => !o && setCreateOpen(false)}
        />
      )}

      {editingId && (
        <EditAccountSheet
          accountId={editingId}
          onOpenChange={(o) => !o && setEditingId(null)}
        />
      )}

      {deleting && (
        <DeleteAccountDialog
          account={deleting}
          onOpenChange={(o) => !o && setDeleting(null)}
        />
      )}
    </>
  );
}

function AccountRow({
  account,
  onEdit,
  onDelete,
}: {
  account: Account;
  onEdit: () => void;
  onDelete: () => void;
}) {
  return (
    <div className="rounded-md border p-3 space-y-2">
      <div className="flex items-start justify-between gap-2">
        <div className="space-y-0.5 min-w-0">
          <div className="flex items-center gap-2">
            <Badge variant="default">{account.provider_id ?? 'unknown'}</Badge>
            <code className="text-xs truncate">{account.account_id}</code>
          </div>
          <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
            <span>
              Created{' '}
              <Time
                date={account.created_at ? new Date(account.created_at) : null}
              />
            </span>
            <span>
              Updated{' '}
              <Time
                date={account.updated_at ? new Date(account.updated_at) : null}
              />
            </span>
            {account.access_token_expires_at && (
              <span>
                Access expires{' '}
                <Time date={new Date(account.access_token_expires_at)} />
              </span>
            )}
            {account.refresh_token_expires_at && (
              <span>
                Refresh expires{' '}
                <Time date={new Date(account.refresh_token_expires_at)} />
              </span>
            )}
          </div>
        </div>
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
          <DropdownMenuContent align="end" className="w-32">
            <DropdownMenuItem onClick={onEdit}>
              <Pencil className="size-4 mr-2" />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem variant="destructive" onClick={onDelete}>
              <Trash2 className="size-4 mr-2" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}

const createSchema = z.object({
  provider_id: z.string().min(1, 'Provider is required'),
  account_id: z.string().min(1, 'Account ID is required'),
  password: z.string(),
  access_token: z.string(),
  refresh_token: z.string(),
  id_token: z.string(),
  scope: z.string(),
  access_token_expires_at: z.union([z.date(), z.undefined()]),
  refresh_token_expires_at: z.union([z.date(), z.undefined()]),
});

function CreateAccountSheet({
  userId,
  onOpenChange,
}: {
  userId: string;
  onOpenChange: (open: boolean) => void;
}) {
  const mutation = useMutation(
    orpc.postAuthAdminUsersByUserIdAccounts.mutationOptions({
      onSuccess: () => {
        toast.success('Account added');
        invalidateAccounts();
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      provider_id: '',
      account_id: '',
      password: '',
      access_token: '',
      refresh_token: '',
      id_token: '',
      scope: '',
      access_token_expires_at: undefined as Date | undefined,
      refresh_token_expires_at: undefined as Date | undefined,
    },
    onSubmit: async ({ value }) => {
      const body: {
        provider_id?: string;
        account_id?: string;
        password?: string;
        access_token?: string;
        refresh_token?: string;
        id_token?: string;
        scope?: string;
        access_token_expires_at?: string;
        refresh_token_expires_at?: string;
      } = {
        provider_id: value.provider_id,
        account_id: value.account_id,
      };
      if (value.password) body.password = value.password;
      if (value.access_token) body.access_token = value.access_token;
      if (value.refresh_token) body.refresh_token = value.refresh_token;
      if (value.id_token) body.id_token = value.id_token;
      if (value.scope) body.scope = value.scope;
      if (value.access_token_expires_at) {
        body.access_token_expires_at =
          value.access_token_expires_at.toISOString();
      }
      if (value.refresh_token_expires_at) {
        body.refresh_token_expires_at =
          value.refresh_token_expires_at.toISOString();
      }
      await mutation.mutateAsync({
        path: { user_id: userId },
        body,
      });
    },
    validators: { onSubmit: createSchema },
  });

  return (
    <Sheet open onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-lg overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Add account</SheetTitle>
          <SheetDescription>
            Link a new credential set to this user.
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
          <form.AppField name="provider_id">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Provider</field.Label>
                <field.Select>
                  <SelectTrigger id="provider_id">
                    <SelectValue placeholder="Select a provider" />
                  </SelectTrigger>
                  <field.SelectContent>
                    {KNOWN_PROVIDERS.map((p) => (
                      <field.SelectItem key={p} value={p}>
                        {p}
                      </field.SelectItem>
                    ))}
                  </field.SelectContent>
                </field.Select>
                <field.Error />
              </div>
            )}
          </form.AppField>

          <form.AppField name="account_id">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Account ID</field.Label>
                <field.Input placeholder="email@example.com or upstream user id" />
                <field.Error />
              </div>
            )}
          </form.AppField>

          <div className="rounded-md border p-3 space-y-3">
            <p className="text-xs text-muted-foreground">
              Credentials (optional).
            </p>
            <form.AppField name="password">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Password</field.Label>
                  <field.Input type="password" autoComplete="new-password" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="access_token">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Access token</field.Label>
                  <field.Input type="password" autoComplete="off" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="refresh_token">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Refresh token</field.Label>
                  <field.Input type="password" autoComplete="off" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="id_token">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>ID token</field.Label>
                  <field.Input type="password" autoComplete="off" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="scope">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Scope</field.Label>
                  <field.Input placeholder="openid email profile" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="access_token_expires_at">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Access token expires at</field.Label>
                  <field.Calendar className="w-full" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="refresh_token_expires_at">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Refresh token expires at</field.Label>
                  <field.Calendar className="w-full" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
          </div>

          <form.AppForm>
            <form.Button
              className="w-full"
              loading={mutation.isPending}
              loadingText="Adding..."
            >
              <KeyRound className="size-4 mr-2" />
              Add account
            </form.Button>
          </form.AppForm>
        </form>
      </SheetContent>
    </Sheet>
  );
}

const editSchema = z.object({
  account_id: z.string().min(1, 'Account ID is required'),
  password: z.string(),
  access_token: z.string(),
  refresh_token: z.string(),
  id_token: z.string(),
  scope: z.string(),
  access_token_expires_at: z.union([z.date(), z.undefined()]),
  refresh_token_expires_at: z.union([z.date(), z.undefined()]),
});

function EditAccountSheet({
  accountId,
  onOpenChange,
}: {
  accountId: string;
  onOpenChange: (open: boolean) => void;
}) {
  const { data, isLoading } = useQuery(
    orpc.getAuthAdminAccountsById.queryOptions({
      input: { path: { id: accountId } },
      enabled: true,
      meta: { silent: true },
    }),
  );

  const existing = data?.body?.account;

  return (
    <Sheet open onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-lg overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Edit account</SheetTitle>
          <SheetDescription>
            Update provider binding. Credential fields are write-only — leave
            them blank to keep existing values.
          </SheetDescription>
        </SheetHeader>
        <div className="grid flex-1 auto-rows-min gap-6 px-4">
          {isLoading || !existing ? (
            <div className="text-sm text-muted-foreground">Loading...</div>
          ) : (
            <EditAccountForm
              accountId={accountId}
              existing={existing}
              onDone={() => onOpenChange(false)}
            />
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}

function EditAccountForm({
  accountId,
  existing,
  onDone,
}: {
  accountId: string;
  existing: Account;
  onDone: () => void;
}) {
  const mutation = useMutation(
    orpc.patchAuthAdminAccountsById.mutationOptions({
      onSuccess: () => {
        toast.success('Account updated');
        invalidateAccounts();
        onDone();
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      account_id: existing.account_id ?? '',
      password: '',
      access_token: '',
      refresh_token: '',
      id_token: '',
      scope: existing.scope ?? '',
      access_token_expires_at: existing.access_token_expires_at
        ? new Date(existing.access_token_expires_at)
        : undefined,
      refresh_token_expires_at: existing.refresh_token_expires_at
        ? new Date(existing.refresh_token_expires_at)
        : undefined,
    },
    onSubmit: async ({ value }) => {
      const body: {
        account_id?: string;
        password?: string;
        access_token?: string;
        refresh_token?: string;
        id_token?: string;
        scope?: string;
        access_token_expires_at?: string;
        refresh_token_expires_at?: string;
      } = {
        account_id: value.account_id,
      };
      if (value.password) body.password = value.password;
      if (value.access_token) body.access_token = value.access_token;
      if (value.refresh_token) body.refresh_token = value.refresh_token;
      if (value.id_token) body.id_token = value.id_token;
      if (value.scope) body.scope = value.scope;
      if (value.access_token_expires_at) {
        body.access_token_expires_at =
          value.access_token_expires_at.toISOString();
      }
      if (value.refresh_token_expires_at) {
        body.refresh_token_expires_at =
          value.refresh_token_expires_at.toISOString();
      }
      await mutation.mutateAsync({
        path: { id: accountId },
        body,
      });
    },
    validators: { onSubmit: editSchema },
  });

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault();
        e.stopPropagation();
        await form.handleSubmit();
      }}
      className="space-y-4"
    >
      <div className="space-y-2">
        <span className="text-sm font-medium">Provider</span>
        <div className="flex h-9 items-center rounded-md border bg-muted px-3 text-sm">
          {existing.provider_id ?? '—'}
        </div>
      </div>

      <form.AppField name="account_id">
        {(field) => (
          <div className="space-y-2">
            <field.Label>Account ID</field.Label>
            <field.Input />
            <field.Error />
          </div>
        )}
      </form.AppField>

      <div className="rounded-md border p-3 space-y-3">
        <p className="text-xs text-muted-foreground">
          Credential overrides. Leave blank to keep existing values.
        </p>
        <form.AppField name="password">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Password</field.Label>
              <field.Input type="password" autoComplete="new-password" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="access_token">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Access token</field.Label>
              <field.Input type="password" autoComplete="off" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="refresh_token">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Refresh token</field.Label>
              <field.Input type="password" autoComplete="off" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="id_token">
          {(field) => (
            <div className="space-y-2">
              <field.Label>ID token</field.Label>
              <field.Input type="password" autoComplete="off" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="scope">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Scope</field.Label>
              <field.Input placeholder="openid email profile" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="access_token_expires_at">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Access token expires at</field.Label>
              <field.Calendar className="w-full" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="refresh_token_expires_at">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Refresh token expires at</field.Label>
              <field.Calendar className="w-full" />
              <field.Error />
            </div>
          )}
        </form.AppField>
      </div>

      <form.AppForm>
        <form.Button
          className="w-full"
          loading={mutation.isPending}
          loadingText="Saving..."
        >
          <KeyRound className="size-4 mr-2" />
          Save changes
        </form.Button>
      </form.AppForm>
    </form>
  );
}

function DeleteAccountDialog({
  account,
  onOpenChange,
}: {
  account: Account;
  onOpenChange: (open: boolean) => void;
}) {
  const mutation = useMutation(
    orpc.deleteAuthAdminAccountsById.mutationOptions({
      onSuccess: () => {
        toast.success('Account deleted');
        invalidateAccounts();
        onOpenChange(false);
      },
    }),
  );

  return (
    <AlertDialog open onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete account?</AlertDialogTitle>
          <AlertDialogDescription>
            The user will no longer be able to sign in with{' '}
            <strong>{account.provider_id}</strong> (
            <code className="text-xs">{account.account_id}</code>). This cannot
            be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={mutation.isPending}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={() =>
              void mutation.mutateAsync({
                path: { id: account.id ?? '' },
              })
            }
            disabled={mutation.isPending}
          >
            {mutation.isPending ? 'Deleting...' : 'Delete'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
