import { useMutation, useQuery } from '@tanstack/react-query';
import {
  Link,
  useNavigate,
  useRouteContext,
  useRouter,
  type LinkProps,
} from '@tanstack/react-router';
import { AlertTriangle, DoorOpen, Moon, Sun, User2 } from 'lucide-react';
import { useTheme } from 'next-themes';
import React from 'react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { orpc, queryClient } from '@/lib/orpc';
import type { routeTree } from '@/routeTree.gen';

type AppLink = {
  to: LinkProps<typeof routeTree>['to'];
  label: string;
  visible?: boolean;
  icon?: React.ReactNode;
};

const links: AppLink[] = [
  { to: '/', label: 'Home' },
  { to: '/dashboard', label: 'Dashboard' },
  // { to: '/todos', label: 'Todos' },
  {
    to: '/admin',
    label: 'Admin',
  },
] as const;

export function Header() {
  const router = useRouter();
  const { session } = useRouteContext({ from: '__root__' });

  const userId = session?.user?.id;

  const impersonations = useQuery(
    orpc.getAuthAdminImpersonations.queryOptions({
      enabled: !!userId,
    }),
  );

  const activeImpersonation = impersonations.data?.body?.find(
    (imp) =>
      !!userId &&
      imp.target_user_id === userId &&
      !!imp.impersonation_session_id &&
      !imp.ended_at &&
      (!imp.expires_at || new Date(imp.expires_at) > new Date()),
  );

  const isImpersonating = !!activeImpersonation;

  const stopImpersonation = useMutation(
    orpc.postAuthAdminImpersonationsByImpersonationIdStop.mutationOptions({
      onSuccess: () => {
        toast.success('Stopped impersonating');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthMe.key(),
        });
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminImpersonations.key(),
        });
        void router.invalidate();
      },
    }),
  );

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 shadow-sm">
      {isImpersonating && (
        <div className="bg-amber-500/15 text-amber-700 dark:text-amber-300 flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium">
          <AlertTriangle className="h-4 w-4" />
          <span>You are impersonating another user</span>
          <Button
            variant="outline"
            size="sm"
            className="ml-2 border-amber-500/30 hover:bg-amber-500/20"
            onClick={async () => {
              if (activeImpersonation?.id) {
                await stopImpersonation.mutateAsync({
                  path: { impersonation_id: activeImpersonation.id },
                });
              }
            }}
          >
            Stop impersonating
          </Button>
        </div>
      )}
      <div className="flex h-14 items-center justify-between px-4">
        <nav className="flex gap-1">
          {links.map((link) => {
            const { to, label, visible } = link;

            if (visible !== undefined && !visible) return null;
            return (
              <Link
                key={to}
                to={to}
                className="flex items-center rounded-md px-3 py-1.5 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                activeProps={{
                  className: 'bg-muted text-foreground',
                }}
              >
                {label}
              </Link>
            );
          })}
        </nav>
        <div className="flex items-center gap-2">
          <ModeToggle />
          <UserMenu />
        </div>
      </div>
    </header>
  );
}

function UserMenu() {
  const navigate = useNavigate();
  const router = useRouter();
  const { session } = useRouteContext({ from: '__root__' });

  const signOut = useMutation(
    orpc.postAuthSignOut.mutationOptions({
      onSuccess: () => {
        toast.success('Successfully signed out!');
        queryClient.removeQueries({
          queryKey: orpc.getAuthMe.key(),
        });
        void router.invalidate();
        void navigate({
          to: '/login',
        });
      },
    }),
  );

  if (!session?.user || !session?.session) {
    return (
      <Link to="/login">
        <Button variant="outline">Sign In</Button>
      </Link>
    );
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger render={<Button variant="outline" />}>
        {session.user.name}
      </DropdownMenuTrigger>
      <DropdownMenuContent className="bg-card w-full">
        <DropdownMenuGroup>
          <DropdownMenuLabel className="text-center">
            My Account
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => {
              void navigate({ to: '/settings' });
            }}
          >
            <User2 />
            {session.user.email}
          </DropdownMenuItem>
          <DropdownMenuItem
            variant="destructive"
            onClick={async () => {
              await signOut.mutateAsync({});
            }}
          >
            <DoorOpen />
            Sign Out
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function ModeToggle() {
  const { setTheme } = useTheme();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={(props) => (
          <Button variant="outline" size="icon" {...props}>
            <Sun className="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90" />
            <Moon className="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0" />
            <span className="sr-only">Toggle theme</span>
          </Button>
        )}
      />
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme('light')}>
          Light
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('dark')}>
          Dark
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('system')}>
          System
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
