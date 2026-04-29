import { Link, type LinkProps } from '@tanstack/react-router';
import { AlertTriangle, Moon, Sun } from 'lucide-react';
import { useTheme } from 'next-themes';
import React from 'react';

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
import { useAuth } from '@/context/auth';
import type { routeTree } from '@/routeTree.gen';

export function Header() {
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
      // visible: !!session && session.user.role?.includes('admin'),
    },
  ] as const;

  const isImpersonating = false;

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
            onClick={async () => {}}
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
  const { user, session, isLoading, signOut } = useAuth();

  if (isLoading) {
    return (
      <Button variant="outline" disabled>
        Loading...
      </Button>
    );
  }

  if (!session || !user) {
    return (
      <Link to="/login">
        <Button variant="outline">Sign In</Button>
      </Link>
    );
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger render={<Button variant="outline" />}>
        {user.name}
      </DropdownMenuTrigger>
      <DropdownMenuContent className="bg-card w-full">
        <DropdownMenuGroup>
          <DropdownMenuLabel>My Account</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem>{user.email}</DropdownMenuItem>
          <DropdownMenuItem
            variant="destructive"
            onClick={async () => {
              await signOut.mutateAsync({});
            }}
          >
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
