import { AddIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Container,
  Skeleton,
  Text,
  useDisclosure,
  VStack,
} from "@chakra-ui/react";
import PageCard from "components/PageCard";
import CreateBudgetModal from "components/CreateBudgetModal";
import { queries, queryKeys } from "constants/queries";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { routes } from "constants/routes";

const BudgetPage = () => {
  const {
    isOpen: isCreateBudgetOpen,
    onOpen: onCreateBudgetOpen,
    onClose: onCreateBudgetClose,
  } = useDisclosure();

  const { isFetching: isLoading, data } = useQuery(
    [queryKeys.budget.getAll],
    queries.budget.getAll
  );
  let budgets = data?.data ?? [];

  return (
    <main>
      <Container mt={8}>
        <PageCard mt={2} p={4} flexDirection="column" d="flex">
          <Text>To continue, please select or create a budget.</Text>
          <Box my={8}>
            <Skeleton isLoaded={!isLoading}>
              <VStack align="flex-start">
                {budgets.length < 1 ? (
                  <Text as="i">No existing budgets found.</Text>
                ) : (
                  budgets.map((budget) => (
                    <Link
                      key={budget.id}
                      href={{
                        pathname: routes.budget.index,
                        query: { budget: budget.id },
                      }}
                      passHref
                    >
                      <Button as="a" size="sm" variant="ghost">
                        {budget.name}
                      </Button>
                    </Link>
                  ))
                )}
              </VStack>
            </Skeleton>
          </Box>
          <Button
            size="sm"
            rightIcon={<AddIcon />}
            onClick={onCreateBudgetOpen}
          >
            New Budget
          </Button>
          <CreateBudgetModal
            isOpen={isCreateBudgetOpen}
            onClose={onCreateBudgetClose}
          />
        </PageCard>
      </Container>
    </main>
  );
};

export default BudgetPage;
