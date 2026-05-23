import React from 'react';

import type { AdminUserState } from '@/lib/api/types.gen';

interface AdminUsersContextValue {
  bannedMap: Map<string, AdminUserState>;
}

const AdminUsersContext = React.createContext<AdminUsersContextValue | null>(
  null,
);

export function AdminUsersProvider({
  bannedMap,
  children,
}: {
  bannedMap: Map<string, AdminUserState>;
  children: React.ReactNode;
}) {
  const value = React.useMemo(() => ({ bannedMap }), [bannedMap]);
  return (
    <AdminUsersContext.Provider value={value}>
      {children}
    </AdminUsersContext.Provider>
  );
}

export function useAdminUsers() {
  const ctx = React.useContext(AdminUsersContext);
  if (!ctx)
    throw new Error('useAdminUsers must be used within AdminUsersProvider');
  return ctx;
}
