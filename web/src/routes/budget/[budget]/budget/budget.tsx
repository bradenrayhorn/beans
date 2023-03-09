import CategoryStats from "@/components/budget/categories/CategoryStats";
import EditButton from "@/components/budget/categories/EditButton";
import PageCard from "@/components/PageCard";
import { MonthCategory } from "@/constants/types";
import { useMonthID } from "@/context/MonthProvider";
import { zeroAmount } from "@/data/format/amount";
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
import MonthHeader from "./components/MonthHeader";
import ToBudget from "./components/ToBudget";

export default function BudgetPage() {
  const monthID = useMonthID();

  const { isLoading: isMonthLoading, month } = useMonth({ monthID });
  const { isLoading: areCategoriesLoading, categoryGroups } = useCategories();

  const categories = useMemo(() => {
    const monthCategories = {} as { [key: string]: MonthCategory | undefined };
    categoryGroups
      .flatMap((group) => group.categories)
      .forEach((category) => {
        monthCategories[category.id] = month?.categories?.find(
          ({ category_id }) => category_id === category.id
        );
      });
    return monthCategories;
  }, [month?.categories, categoryGroups]);

  if (isMonthLoading || areCategoriesLoading || !month) {
    return <Spinner />;
  }

  return (
    <Flex as="main" w="full" flexDir="column">
      <Flex mb={8} alignItems="center" justifyContent="space-between">
        <MonthHeader month={month} />
        <ToBudget month={month} />
      </Flex>
      <VStack
        as={List}
        aria-label="Categories"
        w="full"
        align="flex-start"
        gap={6}
      >
        {categoryGroups
          .filter((group) => !group.is_income)
          .map((group) => (
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
                    <Flex justify="space-between" align="center" mb={2}>
                      <Heading size="md">{category.name}</Heading>
                      <EditButton
                        category={category}
                        monthID={monthID}
                        amount={categories[category.id]?.assigned ?? zeroAmount}
                      />
                    </Flex>
                    <CategoryStats category={categories[category.id]} />
                  </PageCard>
                ))}
              </VStack>
            </ListItem>
          ))}
      </VStack>
    </Flex>
  );
}
