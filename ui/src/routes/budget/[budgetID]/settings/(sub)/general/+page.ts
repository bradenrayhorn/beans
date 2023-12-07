import { getBudget } from "$lib/api/requests/budget";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
  const budget = await getBudget({ fetch, budgetID: params.budgetID });

  return { budget };
};
