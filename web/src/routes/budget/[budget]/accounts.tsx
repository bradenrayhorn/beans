import AddAccountModal from "@/components/AddAccountModal";
import PageCard from "@/components/PageCard";
import { useAccounts } from "@/data/queries/account";
import { AddIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Center,
  Container,
  Flex,
  Heading,
  List,
  ListItem,
  Spinner,
  Stat,
  StatLabel,
  StatNumber,
  Text,
  useDisclosure,
} from "@chakra-ui/react";

export default function AccountsPage() {
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
          <>
            {accounts.length < 1 && <Text as="i">No accounts found.</Text>}

            <List spacing={8}>
              {accounts.map((account) => (
                <PageCard w="full" p={6} key={account.id} as={ListItem}>
                  <Heading size="sm">{account.name}</Heading>
                  <Stat mt={2}>
                    <StatLabel>Balance</StatLabel>
                    <StatNumber>$0.00</StatNumber>
                  </Stat>
                </PageCard>
              ))}
            </List>
          </>
        )}
      </Container>
      <AddAccountModal isOpen={isAddAccountOpen} onClose={onAddAccountClose} />
    </Box>
  );
}
