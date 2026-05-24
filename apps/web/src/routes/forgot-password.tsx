import { useMutation } from '@tanstack/react-query';
import { createFileRoute, Link, redirect } from '@tanstack/react-router';
import { MailCheck } from 'lucide-react';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Button } from '@/components/ui/button';
import { orpc } from '@/lib/orpc';

export const Route = createFileRoute('/forgot-password')({
  beforeLoad: ({ context }) => {
    if (context.session?.user) {
      throw redirect({ to: '/' });
    }
  },
  component: RouteComponent,
});

const schema = z.object({
  email: z.email('Invalid email address'),
});

function RouteComponent() {
  const mutation = useMutation(
    orpc.postAuthEmailPasswordRequestPasswordReset.mutationOptions({
      onSuccess: () => {},
    }),
  );

  const form = useAppForm({
    defaultValues: { email: '' },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        body: {
          email: value.email,
          callback_url: `${window.location.origin}/reset-password`,
        },
      });
    },
    validators: { onSubmit: schema, onBlur: schema, onChange: schema },
    canSubmitWhenInvalid: true,
  });

  if (form.state.isSubmitted) {
    return (
      <div className="mx-auto w-full mt-10 max-w-md p-6 text-center space-y-4">
        <MailCheck className="mx-auto size-10 text-primary" />
        <h1 className="text-2xl font-bold">Check your email</h1>
        <p className="text-muted-foreground">
          If an account exists for that address, we've sent a password reset
          link. The link will expire in one hour.
        </p>
        <Link
          to="/login"
          className="text-sm text-muted-foreground hover:text-foreground underline-offset-4 hover:underline"
        >
          Back to sign in
        </Link>
      </div>
    );
  }

  if (form.state.isSubmitting) {
    return (
      <div className="mx-auto w-full mt-10 max-w-md p-6 text-center space-y-4">
        <MailCheck className="mx-auto size-10 text-primary" />
        <h1 className="text-2xl font-bold">Sending reset link...</h1>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full mt-10 max-w-md p-6">
      <h1 className="mb-2 text-center text-3xl font-bold">Forgot password</h1>
      <p className="mb-6 text-center text-sm text-muted-foreground">
        Enter your email address and we'll send you a link to reset your
        password.
      </p>

      <form
        onSubmit={async (e) => {
          e.preventDefault();
          e.stopPropagation();
          await form.handleSubmit();
        }}
        className="space-y-4"
      >
        <form.AppField name="email">
          {(field) => (
            <div className="space-y-2">
              <field.Label>Email</field.Label>
              <field.Input type="email" />
              <field.Error />
            </div>
          )}
        </form.AppField>

        <form.AppForm>
          <form.Button
            className="w-full"
            loading={mutation.isPending}
            loadingText="Sending..."
          >
            Send reset link
          </form.Button>
        </form.AppForm>
      </form>

      <div className="mt-4 text-center">
        <Button variant="link">
          <Link to="/login">Back to sign in</Link>
        </Button>
      </div>
    </div>
  );
}
