import { getAccounts } from "$lib/api/requests/account";
import { getCategoryGroups } from "$lib/api/requests/category";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch, params }) => {
  const categoryGroups = await getCategoryGroups({ fetch, params });
  const categories = categoryGroups.flatMap((group) => group.categories);

  const accounts = await getAccounts({ fetch, params });

  return { categoryGroups, categories, accounts };
};
