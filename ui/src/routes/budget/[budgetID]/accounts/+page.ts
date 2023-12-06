import { getAccounts } from "$lib/api/requests/account";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
  const accounts = await getAccounts({ fetch, params });

  return { accounts };
};
