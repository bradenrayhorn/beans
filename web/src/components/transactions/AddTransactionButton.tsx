import { useAccounts } from "@/data/queries/account";
import { AddIcon } from "@chakra-ui/icons";
import { Button, Spinner, Tooltip } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { useTransactionTableState } from "./TransactionList";

const ConditonalTooltip = ({
  hasAccounts,
  children,
}: PropsWithChildren & { hasAccounts: boolean }) =>
  !hasAccounts ? (
    <Tooltip label="You must add an account first.">{children}</Tooltip>
  ) : (
    <>{children}</>
  );

export default function AddTransaction() {
  const setIsAdding = useTransactionTableState((state) => state.setIsAdding);

  const { accounts, isLoading } = useAccounts();

  if (isLoading) {
    return <Spinner />;
  }

  const hasAccounts = accounts.length > 0;

  return (
    <ConditonalTooltip hasAccounts={hasAccounts}>
      <Button
        isDisabled={!hasAccounts}
        size="sm"
        rightIcon={<AddIcon />}
        onClick={() => {
          setIsAdding(true);
        }}
      >
        Add
      </Button>
    </ConditonalTooltip>
  );
}
