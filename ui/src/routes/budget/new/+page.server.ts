import { doRequest } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";

export const actions: Actions = {
  save: async ({ fetch, request }) => {
    const res = await doRequest({
      method: "POST",
      path: `/v1/budgets`,
      request,
      fetch,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    throw redirect(302, paths.budgets.list);
  },
};
