import { getSplits, getTransaction } from "$lib/api/requests/transaction";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
  const transaction = await getTransaction({
    id: params.transactionID,
    fetch,
    params,
  });
  const splits = await getSplits({
    id: params.transactionID,
    fetch,
    params,
  });

  return { transaction, splits };
};
