import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/(auth)/admin')({
  beforeLoad: async ({ context }) => {
    if (!context.session) {
      throw redirect({
        to: '/login',
        search: {
          redirect: location.href,
        },
      });
    }
  },
});
