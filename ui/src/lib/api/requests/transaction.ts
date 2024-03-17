import { Amount, type APIAmount } from "$lib/types/amount";
import type { Transaction } from "$lib/types/transaction";
import { api } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { withBudgetHeader, type WithBudgetID } from "./with-budget-header";

type APITransaction = Omit<Transaction, "amount"> & {
  amount: APIAmount;
};

const mapTransaction = (transaction: APITransaction): Transaction => ({
  ...transaction,
  amount: new Amount(transaction.amount),
});

export const getTransactions = async ({
  fetch: _fetch,
  params,
}: WithFetch & WithBudgetID): Promise<Array<Transaction>> => {
  const res = await _fetch(api("/v1/transactions"), withBudgetHeader(params));

  if (!res.ok) {
    return await getError(res);
  }

  return await res
    .json()
    .then((json: DataWrapped<Array<APITransaction>>) =>
      json.data.map(mapTransaction),
    );
};

export const getTransaction = async ({
  fetch: _fetch,
  id,
  params,
}: WithFetch & WithBudgetID & { id: string }): Promise<Transaction> => {
  const res = await _fetch(
    api(`/v1/transactions/${id}`),
    withBudgetHeader(params),
  );

  if (!res.ok) {
    return await getError(res);
  }

  return await res
    .json()
    .then((json: DataWrapped<APITransaction>) => mapTransaction(json.data));
};
