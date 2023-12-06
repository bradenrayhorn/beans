import type { Budget } from "$lib/types/budget";
import { api } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";

export const getBudgets = async ({
  fetch: _fetch,
}: WithFetch): Promise<Array<Budget>> => {
  const res = await _fetch(api("/v1/budgets"));

  if (!res.ok) {
    return await getError(res);
  }

  return await res.json().then((json: DataWrapped<Array<Budget>>) => {
    return json.data;
  });
};

export const getBudget = async ({
  fetch: _fetch,
  budgetID,
}: WithFetch & { budgetID: string }): Promise<Budget> => {
  const res = await _fetch(api(`/v1/budgets/${budgetID}`));

  if (!res.ok) {
    return await getError(res);
  }

  return await res.json().then((json: DataWrapped<Budget>) => json.data);
};
