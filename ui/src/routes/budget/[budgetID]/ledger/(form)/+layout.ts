import { getTransactableAccounts } from "$lib/api/requests/account";
import { getCategoryGroups } from "$lib/api/requests/category";
import { getPayees } from "$lib/api/requests/payee";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch, params }) => {
  const categoryGroups = await getCategoryGroups({ fetch, params });
  const categories = categoryGroups.flatMap((group) => group.categories);

  const accounts = await getTransactableAccounts({ fetch, params });

  const payees = await getPayees({ fetch, params });

  return { categoryGroups, categories, accounts, payees };
};
