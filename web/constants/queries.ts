import ky, { HTTPError } from "ky";
import { KyInstance } from "ky/distribution/types/ky";
import { Account, Budget, User } from "constants/types";

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
    login: ({ username, password }: { username: string; password: string }) =>
      client.post("api/v1/user/login", { json: { username, password } }),

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

    accounts: {
      get: ({ budgetID }: { budgetID: string }) =>
        client
          .get(`api/v1/budgets/${budgetID}/accounts`)
          .json<GetAccountsResponse>(),

      create: ({ budgetID, name }: { budgetID: string; name: string }) =>
        client.post(`api/v1/budgets/${budgetID}/accounts`, { json: { name } }),
    },
  };
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
