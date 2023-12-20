import { doRequest } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths, withParameter } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";

export const actions: Actions = {
  save: async ({ fetch, request, params }) => {
    const res = await doRequest({
      method: "POST",
      path: `/v1/transactions`,
      request,
      fetch,
      params,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    redirect(302, withParameter(paths.budget.ledger.base, params));
  },
};
