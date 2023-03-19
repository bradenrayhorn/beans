import { queryKeys, useQueries } from "@/constants/queries";
import { useMonthDate } from "@/context/MonthProvider";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { useBudgetID } from "./budget";

export const useMonth = () => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });
  const monthDate = useMonthDate();

  const query = useQuery([queryKeys.months.get, budgetID, monthDate], () =>
    queries.months.get({ date: monthDate })
  );

  return {
    ...query,
    month: query.data?.data,
  };
};

export const useUpdateMonth = ({ monthID }: { monthID: string }) => {
  const queryClient = useQueryClient();
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.months.update);

  const submit = useCallback(
    ({ carryover }: { carryover: string }) =>
      mutation
        .mutateAsync({
          carryover: +carryover.replace(/,/g, ""),
          monthID,
        })
        .then(() => {
          queryClient.invalidateQueries([queryKeys.months.get]);
        }),
    [mutation, monthID, queryClient]
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
    [mutation, monthID, categoryID, queryClient]
  );

  return { ...mutation, submit };
};
