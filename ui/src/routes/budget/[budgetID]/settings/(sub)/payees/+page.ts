import { getPayees } from "$lib/api/requests/payee";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
  const payees = await getPayees({ fetch, params });

  return { payees };
};
