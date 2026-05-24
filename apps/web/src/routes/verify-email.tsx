import { useQuery } from '@tanstack/react-query';
import { createFileRoute, Link } from '@tanstack/react-router';
import { CheckCircle2, Loader2, ShieldAlert } from 'lucide-react';
import { z } from 'zod';

import { buttonVariants } from '@/components/ui/button';
import { orpc } from '@/lib/orpc';

export const Route = createFileRoute('/verify-email')({
  validateSearch: z.object({
    token: z.string().min(1).optional(),
    callback_url: z.string().optional(),
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { token, callback_url } = Route.useSearch();

  const { isLoading, isError, isSuccess } = useQuery(
    orpc.getAuthEmailPasswordVerifyEmail.queryOptions({
      input: { query: { token: token ?? '', callback_url } },
      enabled: !!token,
      retry: false,
      meta: { silent: true },
    }),
  );

  if (!token) {
    return (
      <CenteredState
        icon={<ShieldAlert className="size-10 text-destructive" />}
        title="Invalid verification link"
        body="The link is missing a token. Open the link from your verification email."
        cta={
          <Link to="/login" className={buttonVariants()}>
            Back to sign in
          </Link>
        }
      />
    );
  }

  if (isLoading) {
    return (
      <CenteredState
        icon={
          <Loader2 className="size-10 animate-spin text-muted-foreground" />
        }
        title="Verifying your email"
        body="Hang tight — this only takes a moment."
      />
    );
  }

  if (isError) {
    return (
      <CenteredState
        icon={<ShieldAlert className="size-10 text-destructive" />}
        title="This link is invalid or expired"
        body="Request a new verification email from your account settings or sign in again."
        cta={
          <Link to="/login" className={buttonVariants()}>
            Back to sign in
          </Link>
        }
      />
    );
  }

  if (isSuccess) {
    return (
      <CenteredState
        icon={<CheckCircle2 className="size-10 text-primary" />}
        title="Email verified"
        body="Your email address is confirmed. You can continue to your account."
        cta={
          <Link to="/" className={buttonVariants()}>
            Continue
          </Link>
        }
      />
    );
  }

  return null;
}

function CenteredState({
  icon,
  title,
  body,
  cta,
}: {
  icon: React.ReactNode;
  title: string;
  body: string;
  cta?: React.ReactNode;
}) {
  return (
    <div className="mx-auto w-full mt-10 max-w-md p-6 text-center space-y-4">
      <div className="flex justify-center">{icon}</div>
      <h1 className="text-2xl font-bold">{title}</h1>
      <p className="text-muted-foreground">{body}</p>
      {cta}
    </div>
  );
}
