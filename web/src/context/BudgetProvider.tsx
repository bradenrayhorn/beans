import Sidebar from "@/components/Sidebar";
import { currentDate } from "@/data/format/date";
import { useBudget } from "@/data/queries/budget";
import { Center, Flex, Spinner } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import MonthProvider from "./MonthProvider";

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
    <MonthProvider defaultMonthDate={currentDate()}>
      <Flex minH="100vh">
        <Sidebar />
        <Flex w="full">{children}</Flex>
      </Flex>
    </MonthProvider>
  );
}
