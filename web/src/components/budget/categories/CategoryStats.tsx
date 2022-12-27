import { MonthCategory } from "@/constants/types";
import {
  amountToFraction,
  formatAmount,
  formatFraction,
  zeroAmount,
} from "@/data/format/amount";
import { SimpleGrid, Stat, StatLabel, StatNumber } from "@chakra-ui/react";
import { useId } from "react";

export default function CategoryStats({
  category,
}: {
  category?: MonthCategory;
}) {
  const assignedID = useId();
  const spentID = useId();
  const availableID = useId();

  const assigned = category?.assigned ?? zeroAmount;
  const spent = category?.spent ?? zeroAmount;
  const available = amountToFraction(assigned).sub(amountToFraction(spent));

  return (
    <SimpleGrid columns={3}>
      <Stat role="group" flexGrow={0} aria-labelledby={availableID}>
        <StatNumber>{formatFraction(available)}</StatNumber>
        <StatLabel id={availableID}>Available</StatLabel>
      </Stat>

      <Stat
        role="group"
        display="flex"
        alignItems="flex-end"
        size="sm"
        aria-labelledby={assignedID}
      >
        <StatNumber>{formatAmount(assigned)}</StatNumber>
        <StatLabel id={assignedID}>Assigned</StatLabel>
      </Stat>

      <Stat
        role="group"
        display="flex"
        alignItems="flex-end"
        size="sm"
        aria-labelledby={spentID}
      >
        <StatNumber>{formatAmount(spent)}</StatNumber>
        <StatLabel id={spentID}>Spent</StatLabel>
      </Stat>
    </SimpleGrid>
  );
}
