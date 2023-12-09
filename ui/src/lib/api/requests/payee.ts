import type { Payee } from "$lib/types/payee";
import { api } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { withBudgetHeader, type WithBudgetID } from "./with-budget-header";

export const getPayees = async ({
  fetch: _fetch,
  params,
}: WithFetch & WithBudgetID): Promise<Array<Payee>> => {
  const res = await _fetch(api("/v1/payees"), withBudgetHeader(params));

  if (!res.ok) {
    return await getError(res);
  }

  return await res.json().then((json: DataWrapped<Array<Payee>>) => json.data);
};
