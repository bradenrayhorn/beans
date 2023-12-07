export const paths = {
  login: "/login",
  logout: "/logout",

  budgets: {
    list: "/budget",
    new: "/budget/new",
  },

  budget: {
    budget: {
      base: "/budget/[budgetID]/budget",
      month: "/budget/[budgetID]/budget/[month]",
      category: "/budget/[budgetID]/budget/[month]/categories/[categoryID]",
      forNextMonth: "/budget/[budgetID]/budget/[month]/for-next-month",
    },
    accounts: {
      base: "/budget/[budgetID]/accounts",
      new: "/budget/[budgetID]/accounts/new",
    },
    ledger: {
      base: "/budget/[budgetID]/ledger",
      new: "/budget/[budgetID]/ledger/new",
      edit: "/budget/[budgetID]/ledger/edit/[transactionID]",
    },
    settings: {
      base: "/budget/[budgetID]/settings",
      general: "/budget/[budgetID]/settings/general",
      categories: {
        base: "/budget/[budgetID]/settings/categories",
        new: "/budget/[budgetID]/settings/categories/new",
        group: "/budget/[budgetID]/settings/categories/[categoryGroupID]",
        newSubGroup:
          "/budget/[budgetID]/settings/categories/[categoryGroupID]/new",
      },
      payees: "/budget/[budgetID]/settings/payees",
    },
  },
};

export const withParameter = (
  path: string,
  parameters: { [key: string]: string | undefined },
): string => {
  let newPath = path;
  Object.entries(parameters).forEach(([key, value]) => {
    newPath = newPath.replaceAll(`[${key}]`, value ?? "");
  });
  return newPath;
};
