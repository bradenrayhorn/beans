import { doAction } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths, withParameter } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";

export const actions: Actions = {
  save: async ({ fetch, request, params }) => {
    const res = await doAction({
      method: "POST",
      path: `/api/v1/categories/groups`,
      request,
      fetch,
      params,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    redirect(302, withParameter(paths.budget.settings.categories.base, params));
  },
};
