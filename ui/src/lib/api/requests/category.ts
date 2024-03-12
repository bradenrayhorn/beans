import type { CategoryGroup } from "$lib/types/category";
import { api } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { withBudgetHeader, type WithBudgetID } from "./with-budget-header";

export const getCategoryGroups = async ({
  fetch: _fetch,
  params,
}: WithFetch & WithBudgetID): Promise<Array<CategoryGroup>> => {
  const res = await _fetch(api("/v1/categories"), withBudgetHeader(params));

  if (!res.ok) {
    return await getError(res);
  }

  return await res
    .json()
    .then((json: DataWrapped<Array<CategoryGroup>>) => json.data);
};
