import { useQuery } from "@tanstack/react-query";
import { queries, queryKeys } from "constants/queries";
import { PropsWithChildren } from "react";
import create from "zustand";
import shallow from "zustand/shallow";

interface User {
  id: string;
  username: string;
}

interface AuthState {
  isReady: boolean;
  user: User | undefined;
  setUser: (user: User) => void;
  setIsReady: (isReady: boolean) => void;
}

const useAuthStore = create<AuthState>()((set) => ({
  isReady: false,
  user: undefined,
  setUser: (user: User) => set({ user }),
  setIsReady: (isReady: boolean) => set({ isReady }),
}));

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const [setUser, setIsReady] = useAuthStore(
    (state) => [state.setUser, state.setIsReady],
    shallow
  );

  useQuery([queryKeys.me], () => queries.me(), {
    onSuccess: (data) => {
      setUser(data);
      setIsReady(true);
    },
  });

  return <>{children}</>;
};

export const useIsAuthReady = () => useAuthStore((state) => state.isReady);

export const useUser = () => {
  const user = useAuthStore((state) => state.user);
  return user;
};
