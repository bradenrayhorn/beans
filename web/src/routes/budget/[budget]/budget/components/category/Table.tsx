import { Category, MonthCategory } from "@/constants/types";
import { amountToFraction, formatAmount } from "@/data/format/amount";
import { useCategories } from "@/data/queries/category";
import { useMonth } from "@/data/queries/month";
import { Flex, Heading, Skeleton, Tag } from "@chakra-ui/react";
import { useMemo } from "react";
import AssignedCell from "./AssignedCell";
import CategoryRow from "./CategoryRow";
import Header from "./Header";

const spentColumns = [
  {
    name: "Category",
    isNumeric: false,
    get: (category: Category) => category.name,
  },
  {
    name: "Assigned",
    isNumeric: true,
    get: (_: Category, monthCategory?: MonthCategory) => (
      <AssignedCell monthCategory={monthCategory} />
    ),
  },
  {
    name: "Spent",
    isNumeric: true,
    get: (_: Category, monthCategory?: MonthCategory) =>
      formatAmount(monthCategory?.activity),
  },
  {
    name: "Available",
    isNumeric: true,
    get: (_: Category, monthCategory?: MonthCategory) => (
      <Tag
        colorScheme={
          amountToFraction(monthCategory?.available).compare(0) > 0
            ? "green"
            : "gray"
        }
      >
        {formatAmount(monthCategory?.available)}
      </Tag>
    ),
  },
];

const incomeColumns = [
  {
    name: "Income",
    isNumeric: false,
    get: (category: Category) => category.name,
  },
  {
    name: "Received",
    isNumeric: true,
    get: (_: Category, monthCategory?: MonthCategory) =>
      formatAmount(monthCategory?.activity),
  },
];

export default function CategoryTable() {
  const { month, isSuccess: isMonthLoaded } = useMonth();
  const { isSuccess: areCategoriesLoaded, categoryGroups } = useCategories();

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

  return (
    <Skeleton isLoaded={isMonthLoaded && areCategoriesLoaded}>
      <Flex flexDir="column" w="full" role="table" aria-label="Expenses">
        <Header columns={spentColumns} />

        {categoryGroups
          .filter((c) => !c.is_income)
          .map((group) => (
            <Flex
              px={4}
              py={3}
              flexDir="column"
              w="full"
              key={group.id}
              role="rowgroup"
              aria-label={group.name}
            >
              <Flex mb={1} aria-hidden="true">
                <Heading size="xs">{group.name}</Heading>
              </Flex>

              <Flex
                flexDir="column"
                gap={2}
                fontSize="sm"
                w="full"
                role="rowgroup"
              >
                {group.categories.map((category) => (
                  <CategoryRow
                    key={category.id}
                    columns={spentColumns}
                    category={category}
                    monthCategory={categories[category.id]}
                  />
                ))}
              </Flex>
            </Flex>
          ))}
      </Flex>

      <Flex flexDir="column" w="full" mt={12} role="table" aria-label="Income">
        <Header columns={incomeColumns} />

        <Flex px={4} fontSize="sm" w="full" flexDir="column" gap={2} mt={3}>
          {categoryGroups
            .filter((c) => c.is_income)
            .map((group) =>
              group.categories.map((category) => (
                <CategoryRow
                  key={category.id}
                  columns={incomeColumns}
                  category={category}
                  monthCategory={categories[category.id]}
                />
              ))
            )}
        </Flex>
      </Flex>
    </Skeleton>
  );
}
