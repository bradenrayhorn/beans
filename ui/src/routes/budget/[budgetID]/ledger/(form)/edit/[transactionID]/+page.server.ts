import { doAction } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths, withParameter } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";

export const actions: Actions = {
  save: async ({ fetch, request, params }) => {
    const res = await doAction({
      method: "PUT",
      path: `/api/v1/transactions/${params.transactionID}`,
      request,
      fetch,
      params,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    redirect(302, withParameter(paths.budget.ledger.base, params));
  },

  delete: async ({ fetch, request, params }) => {
    const res = await doAction({
      method: "POST",
      path: `/api/v1/transactions/delete`,
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
