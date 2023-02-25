import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getHTTPErrorResponseMessage,
  queryKeys,
  useQueries,
} from "@/constants/queries";
import { useCallback } from "react";
import { useBudgetID } from "./budget";

interface TransactionData {
  accountID: string;
  categoryID?: string;
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
    (values: TransactionData) =>
      mutation
        .mutateAsync({
          ...values,
          amount: +values.amount.replace(/,/g, ""),
        })
        .then(() => {
          queryClient.invalidateQueries([queryKeys.transactions.getAll]);
        }),
    [mutation, queryClient]
  );

  return { ...mutation, errorMessage, submit };
};

export const useEditTransaction = ({ id }: { id: string }) => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });
  const queryClient = useQueryClient();

  const mutation = useMutation(queries.transactions.update);

  const submit = useCallback(
    (values: TransactionData) =>
      mutation
        .mutateAsync({
          ...values,
          amount: +values.amount.replace(/,/g, ""),
          id,
        })
        .then(() => {
          queryClient.invalidateQueries([queryKeys.transactions.getAll]);
        }),
    [id, mutation, queryClient]
  );

  return { ...mutation, submit };
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
