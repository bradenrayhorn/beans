import { MonthCategory } from "@/constants/types";
import { formatAmount, zeroAmount } from "@/data/format/amount";
import { SimpleGrid, Stat, StatLabel, StatNumber } from "@chakra-ui/react";
import { useId } from "react";

export default function CategoryStats({
  category,
}: {
  category?: MonthCategory;
}) {
  const assignedID = useId();
  const activityID = useId();
  const availableID = useId();

  const assigned = category?.assigned ?? zeroAmount;
  const activity = category?.activity ?? zeroAmount;
  const available = category?.available ?? zeroAmount;

  return (
    <SimpleGrid columns={3}>
      <Stat role="group" flexGrow={0} aria-labelledby={availableID}>
        <StatNumber>{formatAmount(available)}</StatNumber>
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
        aria-labelledby={activityID}
      >
        <StatNumber>{formatAmount(activity)}</StatNumber>
        <StatLabel id={activityID}>Activity</StatLabel>
      </Stat>
    </SimpleGrid>
  );
}
