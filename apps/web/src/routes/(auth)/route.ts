import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/(auth)')({
  beforeLoad: ({ context, location }) => {
    if (!context.session?.session) {
      throw redirect({
        to: '/login',
        search: { redirect: location.href },
      });
    }
  },
});
