import { useMutation } from '@tanstack/react-query';
import {
  createFileRoute,
  Link,
  redirect,
  useNavigate,
} from '@tanstack/react-router';
import { ShieldAlert } from 'lucide-react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { buttonVariants } from '@/components/ui/button';
import { orpc } from '@/lib/orpc';

export const Route = createFileRoute('/reset-password')({
  validateSearch: z.object({
    token: z.string().min(1).optional(),
  }),
  beforeLoad: ({ context }) => {
    if (context.session?.user) {
      throw redirect({ to: '/' });
    }
  },
  component: RouteComponent,
});

const schema = z
  .object({
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirm: z.string().min(8, 'Password must be at least 8 characters'),
  })
  .refine((v) => v.password === v.confirm, {
    message: 'Passwords do not match',
    path: ['confirm'],
  });

function RouteComponent() {
  const { token } = Route.useSearch();
  const navigate = useNavigate();

  const mutation = useMutation(
    orpc.postAuthEmailPasswordChangePassword.mutationOptions({
      onSuccess: () => {
        toast.success('Password updated — please sign in.');
        void navigate({ to: '/login' });
      },
    }),
  );

  const form = useAppForm({
    defaultValues: { password: '', confirm: '' },
    onSubmit: async ({ value }) => {
      if (!token) return;
      await mutation.mutateAsync({
        body: { token, password: value.password },
      });
    },
    validators: { onSubmit: schema, onBlur: schema, onChange: schema },
    canSubmitWhenInvalid: true,
  });

  if (!token) {
    return (
      <div className="mx-auto w-full mt-10 max-w-md p-6 text-center space-y-4">
        <ShieldAlert className="mx-auto size-10 text-destructive" />
        <h1 className="text-2xl font-bold">Invalid reset link</h1>
        <p className="text-muted-foreground">
          The reset link is missing or malformed. Request a new one to continue.
        </p>
        <Link to="/forgot-password" className={buttonVariants()}>
          Request a new link
        </Link>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full mt-10 max-w-md p-6">
      <h1 className="mb-2 text-center text-3xl font-bold">Reset password</h1>
      <p className="mb-6 text-center text-sm text-muted-foreground">
        Choose a new password for your account.
      </p>

      <form
        onSubmit={async (e) => {
          e.preventDefault();
          e.stopPropagation();
          await form.handleSubmit();
        }}
        className="space-y-4"
      >
        <form.AppField name="password">
          {(field) => (
            <div className="space-y-2">
              <field.Label>New password</field.Label>
              <field.Input type="password" autoComplete="new-password" />
              <field.Error />
            </div>
          )}
        </form.AppField>

        <form.AppField name="confirm">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Confirm password</field.Label>
              <field.Input type="password" autoComplete="new-password" />
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
            Update password
          </form.Button>
        </form.AppForm>
      </form>

      <div className="mt-4 text-center">
        <Link
          to="/login"
          className="text-sm text-muted-foreground hover:text-foreground underline-offset-4 hover:underline"
        >
          Back to sign in
        </Link>
      </div>
    </div>
  );
}
