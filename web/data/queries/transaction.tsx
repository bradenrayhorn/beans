import { useMutation } from "@tanstack/react-query";
import { useBudgetID } from "components/layouts/BudgetLayout";
import { getHTTPErrorResponseMessage, useQueries } from "constants/queries";
import { useCallback } from "react";

interface AddTransactionData {
  accountID: string;
  amount: string;
  date: string;
  notes?: string;
}

export const useAddTransaction = () => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.transactions.create);
  const errorMessage = getHTTPErrorResponseMessage(mutation.error);

  const submit = useCallback(
    (values: AddTransactionData) => mutation.mutateAsync(values),
    [budgetID]
  );

  return { ...mutation, errorMessage, submit };
};
