import { useMutation, useQuery } from '@tanstack/react-query';
import { ShieldAlert } from 'lucide-react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { orpc, queryClient } from '@/lib/orpc';

const schema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.email('Invalid email address'),
  image: z.string(),
  email_verified: z.boolean(),
});

export function EditUserSheet({
  userId,
  open,
  onOpenChange,
}: {
  userId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { data: meData } = useQuery(orpc.getAuthMe.queryOptions());
  const isSelf = !!meData?.body?.user?.id && meData.body.user.id === userId;

  const { data, isLoading } = useQuery(
    orpc.getAuthAdminUsersByUserId.queryOptions({
      input: { path: { user_id: userId } },
      enabled: open,
    }),
  );

  const user = data?.body?.user;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-lg overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Edit user</SheetTitle>
          <SheetDescription>
            Update name, email, or verification status.
          </SheetDescription>
        </SheetHeader>
        <div className="grid flex-1 auto-rows-min gap-6 px-4">
          {isSelf && (
            <Alert variant="destructive">
              <ShieldAlert />
              <AlertTitle>This is your own account</AlertTitle>
              <AlertDescription>
                Open <strong>Settings</strong> to change your own email or
                password.
              </AlertDescription>
            </Alert>
          )}
          {isLoading || !user ? (
            <div className="text-sm text-muted-foreground">Loading...</div>
          ) : (
            <EditUserForm
              userId={userId}
              initial={{
                name: user.name ?? '',
                email: user.email ?? '',
                image: user.image ?? '',
                email_verified: !!user.email_verified,
              }}
              disabled={isSelf}
              onDone={() => onOpenChange(false)}
            />
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}

function EditUserForm({
  userId,
  initial,
  disabled,
  onDone,
}: {
  userId: string;
  initial: {
    name: string;
    email: string;
    image: string;
    email_verified: boolean;
  };
  disabled: boolean;
  onDone: () => void;
}) {
  const mutation = useMutation(
    orpc.patchAuthAdminUsersByUserId.mutationOptions({
      onSuccess: () => {
        toast.success('User updated');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsers.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsersByUserId.key(),
        });
        onDone();
      },
    }),
  );

  const form = useAppForm({
    defaultValues: initial,
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        path: { user_id: userId },
        body: {
          name: value.name,
          email: value.email,
          email_verified: value.email_verified,
          ...(value.image ? { image: value.image } : {}),
          metadata: [],
        },
      });
    },
    validators: { onSubmit: schema },
  });

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (disabled) return;
        await form.handleSubmit();
      }}
      className="space-y-4"
    >
      <fieldset disabled={disabled} className="space-y-4">
        <form.AppField name="name">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Name</field.Label>
              <field.Input />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="email">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Email</field.Label>
              <field.Input type="email" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="image">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Avatar URL</field.Label>
              <field.Input placeholder="https://…" />
              <field.Error />
            </div>
          )}
        </form.AppField>
        <form.AppField name="email_verified">
          {(field) => (
            <div className="flex items-center gap-2">
              <field.Checkbox />
              <field.Label>Email verified</field.Label>
            </div>
          )}
        </form.AppField>
      </fieldset>
      <form.AppForm>
        <form.Button
          className="w-full"
          loading={mutation.isPending}
          loadingText="Saving..."
          disabled={disabled}
        >
          Save changes
        </form.Button>
      </form.AppForm>
    </form>
  );
}
