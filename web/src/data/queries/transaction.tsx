import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getHTTPErrorResponseMessage,
  queryKeys,
  useQueries,
} from "@/constants/queries";
import { useCallback } from "react";
import { useBudgetID } from "./budget";

interface AddTransactionData {
  accountID: string;
  amount: string;
  date: string;
  notes?: string;
}

export const useAddTransaction = () => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });
  const queryClient = useQueryClient();

  const mutation = useMutation(queries.transactions.create);
  const errorMessage = getHTTPErrorResponseMessage(mutation.error);

  const submit = useCallback(
    (values: AddTransactionData) =>
      mutation.mutateAsync(values).then(() => {
        queryClient.invalidateQueries([queryKeys.transactions.getAll]);
      }),
    [budgetID]
  );

  return { ...mutation, errorMessage, submit };
};

export const useTransactions = () => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const query = useQuery(
    [queryKeys.transactions.getAll, budgetID],
    queries.transactions.getAll
  );

  return { ...query, transactions: query.data?.data ?? [] };
};
