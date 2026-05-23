import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { Ban, MoreHorizontal, Trash2, UserCog } from 'lucide-react';
import React from 'react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
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
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { User } from '@/lib/api/types.gen';
import {
  z_ban_user_request,
  z_start_impersonation_request,
} from '@/lib/api/zod.gen';
import { orpc, queryClient } from '@/lib/orpc';
import { useAdminUsers } from '@/routes/(auth)/admin/-components/context';

export function RowActions({ user }: { user: User }) {
  const userId = user.id ?? '';
  const { bannedMap } = useAdminUsers();
  const isBanned = bannedMap.has(userId);

  const [impersonateOpen, setImpersonateOpen] = React.useState(false);
  const [banOpen, setBanOpen] = React.useState(false);

  return (
    <>
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
          <DropdownMenuItem onClick={() => setImpersonateOpen(true)}>
            <UserCog className="size-4 mr-2" />
            Impersonate
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          {isBanned ? (
            <UnBanUserMenuItem userId={userId} />
          ) : (
            <DropdownMenuItem onClick={() => setBanOpen(true)}>
              <Ban className="size-4 mr-2" />
              Ban
            </DropdownMenuItem>
          )}
          <DropdownMenuSeparator />
          <DeleteUserMenuItem userId={userId} />
        </DropdownMenuContent>
      </DropdownMenu>
      <ImpersonateDialog
        userId={userId}
        open={impersonateOpen}
        onOpenChange={setImpersonateOpen}
      />
      <BanUserDialog userId={userId} open={banOpen} onOpenChange={setBanOpen} />
    </>
  );
}

const IMPERSONATION_MAX_MINUTES = 15;

const impersonateSchema = z_start_impersonation_request
  .omit({ target_user_id: true, expires_in_seconds: true })
  .extend({
    reason: z.string().min(1, 'Reason is required'),
    expires_in_minutes: z
      .number()
      .int()
      .positive('Must be a positive number')
      .max(
        IMPERSONATION_MAX_MINUTES,
        `Cannot exceed ${IMPERSONATION_MAX_MINUTES} minutes`,
      ),
  });

function ImpersonateDialog({
  userId,
  open,
  onOpenChange,
}: {
  userId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const router = useRouter();
  const mutation = useMutation(
    orpc.postAuthAdminImpersonations.mutationOptions({
      onSuccess: () => {
        toast.success('Impersonation started');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminImpersonations.key(),
        });
        void router.invalidate();
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      reason: '',
      expires_in_minutes: IMPERSONATION_MAX_MINUTES,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        body: {
          target_user_id: userId,
          reason: value.reason,
          expires_in_seconds: value.expires_in_minutes * 60,
        },
      });
    },
    validators: {
      onSubmit: impersonateSchema,
      onChange: impersonateSchema,
    },
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Impersonate User</DialogTitle>
          <DialogDescription>
            Start an impersonation session as this user. Provide a reason for
            the audit log.
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
                <field.Textarea placeholder="Why are you impersonating this user?" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="expires_in_minutes">
            {(field) => (
              <div className="space-y-2">
                <field.Label>
                  Duration (minutes, max {IMPERSONATION_MAX_MINUTES})
                </field.Label>
                <field.Input
                  type="number"
                  min={1}
                  max={IMPERSONATION_MAX_MINUTES}
                  placeholder={String(IMPERSONATION_MAX_MINUTES)}
                />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppForm>
            <form.Button
              className="w-full"
              loading={mutation.isPending}
              loadingText="Starting..."
            >
              Start Impersonation
            </form.Button>
          </form.AppForm>
        </form>
      </DialogContent>
    </Dialog>
  );
}

const banSchema = z_ban_user_request.omit({ banned_until: true }).extend({
  reason: z.string().min(1, 'Reason is required'),
  banned_until: z.union([z.date(), z.undefined()]),
});

function BanUserDialog({
  userId,
  open,
  onOpenChange,
}: {
  userId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const mutation = useMutation(
    orpc.postAuthAdminUsersByUserIdBan.mutationOptions({
      onSuccess: () => {
        toast.success('User banned successfully');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsersStatesBanned.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
        onOpenChange(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      reason: '',
      banned_until: undefined as Date | undefined,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { user_id: userId },
        body: {
          reason: value.reason,
          ...(value.banned_until
            ? { banned_until: value.banned_until.toISOString() }
            : {}),
        },
      });
    },
    validators: {
      onSubmit: banSchema,
      onChange: banSchema,
    },
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Ban User</DialogTitle>
          <DialogDescription>
            Ban this user. Leave the end date blank for a permanent ban.
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
                <field.Textarea placeholder="Why is this user being banned?" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppField name="banned_until">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Banned Until</field.Label>
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
              loadingText="Banning..."
            >
              Ban User
            </form.Button>
          </form.AppForm>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function UnBanUserMenuItem({ userId }: { userId: string }) {
  const mutation = useMutation(
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

  return (
    <DropdownMenuItem
      onClick={() =>
        void mutation.mutateAsync({
          path: { user_id: userId },
        })
      }
      disabled={mutation.isPending}
    >
      <Ban className="size-4 mr-2" />
      Unban
    </DropdownMenuItem>
  );
}

function DeleteUserMenuItem({ userId }: { userId: string }) {
  const mutation = useMutation(
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

  return (
    <DropdownMenuItem
      variant="destructive"
      onClick={() =>
        void mutation.mutateAsync({
          path: { user_id: userId },
        })
      }
      disabled={mutation.isPending}
    >
      <Trash2 className="size-4 mr-2" />
      Delete
    </DropdownMenuItem>
  );
}
