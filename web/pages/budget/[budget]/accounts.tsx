import { AddIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Center,
  Container,
  Flex,
  Heading,
  Spinner,
  Stat,
  StatLabel,
  StatNumber,
  Text,
  useDisclosure,
  VStack,
} from "@chakra-ui/react";
import AddAccountModal from "components/AddAccountModal";
import BudgetLayout from "components/layouts/BudgetLayout";
import PageCard from "components/PageCard";
import { useAccounts } from "data/queries/account";
import { NextPageWithLayout } from "pages/_app";

const Accounts: NextPageWithLayout = () => {
  const {
    isOpen: isAddAccountOpen,
    onOpen: onAddAccountOpen,
    onClose: onAddAccountClose,
  } = useDisclosure();
  const { accounts, isLoading, isFetching } = useAccounts();

  return (
    <Box as="main" w="full">
      <Container>
        <Flex justify="space-between" align="center" mb={8}>
          <Flex align="center">
            <Heading size="lg">Accounts</Heading>
            {isFetching && !isLoading && <Spinner ml={4} size="sm" />}
          </Flex>
          <Button size="sm" rightIcon={<AddIcon />} onClick={onAddAccountOpen}>
            Add
          </Button>
        </Flex>
        {isLoading ? (
          <Center>
            <Spinner />
          </Center>
        ) : (
          <VStack spacing={8}>
            {accounts.length < 1 && <Text as="i">No accounts found.</Text>}

            {accounts.map((account) => (
              <PageCard w="full" p={6} key={account.id}>
                <Heading size="sm">{account.name}</Heading>
                <Stat mt={2}>
                  <StatLabel>Balance</StatLabel>
                  <StatNumber>$0.00</StatNumber>
                </Stat>
              </PageCard>
            ))}
          </VStack>
        )}
      </Container>
      <AddAccountModal isOpen={isAddAccountOpen} onClose={onAddAccountClose} />
    </Box>
  );
};

Accounts.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default Accounts;
