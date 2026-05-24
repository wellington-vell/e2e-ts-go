import { useMutation, useQuery } from '@tanstack/react-query';
import { MoreHorizontal, Pencil, ShieldOff } from 'lucide-react';
import React from 'react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Time } from '@/components/ui/time';
import type { AdminUserSession, User } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';
import { SessionStateSheet } from '@/routes/(auth)/admin/-components/session-state-sheet';

function invalidateSessions() {
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminUsersByUserIdSessions.key(),
  });
  void queryClient.invalidateQueries({
    queryKey: orpc.getAuthAdminSessionsStatesRevoked.key(),
  });
}

export function UserSessionsSheet({
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
    orpc.getAuthAdminUsersByUserIdSessions.queryOptions({
      input: { path: { user_id: userId } },
      enabled: open,
    }),
  );

  const sessions: AdminUserSession[] = React.useMemo(
    () => data?.body ?? [],
    [data],
  );

  const [revokeSessionId, setRevokeSessionId] = React.useState<string | null>(
    null,
  );
  const [stateSessionId, setStateSessionId] = React.useState<string | null>(
    null,
  );

  return (
    <>
      <Sheet open={open} onOpenChange={onOpenChange}>
        <SheetContent className="sm:max-w-xl overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Sessions</SheetTitle>
            <SheetDescription>
              Active and historical sessions for{' '}
              <strong>{user.name || user.email}</strong>
            </SheetDescription>
          </SheetHeader>
          <div className="grid flex-1 auto-rows-min gap-6 px-4">
            {isLoading ? (
              <div className="text-sm text-muted-foreground">Loading...</div>
            ) : sessions.length === 0 ? (
              <div className="rounded-md border border-dashed p-4 text-center text-sm text-muted-foreground">
                No sessions found.
              </div>
            ) : (
              <div className="space-y-2">
                {sessions.map((s) => (
                  <SessionRow
                    key={s.session?.id}
                    session={s}
                    onRevoke={() => setRevokeSessionId(s.session?.id ?? null)}
                    onEditState={() => setStateSessionId(s.session?.id ?? null)}
                  />
                ))}
              </div>
            )}
          </div>
        </SheetContent>
      </Sheet>
      {revokeSessionId && (
        <RevokeSessionDialog
          sessionId={revokeSessionId}
          onOpenChange={(o) => !o && setRevokeSessionId(null)}
        />
      )}
      {stateSessionId && (
        <SessionStateSheet
          sessionId={stateSessionId}
          open
          onOpenChange={(o) => !o && setStateSessionId(null)}
        />
      )}
    </>
  );
}

function SessionRow({
  session,
  onRevoke,
  onEditState,
}: {
  session: AdminUserSession;
  onRevoke: () => void;
  onEditState: () => void;
}) {
  const s = session.session;
  const revoked = !!session.state?.revoked_at;

  return (
    <div className="rounded-md border p-3 space-y-2">
      <div className="flex items-start justify-between gap-2">
        <div className="space-y-0.5 min-w-0">
          <code className="text-xs text-muted-foreground block truncate">
            {s?.id}
          </code>
          <div className="flex flex-wrap gap-2 text-xs">
            <span>
              Created{' '}
              <Time date={s?.created_at ? new Date(s.created_at) : null} />
            </span>
            {s?.expires_at && (
              <span>
                Expires <Time date={new Date(s.expires_at)} />
              </span>
            )}
          </div>
          {(s?.ip_address || s?.user_agent) && (
            <div className="text-xs text-muted-foreground truncate wrap-break-word max-w-60">
              {s.ip_address}
              {s.ip_address && s.user_agent ? ' · ' : ''}
              {s.user_agent}
            </div>
          )}
        </div>
        <div className="flex flex-col items-center gap-2">
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
              <DropdownMenuItem onClick={onEditState}>
                <Pencil className="size-4 mr-2" />
                Edit state
              </DropdownMenuItem>
              {!revoked && (
                <DropdownMenuItem variant="destructive" onClick={onRevoke}>
                  <ShieldOff className="size-4 mr-2" />
                  Revoke
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
          {revoked ? (
            <Badge variant="destructive">Revoked</Badge>
          ) : (
            <Badge variant="default">Active</Badge>
          )}
        </div>
      </div>
      {revoked && session.state?.revoked_reason && (
        <p className="text-xs text-muted-foreground">
          <span className="font-medium">Reason:</span>{' '}
          {session.state.revoked_reason}
        </p>
      )}
    </div>
  );
}

const revokeSchema = z.object({
  reason: z.string().min(1, 'Reason is required'),
});

function RevokeSessionDialog({
  sessionId,
  onOpenChange,
}: {
  sessionId: string;
  onOpenChange: (open: boolean) => void;
}) {
  const mutation = useMutation(
    orpc.postAuthAdminSessionsBySessionIdRevoke.mutationOptions({
      onSuccess: () => {
        toast.success('Session revoked');
        invalidateSessions();
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: { reason: '' },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { session_id: sessionId },
        body: { reason: value.reason },
      });
    },
    validators: { onSubmit: revokeSchema, onChange: revokeSchema },
  });

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Revoke session</DialogTitle>
          <DialogDescription>
            Revoke session <code className="text-xs">{sessionId}</code>. Provide
            a reason for the audit log.
          </DialogDescription>
        </DialogHeader>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="space-y-4"
        >
          <form.AppField name="reason">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Reason</field.Label>
                <field.Textarea placeholder="Why are you revoking this session?" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppForm>
            <form.Button
              className="w-full"
              loading={mutation.isPending}
              loadingText="Revoking..."
            >
              Revoke session
            </form.Button>
          </form.AppForm>
        </form>
      </DialogContent>
    </Dialog>
  );
}
