import { Month } from "@/constants/types";
import {
  amountToFraction,
  formatAmount,
  zeroAmount,
} from "@/data/format/amount";
import { useIsMonthLoading } from "@/data/queries/month";
import { Box, Skeleton, Tag } from "@chakra-ui/react";
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

export default function ToBudget({ month }: { month?: Month }) {
  const budgetable = amountToFraction(month?.budgetable ?? zeroAmount);

  const labelId = useId();

  const isMonthLoading = useIsMonthLoading();

  return (
    <Skeleton isLoaded={!isMonthLoading}>
      <Tag colorScheme={colorScheme(budgetable.compare(0))}>
        <Box fontSize="medium" mr={2} id={labelId}>
          To Budget:
        </Box>
        <Box fontSize="large" aria-labelledby={labelId}>
          {formatAmount(month?.budgetable ?? zeroAmount)}
        </Box>
      </Tag>
    </Skeleton>
  );
}
