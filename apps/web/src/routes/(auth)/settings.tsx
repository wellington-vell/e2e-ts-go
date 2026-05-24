import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useRouteContext } from '@tanstack/react-router';
import { KeyRound, MailCheck, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { orpc } from '@/lib/orpc';

export const Route = createFileRoute('/(auth)/settings')({
  component: SettingsPage,
});

function SettingsPage() {
  const { session } = useRouteContext({ from: '__root__' });

  return (
    <div className="mx-auto w-full max-w-3xl px-4 py-8 space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Settings</h1>
        <p className="text-muted-foreground">
          Manage your profile and account security.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Profile</CardTitle>
          <CardDescription>Read-only snapshot of your account.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Name</span>
            <span className="font-medium">{session?.user?.name ?? '—'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Email</span>
            <span className="font-medium">{session?.user?.email ?? '-'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Email verified</span>
            {session?.user?.email_verified ? (
              <Badge variant="default">Verified</Badge>
            ) : (
              <Badge variant="secondary">Not verified</Badge>
            )}
          </div>
        </CardContent>
      </Card>

      {!session?.user?.email_verified && <ResendVerificationCard />}

      <ChangeEmailCard />

      <ResetPasswordCard email={session?.user?.email ?? ''} />
    </div>
  );
}

function ResendVerificationCard() {
  const mutation = useMutation(
    orpc.postAuthEmailPasswordSendEmailVerification.mutationOptions({
      onSuccess: () => {
        toast.success('Verification email sent — check your inbox.');
      },
    }),
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <MailCheck className="size-4" />
          Verify your email
        </CardTitle>
        <CardDescription>
          We'll send you a link to confirm your address.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button
          onClick={() =>
            void mutation.mutateAsync({
              body: {
                callback_url: `${window.location.origin}/verify-email`,
              },
            })
          }
          disabled={mutation.isPending}
        >
          {mutation.isPending ? 'Sending...' : 'Send verification email'}
        </Button>
      </CardContent>
    </Card>
  );
}

const emailSchema = z.object({
  new_email: z.email('Invalid email address'),
});

function ChangeEmailCard() {
  const mutation = useMutation(
    orpc.postAuthEmailPasswordRequestEmailChange.mutationOptions({
      onSuccess: () => {
        toast.success('Confirmation link sent to your new email address.');
        form.reset();
      },
    }),
  );

  const form = useAppForm({
    defaultValues: { new_email: '' },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        body: {
          new_email: value.new_email,
          callback_url: `${window.location.origin}/verify-email`,
        },
      });
    },
    validators: { onSubmit: emailSchema },
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <ShieldCheck className="size-4" />
          Change email
        </CardTitle>
        <CardDescription>
          You'll receive a confirmation link at the new address.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={async (e) => {
            e.preventDefault();
            e.stopPropagation();
            await form.handleSubmit();
          }}
          className="space-y-3"
        >
          <form.AppField name="new_email">
            {(field) => (
              <div className="space-y-2">
                <field.Label>New email</field.Label>
                <field.Input type="email" />
                <field.Error />
              </div>
            )}
          </form.AppField>
          <form.AppForm>
            <form.Button loading={mutation.isPending} loadingText="Sending...">
              Send confirmation
            </form.Button>
          </form.AppForm>
        </form>
      </CardContent>
    </Card>
  );
}

function ResetPasswordCard({ email }: { email: string }) {
  const mutation = useMutation(
    orpc.postAuthEmailPasswordRequestPasswordReset.mutationOptions({
      onSuccess: () => {
        toast.success('Password reset link sent — check your inbox.');
      },
    }),
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <KeyRound className="size-4" />
          Reset password
        </CardTitle>
        <CardDescription>
          We'll email you a link to set a new password.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button
          onClick={() =>
            void mutation.mutateAsync({
              body: {
                email,
                callback_url: `${window.location.origin}/reset-password`,
              },
            })
          }
          disabled={mutation.isPending || !email}
        >
          {mutation.isPending ? 'Sending...' : 'Send password reset email'}
        </Button>
      </CardContent>
    </Card>
  );
}
