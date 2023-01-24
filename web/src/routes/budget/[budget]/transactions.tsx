import { Box, Flex, Heading } from "@chakra-ui/react";
import TransactionList from "@/components/transactions/TransactionList";
import AddTransactionButton from "@/components/transactions/AddTransactionButton";

export default function TransactionsPage() {
  return (
    <Box as="main" w="full">
      <Flex justify="space-between" align="center" mb={8}>
        <Flex align="center">
          <Heading size="lg">Transactions</Heading>
        </Flex>
        <AddTransactionButton />
      </Flex>
      <TransactionList />
    </Box>
  );
}
