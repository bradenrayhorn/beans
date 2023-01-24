import { Transaction } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import { useEditTransaction } from "@/data/queries/transaction";
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
  const { submit } = useEditTransaction({ id: transaction.id });

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
      onSubmit={(values) =>
        submit({
          accountID: values.account.id,
          categoryID: values.category?.id,
          amount: values.amount,
          date: values.date,
          notes: values.notes.trim() ? values.notes : undefined,
        })
          .then(() => {
            onSuccess();
          })
          .catch((error) => {
            console.error(error);
          })
      }
    >
      {children}
    </TransactionFormProvider>
  );
}
