import { Amount } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import { Stat, StatLabel, StatNumber } from "@chakra-ui/react";
import { useId } from "react";

export default function CategoryStats({ assigned }: { assigned: Amount }) {
  const assignedID = useId();

  return (
    <Stat role="group" aria-labelledby={assignedID}>
      <StatNumber>{formatAmount(assigned)}</StatNumber>
      <StatLabel id={assignedID}>Assigned</StatLabel>
    </Stat>
  );
}
