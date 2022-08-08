import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getHTTPErrorResponseMessage,
  queries,
  queryKeys,
} from "constants/queries";
import { useRouter } from "next/router";
import { useCallback } from "react";

export const useAccounts = () => {
  const router = useRouter();
  const { budget: budgetID } = router.query;

  const query = useQuery([queryKeys.accounts.get, budgetID], () =>
    queries.accounts.get({ budgetID: (budgetID ?? "") as string })
  );

  return { ...query, accounts: query.data?.data ?? [] };
};

interface AddAccountData {
  name: string;
}

export const useAddAccount = () => {
  const queryClient = useQueryClient();
  const router = useRouter();
  const { budget: budgetID } = router.query;

  const mutation = useMutation(queries.accounts.create);
  const errorMessage = getHTTPErrorResponseMessage(mutation.error);

  const submit = useCallback(
    (values: AddAccountData) =>
      mutation
        .mutateAsync({ ...values, budgetID: (budgetID ?? "") as string })
        .then(() => {
          queryClient.invalidateQueries([queryKeys.accounts.get]);
        }),
    [budgetID]
  );

  return { ...mutation, errorMessage, submit };
};
