import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/(auth)/admin')({
  beforeLoad: async ({ context, location }) => {
    const userId = context.session?.user?.id;
    if (!userId) {
      throw redirect({
        to: '/login',
        search: { redirect: location.href },
      });
    }

    const roles = await context.queryClient.fetchQuery(
      context.orpc.getAuthAccessControlUsersByUserIdRoles.queryOptions({
        input: { path: { user_id: userId } },
      }),
    );

    const hasAdmin =
      roles?.body?.some((role) => role.role_name === 'admin') ?? false;

    if (!hasAdmin) {
      throw redirect({ to: '/' });
    }
  },
});
