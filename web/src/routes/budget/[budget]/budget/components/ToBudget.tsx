import { Month } from "@/constants/types";
import { amountToFraction, formatAmount } from "@/data/format/amount";
import { Box, Tag } from "@chakra-ui/react";
import { useId } from "react";

const colorScheme = (res: number): string => {
  if (res > 0) {
    return "blue";
  } else if (res < 0) {
    return "red";
  } else {
    return "green";
  }
};

export default function ToBudget({ month }: { month: Month }) {
  const budgetable = amountToFraction(month.budgetable);

  const labelId = useId();

  return (
    <Tag colorScheme={colorScheme(budgetable.compare(0))}>
      <Box fontSize="medium" mr={2} id={labelId}>
        To Budget:
      </Box>
      <Box fontSize="large" aria-labelledby={labelId}>
        {formatAmount(month.budgetable)}
      </Box>
    </Tag>
  );
}
