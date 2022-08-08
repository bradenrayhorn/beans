import { Center, Flex, Spinner } from "@chakra-ui/react";
import { AuthStatus, useAuthStatus } from "components/AuthProvider";
import Sidebar from "components/Sidebar";
import { routes } from "constants/routes";
import { useBudget } from "data/queries/budget";
import { useRouter } from "next/router";
import { PropsWithChildren, useEffect } from "react";

const BudgetLayout = ({ children }: PropsWithChildren) => {
  const authStatus = useAuthStatus();
  const router = useRouter();

  const { isSuccess: isBudgetLoaded } = useBudget();

  useEffect(() => {
    if (authStatus === AuthStatus.Unauthenticated) {
      router.push(routes.login);
    }
  }, [authStatus]);

  if (authStatus !== AuthStatus.Authenticated || !isBudgetLoaded) {
    return (
      <Center h="100vh">
        <Spinner size="xl" />
      </Center>
    );
  }

  return (
    <Flex minH="100vh">
      <Sidebar />
      <Flex p={4} w="full">
        {children}
      </Flex>
    </Flex>
  );
};

export default BudgetLayout;
