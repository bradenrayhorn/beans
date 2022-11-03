import { routes } from "@/constants/routes";
import { AuthStatus, useAuthStatus } from "@/context/AuthContext";
import { Center, Spinner } from "@chakra-ui/react";
import { Navigate, Outlet } from "react-router-dom";

export default function MustBeNotAuthed() {
  const status = useAuthStatus();

  if (status === AuthStatus.Loading || status === AuthStatus.Error) {
    return (
      <Center mt="4">
        <Spinner />
      </Center>
    );
  }

  if (status === AuthStatus.Authenticated) {
    return <Navigate to={routes.defaultAfterAuth} replace={true} />;
  }

  return (
    <>
      <Outlet />
    </>
  );
}
