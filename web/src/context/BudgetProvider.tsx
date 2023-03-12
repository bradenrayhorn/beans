import Sidebar from "@/components/Sidebar";
import { useBudget } from "@/data/queries/budget";
import { Center, Flex, Spinner } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import MonthProvider from "./MonthProvider";

export default function BudgetProvider({ children }: PropsWithChildren) {
  const { isSuccess: isBudgetLoaded, budget } = useBudget();

  if (!isBudgetLoaded) {
    return (
      <Center h="100vh">
        <Spinner size="xl" />
      </Center>
    );
  }

  return (
    <MonthProvider defaultMonthID={budget.latest_month_id}>
      <Flex minH="100vh">
        <Sidebar />
        <Flex w="full">{children}</Flex>
      </Flex>
    </MonthProvider>
  );
}
