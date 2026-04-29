import { useMutation, useQuery } from '@tanstack/react-query';
import type { UseMutationResult } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import type { GetMeResponse, SignOutRequest, Session, User } from 'authula';
import type {
  BanUserRequest,
  GetAllImpersonationsResponse,
  GetAllUsersResponse,
  GetBannedUserStatesResponse,
  RevokeSessionRequest,
  SignInRequest,
  SignUpRequest,
  StartImpersonationRequest,
  StopImpersonationResponse,
  UserRoleInfo,
} from 'authula/plugins';
import React from 'react';
import { toast } from 'sonner';

import { authClient } from '@/lib/auth-client';
import { queryClient } from '@/lib/orpc';

export const AUTH_QUERY_KEY = ['auth'] as const;
export const ADMIN_USERS_QUERY_KEY = ['admin', 'users'] as const;
export const ADMIN_BANNED_STATES_QUERY_KEY = [
  'admin',
  'banned-states',
] as const;

type AuthContextValue = {
  user: User | null;
  session: Session | null;
  isLoading: boolean;
  error: Error | null;
  signIn: UseMutationResult<void, Error, SignInRequest>;
  signUp: UseMutationResult<void, Error, SignUpRequest>;
  signOut: UseMutationResult<void, Error, SignOutRequest | undefined>;
  banUser: UseMutationResult<
    void,
    Error,
    { userId: string; params?: BanUserRequest }
  >;
  unbanUser: UseMutationResult<void, Error, string>;
  deleteUser: UseMutationResult<void, Error, string>;
  revokeSession: UseMutationResult<
    void,
    Error,
    { sessionId: string; params?: RevokeSessionRequest }
  >;
  stopImpersonation: UseMutationResult<
    StopImpersonationResponse,
    Error,
    string
  >;
  getAllImpersonations: () => Promise<GetAllImpersonationsResponse>;
  getUserRoles: (userId: string) => Promise<UserRoleInfo[]>;
  impersonate: UseMutationResult<void, Error, StartImpersonationRequest>;
};

const AuthContext = React.createContext<AuthContextValue | null>(null);

function useAuthQuery() {
  return useQuery<GetMeResponse>({
    queryKey: [...AUTH_QUERY_KEY, 'me'],
    queryFn: () => authClient.getMe(),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}

export function useAdminUsers() {
  return useQuery<GetAllUsersResponse>({
    queryKey: ADMIN_USERS_QUERY_KEY,
    queryFn: () => authClient.admin.getAllUsers(),
    staleTime: 30 * 1000,
  });
}

export function useBannedUsers() {
  return useQuery<GetBannedUserStatesResponse>({
    queryKey: ADMIN_BANNED_STATES_QUERY_KEY,
    queryFn: () => authClient.admin.getBannedUserStates(),
    staleTime: 30 * 1000,
  });
}

function invalidateAuth() {
  return queryClient.invalidateQueries({ queryKey: AUTH_QUERY_KEY });
}

function invalidateAdmin() {
  return Promise.all([
    queryClient.invalidateQueries({ queryKey: ADMIN_USERS_QUERY_KEY }),
    queryClient.invalidateQueries({ queryKey: ADMIN_BANNED_STATES_QUERY_KEY }),
  ]);
}

function useSignIn() {
  const navigate = useNavigate();
  return useMutation<void, Error, SignInRequest>({
    mutationFn: async (params) => {
      await authClient.emailPassword.signIn(params);
    },
    onSuccess: () => {
      void invalidateAuth();
      toast.success('Successfully signed in!');
      void navigate({ to: '/admin' });
    },
  });
}

function useSignUp() {
  const navigate = useNavigate();
  return useMutation<void, Error, SignUpRequest>({
    mutationFn: async (params) => {
      await authClient.emailPassword.signUp(params);
    },
    onSuccess: () => {
      void invalidateAuth();
      toast.success('Successfully signed up!');
      void navigate({ to: '/' });
    },
  });
}

function useSignOut() {
  return useMutation<void, Error, SignOutRequest | undefined>({
    mutationFn: async (params) => {
      await authClient.signOut(params ?? {});
    },
    onSuccess: () => {
      void invalidateAuth();
      toast.success('Successfully signed out!');
    },
  });
}

function useBanUser() {
  return useMutation<void, Error, { userId: string; params?: BanUserRequest }>({
    mutationFn: async ({ userId, params }) => {
      await authClient.admin.banUser(userId, params ?? {});
    },
    onSuccess: () => {
      void invalidateAuth();
      void invalidateAdmin();
      toast.success('User banned successfully');
    },
  });
}

function useUnbanUser() {
  return useMutation<void, Error, string>({
    mutationFn: async (userId) => {
      await authClient.admin.unbanUser(userId);
    },
    onSuccess: () => {
      void invalidateAuth();
      void invalidateAdmin();
      toast.success('User unbanned successfully');
    },
  });
}

function useDeleteUser() {
  return useMutation<void, Error, string>({
    mutationFn: async (userId) => {
      await authClient.admin.deleteUser(userId);
    },
    onSuccess: () => {
      void invalidateAuth();
      void invalidateAdmin();
      toast.success('User deleted successfully');
    },
  });
}

function useRevokeSession() {
  return useMutation<
    void,
    Error,
    { sessionId: string; params?: RevokeSessionRequest }
  >({
    mutationFn: async ({ sessionId, params }) => {
      await authClient.admin.revokeSession(sessionId, params ?? {});
    },
    onSuccess: () => {
      void invalidateAuth();
      toast.success('Session revoked successfully');
    },
  });
}

function useImpersonate() {
  return useMutation<void, Error, StartImpersonationRequest>({
    mutationFn: async (params) => {
      await authClient.admin.startImpersonation({
        reason: params.reason,
        targetUserId: params.targetUserId,
        expiresInSeconds: params.expiresInSeconds,
      });
    },
    onSuccess: () => {
      void invalidateAuth();
      void queryClient.invalidateQueries({
        queryKey: ADMIN_USERS_QUERY_KEY,
      });
      toast.success('Impersonation started');
    },
    // onError: (error) => {
    //   toast.error(
    //     `Failed to start impersonation: ${
    //       error instanceof Error ? error.message : 'Unknown error'
    //     }`,
    //   );
    // },
  });
}

function useStopImpersonation() {
  return useMutation<StopImpersonationResponse, Error, string>({
    mutationFn: async (impersonationId) => {
      return authClient.admin.stopImpersonation(impersonationId);
    },
    onSuccess: () => {
      void invalidateAuth();
      toast.success('Impersonation stopped');
    },
  });
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const query = useAuthQuery();

  const signIn = useSignIn();
  const signUp = useSignUp();
  const signOut = useSignOut();
  const banUser = useBanUser();
  const unbanUser = useUnbanUser();
  const deleteUser = useDeleteUser();
  const revokeSession = useRevokeSession();
  const impersonate = useImpersonate();
  const stopImpersonation = useStopImpersonation();

  const value = React.useMemo<AuthContextValue>(
    () => ({
      user: query.error ? null : (query.data?.user ?? null),
      session: query.error ? null : (query.data?.session ?? null),
      isLoading: query.isLoading,
      error: query.error,
      signIn,
      signUp,
      signOut,
      banUser,
      unbanUser,
      deleteUser,
      revokeSession,
      stopImpersonation,
      getAllImpersonations: () => authClient.admin.getAllImpersonations(),
      getUserRoles: (userId) => authClient.accessControl.getUserRoles(userId),
      impersonate,
    }),
    [
      query.error,
      query.data?.user,
      query.data?.session,
      query.isLoading,
      signIn,
      signUp,
      signOut,
      banUser,
      unbanUser,
      deleteUser,
      revokeSession,
      impersonate,
      stopImpersonation,
    ],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = React.useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return ctx;
}
