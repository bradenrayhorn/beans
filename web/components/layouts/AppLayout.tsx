import { Center, Flex, Spinner } from "@chakra-ui/react";
import { AuthStatus, useAuthStatus } from "components/AuthProvider";
import Sidebar from "components/Sidebar";
import { useRouter } from "next/router";
import { PropsWithChildren, useEffect } from "react";

const AppLayout = ({ children }: PropsWithChildren) => {
  const authStatus = useAuthStatus();
  const router = useRouter();

  useEffect(() => {
    if (authStatus === AuthStatus.Unauthenticated) {
      router.push("/login");
    }
  }, [authStatus]);

  if (authStatus !== AuthStatus.Authenticated) {
    return (
      <Center h="full">
        <Spinner size="xl" />
      </Center>
    );
  }

  return (
    <Flex h="full">
      <Sidebar />
      <Flex p={4}>{children}</Flex>
    </Flex>
  );
};

export default AppLayout;
