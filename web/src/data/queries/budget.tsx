import { queries, queryKeys } from "@/constants/queries";
import { Budget } from "@/constants/types";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

export const useBudgetID = () => {
  const { budget: budgetID } = useParams();

  return budgetID as string;
};

export const useBudget = () => {
  const budgetID = useBudgetID();

  const query = useQuery([queryKeys.budget.get, budgetID], () =>
    queries.budget.get({ budgetID: (budgetID ?? "") as string })
  );

  return { ...query, budget: query.data?.data as Budget };
};
