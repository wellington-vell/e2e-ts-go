import { createFileRoute } from '@tanstack/react-router';
import type { User } from 'authula';
import React from 'react';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/context/auth';

export const Route = createFileRoute('/login')({
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
  const { signIn } = useAuth();
  const defaultValues: Pick<User, 'email'> & { password: string } = {
    email: '',
    password: '',
  };

  const schema = z.object({
    email: z.email('Invalid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
  });

  const form = useAppForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await signIn.mutateAsync({
        email: value.email,
        password: value.password,
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
  const { signUp } = useAuth();
  const defaultValues: Pick<User, 'email' | 'name'> & { password: string } = {
    email: '',
    password: '',
    name: '',
  };

  const schema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.email('Invalid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
  });

  const form = useAppForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await signUp.mutateAsync({
        email: value.email,
        password: value.password,
        name: value.name,
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
