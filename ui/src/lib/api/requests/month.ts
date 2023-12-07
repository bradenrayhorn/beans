import { Amount, type APIAmount } from "$lib/types/amount";
import type { Month, MonthCategory } from "$lib/types/month";
import { api } from "../api";
import { getError } from "../fetch-error";
import type { DataWrapped } from "./data-wrapped";
import type { WithFetch } from "./fetch";
import { withBudgetHeader, type WithBudgetID } from "./with-budget-header";

type APIMonth = {
  id: string;
  date: string;
  budgetable: APIAmount;
  carryover: APIAmount;
  income: APIAmount;
  assigned: APIAmount;
  carried_over: APIAmount;
  categories: Array<APIMonthCategory>;
};

type APIMonthCategory = {
  id: string;
  category_id: string;
  assigned: APIAmount;
  activity: APIAmount;
  available: APIAmount;
};

export const getMonth = async ({
  fetch: _fetch,
  date,
  params,
}: WithFetch & WithBudgetID & { date: string }): Promise<Month> => {
  const res = await _fetch(api(`/v1/months/${date}`), withBudgetHeader(params));

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
    carriedOver: new Amount(json.carried_over),
    categories: json.categories.map<MonthCategory>((apiCategory) => ({
      id: apiCategory.id,
      categoryID: apiCategory.category_id,
      assigned: new Amount(apiCategory.assigned),
      activity: new Amount(apiCategory.activity),
      available: new Amount(apiCategory.available),
    })),
  }));
};
