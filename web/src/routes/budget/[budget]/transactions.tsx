import { Box, Flex, Heading } from "@chakra-ui/react";
import AddTransaction from "@/components/transactions/AddTransaction";
import TransactionList from "@/components/transactions/TransactionList";

export default function TransactionsPage() {
  return (
    <Box as="main" w="full">
      <Flex justify="space-between" align="center" mb={8}>
        <Flex align="center">
          <Heading size="lg">Transactions</Heading>
        </Flex>
        <AddTransaction />
      </Flex>
      <TransactionList />
    </Box>
  );
}
