import { Transaction } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import { PropsWithChildren } from "react";
import { TransactionFormProvider } from "./TransactionForm";

export default function EditTransactionProvider({
  children,
  onSuccess,
  transaction,
  isEditing,
}: {
  transaction: Transaction;
  isEditing: boolean;
  onSuccess: () => void;
} & PropsWithChildren) {
  return (
    <TransactionFormProvider
      key={`${isEditing}`}
      defaultValues={{
        date: transaction.date,
        category: transaction.category,
        account: transaction.account,
        notes: transaction.notes ?? "",
        amount: formatAmount(transaction.amount),
      }}
      onSubmit={(values) => {
        console.log(values);
        onSuccess();
      }}
    >
      {children}
    </TransactionFormProvider>
  );
}
