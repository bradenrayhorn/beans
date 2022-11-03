import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useBudgetID } from "components/layouts/BudgetLayout";
import { queryKeys, useQueries } from "constants/queries";
import { useCallback } from "react";

export const useCategories = () => {
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const query = useQuery(
    [queryKeys.categories.get, budgetID],
    queries.categories.get
  );

  return { ...query, categoryGroups: query.data?.data ?? [] };
};

interface AddCategoryGroupData {
  name: string;
}

export const useAddCategoryGroup = () => {
  const queryClient = useQueryClient();
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.categories.createGroup);

  const submit = useCallback(
    (values: AddCategoryGroupData) =>
      mutation.mutateAsync(values).then(() => {
        queryClient.invalidateQueries([queryKeys.categories.get]);
      }),
    []
  );

  return { ...mutation, submit };
};

interface AddCategoryData {
  name: string;
}

export const useAddCategory = ({ groupID }: { groupID: string }) => {
  const queryClient = useQueryClient();
  const budgetID = useBudgetID();
  const queries = useQueries({ budgetID });

  const mutation = useMutation(queries.categories.createCategory);

  const submit = useCallback(
    (values: AddCategoryData) =>
      mutation.mutateAsync({ ...values, groupID }).then(() => {
        queryClient.invalidateQueries([queryKeys.categories.get]);
      }),
    [groupID]
  );

  return { ...mutation, submit };
};
