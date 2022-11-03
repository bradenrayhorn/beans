import {
  Account,
  Budget,
  CategoryGroup,
  Transaction,
  User,
} from "constants/types";
import ky, { HTTPError } from "ky";
import { KyInstance } from "ky/distribution/types/ky";
import { useEffect, useState } from "react";

const queryKeys = {
  accounts: {
    get: "accounts_get",
  },

  budget: {
    get: "budget_get",
    getAll: "budget_get_all",
  },

  categories: {
    get: "categories_get",
    addGroup: "categories_add_group",
  },

  login: "login",

  me: "me",

  transactions: {
    getAll: "transactions_get_all",
  },
};

interface GetAccountsResponse {
  data: Account[];
}
interface GetAllBudgetsResponse {
  data: Budget[];
}
interface GetBudgetResponse {
  data: Budget;
}
interface GetTransactionsResponse {
  data: Transaction[];
}
interface GetCategoriesResponse {
  data: CategoryGroup[];
}

const buildQueries = (client: KyInstance) => {
  client = client.extend({
    hooks: {
      beforeError: [
        async (error) => {
          if (error.response) {
            try {
              const errorJSON = await error.response.json();
              error.message = (errorJSON as { error: string }).error;
            } catch (e) {
              console.error("failed to parse error response");
            }
          }
          return error;
        },
      ],
    },
  });

  return {
    // user
    login: ({
      username,
      password,
    }: {
      username: string;
      password: string;
    }): Promise<User> =>
      client.post("api/v1/user/login", { json: { username, password } }).json(),

    me: ({ cookie }: { cookie?: string } = {}): Promise<User> =>
      client.get("api/v1/user/me", { headers: { cookie } }).json(),

    // accounts
    accounts: {
      get: () => client.get(`api/v1/accounts`).json<GetAccountsResponse>(),

      create: ({ name }: { name: string }) =>
        client.post(`api/v1/accounts`, { json: { name } }),
    },

    categories: {
      get: () => client.get(`api/v1/categories`).json<GetCategoriesResponse>(),
      createCategory: ({ name, groupID }: { name: string; groupID: string }) =>
        client.post(`api/v1/categories`, {
          json: { name, group_id: groupID },
        }),
      createGroup: ({ name }: { name: string }) =>
        client.post(`api/v1/categories/groups`, { json: { name } }),
    },

    // budget
    budget: {
      get: ({ budgetID }: { budgetID: string }) =>
        client.get(`api/v1/budgets/${budgetID}`).json<GetBudgetResponse>(),

      getAll: () => client.get("api/v1/budgets").json<GetAllBudgetsResponse>(),

      create: ({ name }: { name: string }) =>
        client.post("api/v1/budgets", { json: { name } }),
    },

    // transactions
    transactions: {
      create: ({
        accountID,
        amount,
        date,
        notes,
      }: {
        accountID: string;
        amount: string;
        date: string;
        notes?: string;
      }) =>
        client.post(`api/v1/transactions`, {
          json: { account_id: accountID, amount, date, notes },
        }),

      getAll: () =>
        client.get(`api/v1/transactions`).json<GetTransactionsResponse>(),
    },
  };
};

type Props = {
  budgetID?: string;
};
const getQueries = ({ budgetID }: Props) => {
  const client = ky.extend({
    prefixUrl: "/",
    hooks: {
      beforeRequest: [
        (request) => {
          if (budgetID) {
            request.headers.set("Budget-ID", budgetID);
          }
        },
      ],
    },
  });
  return buildQueries(client);
};

export const useQueries = ({ budgetID }: Props) => {
  const [queries, setQueries] = useState(() => getQueries({ budgetID }));

  useEffect(() => {
    setQueries(getQueries({ budgetID }));
  }, [budgetID]);

  return queries;
};

export const queries = buildQueries(ky.extend({ prefixUrl: "/" }));

export function getHTTPErrorResponseMessage(error: unknown) {
  if (!error) {
    return "";
  }
  if (error instanceof HTTPError) {
    return error.message;
  }

  return "Unknown error.";
}

export { buildQueries, queryKeys };
