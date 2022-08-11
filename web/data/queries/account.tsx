import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useBudgetID } from "components/layouts/BudgetLayout";
import {
  getHTTPErrorResponseMessage,
  queryKeys,
  useQueries,
} from "constants/queries";
import { useCallback } from "react";

export const useAccounts = () => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const query = useQuery(
    [queryKeys.accounts.get, budgetID],
    queries.accounts.get
  );

  return { ...query, accounts: query.data?.data ?? [] };
};

interface AddAccountData {
  name: string;
}

export const useAddAccount = () => {
  const queryClient = useQueryClient();
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.accounts.create);
  const errorMessage = getHTTPErrorResponseMessage(mutation.error);

  const submit = useCallback(
    (values: AddAccountData) =>
      mutation.mutateAsync(values).then(() => {
        queryClient.invalidateQueries([queryKeys.accounts.get]);
      }),
    [budgetID]
  );

  return { ...mutation, errorMessage, submit };
};
