import { Flex } from "@chakra-ui/react";
import CategoryTable from "./components/category/Table";
import ForNextMonth from "./components/ForNextMonth";
import MonthHeader from "./components/MonthHeader";
import ToBudget from "./components/ToBudget";

export default function BudgetPage() {
  return (
    <Flex as="main" w="full">
      <Flex grow="1" flexDir="column">
        <Flex alignItems="center" justifyContent="space-between" p={4}>
          <MonthHeader />
        </Flex>

        <CategoryTable />
      </Flex>
      <Flex
        shrink={0}
        bg={"gray.50"}
        p={2}
        shadow="md"
        minW={72}
        gap={6}
        flexDir="column"
      >
        <ToBudget />

        <ForNextMonth />
      </Flex>
    </Flex>
  );
}
