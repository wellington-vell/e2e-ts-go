import { useMutation } from '@tanstack/react-query';
import { createFileRoute, Link, useRouteContext } from '@tanstack/react-router';
import { MailWarning } from 'lucide-react';
import { toast } from 'sonner';

import {
  Alert,
  AlertAction,
  AlertDescription,
  AlertTitle,
} from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { orpc } from '@/lib/orpc';

export const Route = createFileRoute('/(auth)/dashboard')({
  component: RouteComponent,
});

function RouteComponent() {
  const { session } = useRouteContext({ from: '__root__' });

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-8 space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      {session?.user && !session.user.email_verified && (
        <UnverifiedEmailBanner />
      )}
      <p className="text-muted-foreground">
        Welcome back{session?.user?.name ? `, ${session.user.name}` : ''}.
      </p>
    </div>
  );
}

function UnverifiedEmailBanner() {
  const mutation = useMutation(
    orpc.postAuthEmailPasswordSendEmailVerification.mutationOptions({
      onSuccess: () => {
        toast.success('Verification email sent — check your inbox.');
      },
    }),
  );

  return (
    <Alert>
      <MailWarning />
      <AlertTitle>Verify your email address</AlertTitle>
      <AlertDescription>
        Confirm your email to unlock all account features. Open{' '}
        <Button variant="link">
          <Link to="/settings">Settings</Link>
        </Button>{' '}
        for more options.
      </AlertDescription>
      <AlertAction>
        <Button
          size="sm"
          variant="outline"
          onClick={() =>
            void mutation.mutateAsync({
              body: {
                callback_url: `${window.location.origin}/verify-email`,
              },
            })
          }
          disabled={mutation.isPending}
        >
          {mutation.isPending ? 'Sending...' : 'Resend'}
        </Button>
      </AlertAction>
    </Alert>
  );
}
