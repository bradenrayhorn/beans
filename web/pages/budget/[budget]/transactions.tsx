import { Box, Flex, Heading } from "@chakra-ui/react";
import BudgetLayout from "components/layouts/BudgetLayout";
import { NextPageWithLayout } from "pages/_app";
import AddTransaction from "components/transactions/AddTransaction";

const Transactions: NextPageWithLayout = () => {
  return (
    <Box as="main" w="full">
      <Flex justify="space-between" align="center" mb={8}>
        <Flex align="center">
          <Heading size="lg">Transactions</Heading>
        </Flex>
        <AddTransaction />
      </Flex>
    </Box>
  );
};

Transactions.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default Transactions;
