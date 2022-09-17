import ky, { HTTPError } from "ky";
import { KyInstance } from "ky/distribution/types/ky";
import { Account, Budget, User } from "constants/types";
import { useEffect, useState } from "react";

const queryKeys = {
  login: "login",
  me: "me",
  budget: {
    get: "budget_get",
    getAll: "budget_get_all",
  },
  accounts: {
    get: "accounts_get",
  },
};

interface GetAllBudgetsResponse {
  data: Budget[];
}

interface GetBudgetResponse {
  data: Budget;
}

interface GetAccountsResponse {
  data: Account[];
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

    // budget
    budget: {
      get: ({ budgetID }: { budgetID: string }) =>
        client.get(`api/v1/budgets/${budgetID}`).json<GetBudgetResponse>(),

      getAll: () => client.get("api/v1/budgets").json<GetAllBudgetsResponse>(),

      create: ({ name }: { name: string }) =>
        client.post("api/v1/budgets", { json: { name } }),
    },

    // accounts
    accounts: {
      get: () => client.get(`api/v1/accounts`).json<GetAccountsResponse>(),

      create: ({ name }: { name: string }) =>
        client.post(`api/v1/accounts`, { json: { name } }),
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
  const [queries, setQueries] = useState(getQueries({ budgetID }));

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
