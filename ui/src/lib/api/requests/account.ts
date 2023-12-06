import type { Account } from "$lib/types/account";
import { Amount, type APIAmount } from "$lib/types/amount";
import { api } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { withBudgetHeader, type WithBudgetID } from "./with-budget-header";

type APIAccount = {
  id: string;
  name: string;
  balance: APIAmount;
};

export const getAccounts = async ({
  fetch: _fetch,
  params,
}: WithFetch & WithBudgetID): Promise<Array<Account>> => {
  const res = await _fetch(api("/v1/accounts"), withBudgetHeader(params));

  if (!res.ok) {
    return await getError(res);
  }

  return await res.json().then((json: DataWrapped<Array<APIAccount>>) =>
    json.data.map((account) => ({
      id: account.id,
      name: account.name,
      balance: new Amount(account.balance),
    })),
  );
};
