import { currentDate } from "@/data/format/date";
import { useAccounts } from "@/data/queries/account";
import { useAddTransaction } from "@/data/queries/transaction";
import { Spinner } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { TransactionFormProvider } from "./TransactionForm";

export default function AddTransactionProvider({
  children,
  onSuccess,
}: { onSuccess: () => void } & PropsWithChildren) {
  const { submit } = useAddTransaction();

  const { accounts, isLoading } = useAccounts();

  if (isLoading) {
    return <Spinner />;
  }

  return (
    <TransactionFormProvider
      defaultValues={{
        date: currentDate(),
        category: null,
        account: accounts[0],
        notes: "",
        amount: "0",
      }}
      onSubmit={(values) =>
        submit({
          accountID: values.account.id,
          categoryID: values.category?.id,
          amount: values.amount,
          date: values.date,
          notes: values.notes ? values.notes : undefined,
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
