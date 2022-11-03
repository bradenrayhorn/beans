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
import PageCard from "@/components/PageCard";
import CreateBudgetModal from "@/components/CreateBudgetModal";
import { queries, queryKeys } from "@/constants/queries";
import { useQuery } from "@tanstack/react-query";
import { routes } from "@/constants/routes";
import { generatePath, Link } from "react-router-dom";

export default function BudgetsPage() {
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
                    <Button
                      key={budget.id}
                      to={generatePath(routes.budget.index, {
                        budget: budget.id,
                      })}
                      size="sm"
                      variant="ghost"
                      as={Link}
                    >
                      {budget.name}
                    </Button>
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
}
