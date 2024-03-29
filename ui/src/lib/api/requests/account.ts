import type { Account, AccountWithBalance } from "$lib/types/account";
import { Amount, type APIAmount } from "$lib/types/amount";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { withBudgetHeader, type WithBudgetID } from "./with-budget-header";

type APIAccountWithBalance = {
  id: string;
  name: string;
  balance: APIAmount;
  offBudget: boolean;
};

export const getAccounts = async ({
  fetch: _fetch,
  params,
}: WithFetch & WithBudgetID): Promise<Array<AccountWithBalance>> => {
  const res = await _fetch("/api/v1/accounts", withBudgetHeader(params));

  if (!res.ok) {
    return await getError(res);
  }

  return await res
    .json()
    .then((json: DataWrapped<Array<APIAccountWithBalance>>) =>
      json.data.map((account) => ({
        id: account.id,
        name: account.name,
        balance: new Amount(account.balance),
        offBudget: account.offBudget,
      })),
    );
};

export const getTransactableAccounts = async ({
  fetch: _fetch,
  params,
}: WithFetch & WithBudgetID): Promise<Array<Account>> => {
  const res = await _fetch(
    "/api/v1/accounts/transactable",
    withBudgetHeader(params),
  );

  if (!res.ok) {
    return await getError(res);
  }

  return await res
    .json()
    .then((json: DataWrapped<Array<Account>>) => json.data);
};
