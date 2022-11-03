import { routes } from "@/constants/routes";
import { AuthStatus, useAuthStatus } from "@/context/AuthContext";
import { Center, Spinner } from "@chakra-ui/react";
import { Navigate, Outlet } from "react-router-dom";

export default function MustBeAuthed() {
  const status = useAuthStatus();

  if (status === AuthStatus.Loading) {
    return (
      <Center mt="4">
        <Spinner />
      </Center>
    );
  }

  if (status === AuthStatus.Error || status === AuthStatus.Unauthenticated) {
    return <Navigate to={routes.login} replace={true} />;
  }

  return (
    <>
      <Outlet />
    </>
  );
}
