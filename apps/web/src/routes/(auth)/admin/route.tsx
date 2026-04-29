import { createFileRoute, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/(auth)/admin')({
  beforeLoad: async () => {},
  component: () => <Outlet />,
});
