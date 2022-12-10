import {
  CreateMonthResponse,
  queryKeys,
  useQueries,
} from "@/constants/queries";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { useBudgetID } from "./budget";

export const useMonth = ({ monthID }: { monthID: string }) => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const query = useQuery([queryKeys.months.get, budgetID, monthID], () =>
    queries.months.get({ monthID })
  );

  return { ...query, month: query.data?.data };
};

export const useCreateMonth = ({
  onError,
  onSuccess,
}: {
  onError: () => void;
  onSuccess: (data: CreateMonthResponse) => void;
}) => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.months.create, { onError, onSuccess });

  const submit = useCallback(
    ({ date }: { date: string }) => mutation.mutateAsync({ date }),
    []
  );

  return { ...mutation, submit };
};

export const useUpdateMonthCategory = ({
  monthID,
  categoryID,
}: {
  monthID: string;
  categoryID: string;
}) => {
  const queryClient = useQueryClient();
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.months.categories.update);

  const submit = useCallback(
    ({ amount }: { amount: string }) =>
      mutation
        .mutateAsync({
          amount: !!amount && +amount !== 0 ? +amount.replace(/,/g, "") : null,
          monthID,
          categoryID,
        })
        .then(() => {
          queryClient.invalidateQueries([queryKeys.months.get]);
        }),
    [monthID, categoryID]
  );

  return { ...mutation, submit };
};
