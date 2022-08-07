import { useQuery } from "@tanstack/react-query";
import { queries, queryKeys } from "constants/queries";
import { Budget } from "constants/types";
import { useRouter } from "next/router";

export const useBudget = () => {
  const router = useRouter();
  const { budget: budgetID } = router.query;

  const query = useQuery([queryKeys.budget.get, budgetID], () =>
    queries.budget.get({ budgetID: (budgetID ?? "") as string })
  );

  return { ...query, budget: query.data?.data as Budget };
};
