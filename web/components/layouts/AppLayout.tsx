import { Center, Flex, Spinner } from "@chakra-ui/react";
import { useIsAuthReady } from "components/AuthProvider";
import Sidebar from "components/Sidebar";
import { PropsWithChildren } from "react";

const AppLayout = ({ children }: PropsWithChildren) => {
  const isReady = useIsAuthReady();


  if (!isReady) {
    return (
      <Center h="full">
        <Spinner size="xl" />
      </Center>
    )
  }

  return (
    <Flex h="full">
      <Sidebar />
      <Flex p={4}>{children}</Flex>
    </Flex>
  );
};

export default AppLayout;
