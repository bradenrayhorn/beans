export const withBudgetHeader = (params: { budgetID: string }) => ({
  headers: { "Budget-ID": params.budgetID },
});

export type WithBudgetID = {
  params: { budgetID: string };
};
