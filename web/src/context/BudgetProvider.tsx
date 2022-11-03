import Sidebar from "@/components/Sidebar";
import { useBudget } from "@/data/queries/budget";
import { Center, Flex, Spinner } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

export default function BudgetProvider({ children }: PropsWithChildren) {
  const { isSuccess: isBudgetLoaded } = useBudget();

  if (!isBudgetLoaded) {
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
}
