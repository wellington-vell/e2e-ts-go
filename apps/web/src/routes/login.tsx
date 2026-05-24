import { useMutation } from '@tanstack/react-query';
import {
  createFileRoute,
  Link,
  redirect,
  useNavigate,
  useRouter,
} from '@tanstack/react-router';
import React from 'react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Button } from '@/components/ui/button';
import type { User } from '@/lib/api/types.gen';
import { orpc, queryClient } from '@/lib/orpc';

function sanitizeRedirect(url: unknown): string {
  if (typeof url !== 'string' || !url.startsWith('/') || url.startsWith('//')) {
    return '/';
  }
  return url;
}

export const Route = createFileRoute('/login')({
  validateSearch: z.object({
    redirect: z.string().optional().transform(sanitizeRedirect).default('/'),
  }),
  beforeLoad: async ({ context, search }) => {
    if (context.session?.user) {
      throw redirect({ to: search.redirect });
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  const [showSignIn, setShowSignIn] = React.useState(true);

  return (
    <>
      <React.Activity mode={showSignIn ? 'visible' : 'hidden'}>
        <SignInForm onSwitchToSignUp={() => setShowSignIn(false)} />
      </React.Activity>

      <React.Activity mode={showSignIn ? 'hidden' : 'visible'}>
        <SignUpForm onSwitchToSignIn={() => setShowSignIn(true)} />
      </React.Activity>
    </>
  );
}

export function SignInForm({
  onSwitchToSignUp,
}: {
  onSwitchToSignUp: () => void;
}) {
  const search = Route.useSearch();
  const navigate = useNavigate();
  const router = useRouter();
  const defaultValues: Pick<User, 'email'> & { password: string } = {
    email: '',
    password: '',
  };

  const schema = z.object({
    email: z.email('Invalid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
  });

  const signIn = useMutation(
    orpc.postAuthEmailPasswordSignIn.mutationOptions({
      onSuccess: async () => {
        toast.success('Successfully signed in!');
        await queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
        await router.invalidate();
        void navigate({ to: search.redirect });
      },
    }),
  );

  const form = useAppForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await signIn.mutateAsync({
        body: {
          email: value.email,
          password: value.password,
        },
      });
    },
    validators: {
      onSubmit: schema,
      onBlur: schema,
      onChange: schema,
    },
    canSubmitWhenInvalid: true,
  });

  return (
    <div className="mx-auto w-full mt-10 max-w-md p-6">
      <h1 className="mb-6 text-center text-3xl font-bold">Welcome Back</h1>

      <form
        onSubmit={async (e) => {
          e.preventDefault();
          e.stopPropagation();
          await form.handleSubmit();
        }}
        className="space-y-4"
      >
        <div>
          <form.AppField name="email">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Email</field.Label>
                <field.Input type="email" />
                <field.Error />
              </div>
            )}
          </form.AppField>
        </div>

        <div>
          <form.AppField name="password">
            {(field) => (
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <field.Label>Password</field.Label>
                  <Button variant="link">
                    <Link to="/forgot-password">Forgot password?</Link>
                  </Button>
                </div>
                <field.Input type="password" />
                <field.Error />
              </div>
            )}
          </form.AppField>
        </div>

        <form.AppForm>
          <form.Button
            className="w-full"
            disabled={form.state.isSubmitting}
            loading={form.state.isSubmitting}
            loadingText="Submitting..."
          >
            Sign In
          </form.Button>
        </form.AppForm>
      </form>

      <div className="mt-4 text-center">
        <Button
          variant="link"
          onClick={onSwitchToSignUp}
          className="text-indigo-600 hover:text-indigo-800"
        >
          Need an account? Sign Up
        </Button>
      </div>
    </div>
  );
}

export function SignUpForm({
  onSwitchToSignIn,
}: {
  onSwitchToSignIn: () => void;
}) {
  const router = useRouter();
  const defaultValues: Pick<User, 'email' | 'name'> & {
    password: string;
  } = {
    email: '',
    password: '',
    name: '',
  };

  const schema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.email('Invalid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
  });

  const signUp = useMutation(
    orpc.postAuthEmailPasswordSignUp.mutationOptions({
      onSuccess: async () => {
        toast.success('Successfully signed up!');
        await queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
        await router.invalidate();
      },
    }),
  );

  const form = useAppForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await signUp.mutateAsync({
        body: {
          email: value.email,
          password: value.password,
          name: value.name,
        },
      });
    },
    validators: {
      onSubmit: schema,
      onBlur: schema,
      onChange: schema,
    },
    canSubmitWhenInvalid: true,
  });

  return (
    <div className="mx-auto w-full mt-10 max-w-md p-6">
      <h1 className="mb-6 text-center text-3xl font-bold">Create Account</h1>

      <form
        onSubmit={async (e) => {
          e.preventDefault();
          e.stopPropagation();
          await form.handleSubmit();
        }}
        className="space-y-4"
      >
        <div>
          <form.AppField name="name">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Name</field.Label>
                <field.Input />
                <field.Error />
              </div>
            )}
          </form.AppField>
        </div>

        <div>
          <form.AppField name="email">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Email</field.Label>
                <field.Input type="email" />
                <field.Error />
              </div>
            )}
          </form.AppField>
        </div>

        <div>
          <form.AppField name="password">
            {(field) => (
              <div className="space-y-2">
                <field.Label>Password</field.Label>
                <field.Input type="password" />
                <field.Error />
              </div>
            )}
          </form.AppField>
        </div>

        <form.AppForm>
          <form.Button
            className="w-full"
            disabled={form.state.isSubmitting}
            loading={form.state.isSubmitting}
            loadingText="Submitting..."
          >
            Sign Up
          </form.Button>
        </form.AppForm>
      </form>

      <div className="mt-4 text-center">
        <Button
          variant="link"
          onClick={onSwitchToSignIn}
          className="text-indigo-600 hover:text-indigo-800"
        >
          Already have an account? Sign In
        </Button>
      </div>
    </div>
  );
}
