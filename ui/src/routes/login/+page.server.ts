import { doRequest } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  if (locals.isLoggedIn) {
    throw redirect(302, paths.budgets.list);
  }
};

export const actions: Actions = {
  login: async ({ fetch, request }) => {
    const res = await doRequest({
      method: "POST",
      path: `/v1/user/login`,
      request,
      fetch,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    throw redirect(302, paths.budgets.list);
  },
};
