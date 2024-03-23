import { doAction } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths, withParameter } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";

export const actions: Actions = {
  save: async ({ fetch, request, params }) => {
    const res = await doAction({
      method: "POST",
      path: `/api/v1/accounts`,
      request,
      fetch,
      params,
      mapFormData: (obj) => {
        return { ...obj, off_budget: obj.off_budget === "true" };
      },
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    redirect(302, withParameter(paths.budget.accounts.base, params));
  },
};
