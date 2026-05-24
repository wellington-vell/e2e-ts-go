import { useMutation, useQuery } from '@tanstack/react-query';
import { Trash2 } from 'lucide-react';
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
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import type { AdminSessionState } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

function invalidateSessionState() {
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminSessionsBySessionIdState.key(),
  });
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminUsersByUserIdSessions.key(),
  });
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminSessionsStatesRevoked.key(),
  });
}

export function SessionStateSheet({
  sessionId,
  open,
  onOpenChange,
}: {
  sessionId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { data, isLoading } = useQuery(
    orpc.getAuthAdminSessionsBySessionIdState.queryOptions({
      input: { path: { session_id: sessionId } },
      enabled: open,
      retry: false,
    }),
  );

  if (!open) return null;

  if (isLoading) {
    return (
      <Sheet open onOpenChange={onOpenChange}>
        <SheetContent className="sm:max-w-lg overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Session state</SheetTitle>
            <SheetDescription>
              Loading state for <code className="text-xs">{sessionId}</code>
            </SheetDescription>
          </SheetHeader>
          <div className="py-4 text-sm text-muted-foreground">Loading...</div>
        </SheetContent>
      </Sheet>
    );
  }

  const existing = data?.body?.state;

  if (existing) {
    return (
      <UpdateSessionStateSheet
        sessionId={sessionId}
        existing={existing}
        onOpenChange={onOpenChange}
      />
    );
  }

  return (
    <CreateSessionStateSheet
      sessionId={sessionId}
      onOpenChange={onOpenChange}
    />
  );
}

const createSchema = z.object({
  revoke: z.boolean(),
  revoked_reason: z.string(),
  impersonator_user_id: z.string(),
  impersonation_reason: z.string(),
  impersonation_expires_at: z.union([z.date(), z.undefined()]),
});

function CreateSessionStateSheet({
  sessionId,
  onOpenChange,
}: {
  sessionId: string;
  onOpenChange: (open: boolean) => void;
}) {
  const mutation = useMutation(
    orpc.postAuthAdminSessionsBySessionIdState.mutationOptions({
      onSuccess: () => {
        toast.success('Session state created');
        invalidateSessionState();
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      revoke: false,
      revoked_reason: '',
      impersonator_user_id: '',
      impersonation_reason: '',
      impersonation_expires_at: undefined as Date | undefined,
    },
    onSubmit: async ({ value }) => {
      const body: {
        revoke: boolean;
        revoked_reason?: string;
        impersonator_user_id?: string;
        impersonation_reason?: string;
        impersonation_expires_at?: string;
      } = { revoke: value.revoke };
      if (value.revoked_reason) body.revoked_reason = value.revoked_reason;
      if (value.impersonator_user_id) {
        body.impersonator_user_id = value.impersonator_user_id;
      }
      if (value.impersonation_reason) {
        body.impersonation_reason = value.impersonation_reason;
      }
      if (value.impersonation_expires_at) {
        body.impersonation_expires_at =
          value.impersonation_expires_at.toISOString();
      }
      await mutation.mutateAsync({
        path: { session_id: sessionId },
        body,
      });
    },
    validators: { onSubmit: createSchema, onChange: createSchema },
  });

  return (
    <Sheet open onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-lg overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Create session state</SheetTitle>
          <SheetDescription>
            Create the raw state record for session{' '}
            <code className="text-xs">{sessionId}</code>
          </SheetDescription>
        </SheetHeader>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="space-y-4 py-4"
        >
          <form.AppField name="revoke">
            {(field) => (
              <div className="flex items-center gap-2">
                <field.Checkbox />
                <field.Label>Revoked</field.Label>
              </div>
            )}
          </form.AppField>
          <form.AppField name="revoked_reason">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Revoked reason</field.Label>
                <field.Textarea placeholder="Why was this session revoked?" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="impersonator_user_id">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Impersonator user ID</field.Label>
                <field.Input placeholder="User ID acting as impersonator" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="impersonation_reason">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Impersonation reason</field.Label>
                <field.Textarea placeholder="Reason for impersonation" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="impersonation_expires_at">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Impersonation expires at</field.Label>
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
              loadingText="Creating..."
            >
              Create state
            </form.Button>
          </form.AppForm>
        </form>
      </SheetContent>
    </Sheet>
  );
}

const updateSchema = z.object({
  revoke: z.boolean(),
  revoked_reason: z.string(),
  impersonator_user_id: z.string(),
  impersonation_reason: z.string(),
  impersonation_expires_at: z.union([z.date(), z.undefined()]),
});

function UpdateSessionStateSheet({
  sessionId,
  existing,
  onOpenChange,
}: {
  sessionId: string;
  existing: AdminSessionState;
  onOpenChange: (open: boolean) => void;
}) {
  const update = useMutation(
    orpc.patchAuthAdminSessionsBySessionIdState.mutationOptions({
      onSuccess: () => {
        toast.success('Session state updated');
        invalidateSessionState();
        onOpenChange(false);
      },
    }),
  );

  const remove = useMutation(
    orpc.deleteAuthAdminSessionsBySessionIdState.mutationOptions({
      onSuccess: () => {
        toast.success('Session state deleted');
        invalidateSessionState();
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      revoke: !!existing.revoked_at,
      revoked_reason: existing.revoked_reason ?? '',
      impersonator_user_id: existing.impersonator_user_id ?? '',
      impersonation_reason: existing.impersonation_reason ?? '',
      impersonation_expires_at: existing.impersonation_expires_at
        ? new Date(existing.impersonation_expires_at)
        : undefined,
    },
    onSubmit: async ({ value }) => {
      const body: {
        revoke: boolean;
        revoked_reason?: string;
        impersonator_user_id?: string;
        impersonation_reason?: string;
        impersonation_expires_at?: string;
      } = { revoke: value.revoke };
      if (value.revoked_reason) body.revoked_reason = value.revoked_reason;
      if (value.impersonator_user_id) {
        body.impersonator_user_id = value.impersonator_user_id;
      }
      if (value.impersonation_reason) {
        body.impersonation_reason = value.impersonation_reason;
      }
      if (value.impersonation_expires_at) {
        body.impersonation_expires_at =
          value.impersonation_expires_at.toISOString();
      }
      await update.mutateAsync({
        path: { session_id: sessionId },
        body,
      });
    },
    validators: { onSubmit: updateSchema, onChange: updateSchema },
  });

  return (
    <Sheet open onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-lg overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Edit session state</SheetTitle>
          <SheetDescription>
            Update the raw state record for session{' '}
            <code className="text-xs">{sessionId}</code>
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
          <form.AppField name="revoke">
            {(field) => (
              <div className="flex items-center gap-2">
                <field.Checkbox />
                <field.Label>Revoked</field.Label>
              </div>
            )}
          </form.AppField>
          <form.AppField name="revoked_reason">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Revoked reason</field.Label>
                <field.Textarea placeholder="Why was this session revoked?" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="impersonator_user_id">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Impersonator user ID</field.Label>
                <field.Input placeholder="User ID acting as impersonator" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="impersonation_reason">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Impersonation reason</field.Label>
                <field.Textarea placeholder="Reason for impersonation" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="impersonation_expires_at">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Impersonation expires at</field.Label>
                <field.Calendar
                  disabled={{ before: new Date() }}
                  className="w-full"
                />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <div className="flex items-center gap-2 pt-2">
            <form.AppForm>
              <form.Button
                className="flex-1"
                loading={update.isPending}
                loadingText="Saving..."
              >
                Save changes
              </form.Button>
            </form.AppForm>
            <AlertDialog>
              <AlertDialogTrigger
                render={(props) => (
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    disabled={remove.isPending}
                    {...props}
                  >
                    <Trash2 className="size-4" />
                  </Button>
                )}
              />
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete session state?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This removes the state record for session{' '}
                    <code className="text-xs">{sessionId}</code>. The session
                    itself is not affected.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel disabled={remove.isPending}>
                    Cancel
                  </AlertDialogCancel>
                  <AlertDialogAction
                    onClick={() =>
                      void remove.mutateAsync({
                        path: { session_id: sessionId },
                      })
                    }
                    disabled={remove.isPending}
                  >
                    {remove.isPending ? 'Deleting...' : 'Delete'}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        </form>
      </SheetContent>
    </Sheet>
  );
}
