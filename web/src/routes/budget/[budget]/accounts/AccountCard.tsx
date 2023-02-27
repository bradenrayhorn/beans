import PageCard from "@/components/PageCard";
import { FullAccount } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import {
  Heading,
  ListItem,
  Stat,
  StatLabel,
  StatNumber,
} from "@chakra-ui/react";
import { useId } from "react";

export default function AccountCard({ account }: { account: FullAccount }) {
  const balanceID = useId();

  return (
    <PageCard w="full" p={6} key={account.id} as={ListItem}>
      <Heading size="sm">{account.name}</Heading>
      <Stat mt={2} role="group" aria-labelledby={balanceID}>
        <StatLabel id={balanceID}>Balance</StatLabel>
        <StatNumber>{formatAmount(account.balance)}</StatNumber>
      </Stat>
    </PageCard>
  );
}
