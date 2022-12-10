import CategoryStats from "@/components/budget/categories/CategoryStats";
import EditButton from "@/components/budget/categories/EditButton";
import MonthHeader from "@/components/budget/MonthHeader";
import PageCard from "@/components/PageCard";
import { Amount } from "@/constants/types";
import { useMonthID } from "@/context/MonthProvider";
import { useCategories } from "@/data/queries/category";
import { useMonth } from "@/data/queries/month";
import {
  Flex,
  Heading,
  List,
  ListItem,
  Spinner,
  VStack,
} from "@chakra-ui/react";
import { useMemo } from "react";

export default function BudgetPage() {
  const monthID = useMonthID();

  const { isLoading: isMonthLoading, month } = useMonth({ monthID });
  const { isLoading: areCategoriesLoading, categoryGroups } = useCategories();

  const categoryBudgets = useMemo(() => {
    const budgets = {} as { [key: string]: Amount };
    categoryGroups
      .flatMap((group) => group.categories)
      .forEach((category) => {
        budgets[category.id] = month?.categories?.find(
          ({ category_id }) => category_id === category.id
        )?.assigned ?? { exponent: 0, coefficient: 0 };
      });
    return budgets;
  }, [month?.categories, categoryGroups]);

  if (isMonthLoading || areCategoriesLoading || !month) {
    return <Spinner />;
  }

  return (
    <Flex as="main" w="full" flexDir="column">
      <Flex mb={8} alignItems="center" gap={2}>
        <MonthHeader month={month} />
      </Flex>
      <VStack
        as={List}
        aria-label="Categories"
        w="full"
        align="flex-start"
        gap={6}
      >
        {categoryGroups.map((group) => (
          <ListItem w="full" key={group.id}>
            <Heading textTransform="uppercase" mb={4} size="lg">
              {group.name}
            </Heading>
            <VStack as={List} gap={4} w="full" align="flex-start">
              {group.categories.map((category) => (
                <PageCard
                  key={category.id}
                  as={ListItem}
                  w="full"
                  p={4}
                  display="flex"
                  flexDir="column"
                >
                  <Flex justify="space-between" align="center">
                    <Heading size="md">{category.name}</Heading>
                    <EditButton
                      category={category}
                      monthID={monthID}
                      amount={categoryBudgets[category.id]}
                    />
                  </Flex>
                  <Flex mt={2}>
                    <CategoryStats assigned={categoryBudgets[category.id]} />
                  </Flex>
                </PageCard>
              ))}
            </VStack>
          </ListItem>
        ))}
      </VStack>
    </Flex>
  );
}
