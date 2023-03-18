import { Category, MonthCategory } from "@/constants/types";
import { Flex } from "@chakra-ui/react";
import { ReactNode } from "react";

type Column = {
  name: string;
  isNumeric: boolean;
  get: (
    category: Category,
    monthCategory?: MonthCategory
  ) => string | ReactNode;
};

export default function CategoryRow({
  columns,
  category,
  monthCategory,
}: {
  columns: Column[];
  category: Category;
  monthCategory?: MonthCategory;
}) {
  return (
    <Flex key={category.id} w="full" alignItems="center" role="row">
      {columns.map((column) => (
        <Flex
          key={column.name}
          flex={1}
          justifyContent={column.isNumeric ? "flex-end" : "flex-start"}
          role="cell"
        >
          {column.get(category, monthCategory)}
        </Flex>
      ))}
    </Flex>
  );
}
