import { Checkbox } from "@chakra-ui/react";
import { useTransactionTableState } from "./TransactionList";

export default function TransactionRowCheckbox({ id }: { id: string }) {
  const isSelected = useTransactionTableState((state) =>
    state.isRowSelected(id)
  );

  const set = useTransactionTableState((state) => state.setRowSelection);

  return (
    <Checkbox
      aria-label="Select transaction"
      isChecked={isSelected}
      onChange={(e) => {
        set(id, e.target.checked);
      }}
    />
  );
}
