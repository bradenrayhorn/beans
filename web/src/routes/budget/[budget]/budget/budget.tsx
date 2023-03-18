import { TimeIcon } from "@chakra-ui/icons";
import { Button, Divider, Flex } from "@chakra-ui/react";
import CategoryTable from "./components/category/Table";
import MonthHeader from "./components/MonthHeader";
import ToBudget from "./components/ToBudget";

export default function BudgetPage() {
  return (
    <Flex as="main" w="full">
      <Flex grow="1" flexDir="column">
        <Flex alignItems="center" justifyContent="space-between" p={4}>
          <MonthHeader />
        </Flex>
        <Divider mb={4} />
        <Flex mb={4} px={4} justifyContent="flex-end">
          <Button
            leftIcon={<TimeIcon />}
            colorScheme="blue"
            size="xs"
            variant="ghost"
          >
            For Next Month
          </Button>
        </Flex>

        <CategoryTable />
      </Flex>
      <Flex shrink={0} bg={"gray.50"} p={4} shadow="md" minW={72}>
        <ToBudget />
      </Flex>
    </Flex>
  );
}
