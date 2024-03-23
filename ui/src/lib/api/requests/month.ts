import { Amount, type APIAmount } from "$lib/types/amount";
import type { Month, MonthCategory } from "$lib/types/month";
import { doRequest } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { type WithBudgetID } from "./with-budget-header";

type APIMonth = {
  id: string;
  date: string;
  budgetable: APIAmount;
  carryover: APIAmount;
  income: APIAmount;
  assigned: APIAmount;
  carriedOver: APIAmount;
  categories: Array<APIMonthCategory>;
};

type APIMonthCategory = {
  id: string;
  categoryId: string;
  assigned: APIAmount;
  activity: APIAmount;
  available: APIAmount;
};

export const getMonth = async ({
  fetch: _fetch,
  date,
  params,
}: WithFetch & WithBudgetID & { date: string }): Promise<Month> => {
  const res = await doRequest({
    method: "GET",
    path: `/v1/months/${date}`,
    fetch: _fetch,
    params,
  });

  if (!res.ok) {
    return await getError(res);
  }

  return await res.json().then(({ data: json }: DataWrapped<APIMonth>) => ({
    id: json.id,
    date: json.date,
    budgetable: new Amount(json.budgetable),
    carryover: new Amount(json.carryover),
    income: new Amount(json.income),
    assigned: new Amount(json.assigned),
    carriedOver: new Amount(json.carriedOver),
    categories: json.categories.map<MonthCategory>((apiCategory) => ({
      id: apiCategory.id,
      categoryID: apiCategory.categoryId,
      assigned: new Amount(apiCategory.assigned),
      activity: new Amount(apiCategory.activity),
      available: new Amount(apiCategory.available),
    })),
  }));
};
