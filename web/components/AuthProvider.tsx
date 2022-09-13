import { useQuery } from "@tanstack/react-query";
import { queries, queryKeys } from "constants/queries";
import { HTTPError } from "ky";
import { PropsWithChildren, useCallback, useState } from "react";
import create from "zustand";
import shallow from "zustand/shallow";

interface User {
  id: string;
  username: string;
}

export enum AuthStatus {
  Loading = 1,
  Authenticated,
  Unauthenticated,
  Error,
}

interface AuthState {
  status: AuthStatus;
  user: User | undefined;
  setUser: (user: User) => void;
  setStatus: (status: AuthStatus) => void;
}

const useAuthStore = create<AuthState>()((set) => ({
  status: AuthStatus.Loading,
  user: undefined,
  setUser: (user: User) => set({ user }),
  setStatus: (status: AuthStatus) => set({ status }),
}));

type Props = PropsWithChildren & {
  initialUser?: User;
};

export const AuthProvider = ({ children, initialUser }: Props) => {
  const [setUser, setStatus] = useAuthStore(
    (state) => [state.setUser, state.setStatus],
    shallow
  );

  // set initial data with initialUser prop from server
  useState(() => {
    if (initialUser) {
      setUser(initialUser);
    }
    setStatus(
      initialUser ? AuthStatus.Authenticated : AuthStatus.Unauthenticated
    );
  });

  useQuery([queryKeys.me], () => queries.me(), {
    retry: false,
    onSuccess: (data) => {
      setUser(data);
      setStatus(AuthStatus.Authenticated);
    },
    onError: (error: HTTPError | any) => {
      const status = error?.response?.status;
      if (status === 401) {
        setStatus(AuthStatus.Unauthenticated);
      } else {
        setStatus(AuthStatus.Error);
      }
    },
  });

  return <>{children}</>;
};

export const useOnLogin = () => {
  const [setUser, setStatus] = useAuthStore(
    (state) => [state.setUser, state.setStatus],
    shallow
  );

  return useCallback(
    (user: User) => {
      setUser(user);
      setStatus(AuthStatus.Authenticated);
    },
    [setUser, setStatus]
  );
};

export const useAuthStatus = () => useAuthStore((state) => state.status);

export const useUser = () => useAuthStore((state) => state.user);
