import type { PageLoad } from "./$types";
import { getBudgets } from "$lib/api/requests/budget";

export const load: PageLoad = async ({ fetch }) => {
  const budgets = await getBudgets({ fetch });

  return { budgets };
};
