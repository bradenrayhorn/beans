import {
  CreateMonthResponse,
  queryKeys,
  useQueries,
} from "@/constants/queries";
import { useMonthID } from "@/context/MonthProvider";
import {
  useIsMutating,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { useCallback } from "react";
import { useBudgetID } from "./budget";

export const useMonth = ({
  monthID: propMonthID,
}: { monthID?: string } = {}) => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });
  const ctxMonthID = useMonthID();

  const isMutating = useIsMutating({ mutationKey: [queryKeys.months.create] });

  const monthID = propMonthID ?? ctxMonthID;
  const query = useQuery([queryKeys.months.get, budgetID, monthID], () =>
    queries.months.get({ monthID })
  );

  return {
    ...query,
    isLoading: !!query.isLoading || !!isMutating,
    month: query.data?.data,
  };
};

export const useCreateMonth = ({
  onError,
  onSuccess,
}: {
  onError?: () => void;
  onSuccess: (data: CreateMonthResponse) => void;
}) => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.months.create, {
    onError,
    onSuccess,
    mutationKey: [queryKeys.months.create],
  });

  const submit = useCallback(
    ({ date }: { date: string }) => mutation.mutateAsync({ date }),
    [mutation]
  );

  return { ...mutation, submit };
};

export const useIsMonthLoading = (): boolean => {
  const monthID = useMonthID();
  const isMutating = useIsMutating({ mutationKey: [queryKeys.months.create] });

  const { isLoading: isMonthLoading } = useMonth({ monthID });

  return !!isMutating || isMonthLoading;
};

export const useUpdateMonth = () => {
  const queryClient = useQueryClient();
  const budgetID = useBudgetID();
  const monthID = useMonthID();
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
