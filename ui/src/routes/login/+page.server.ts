import { doRequest } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths } from "$lib/paths";
import { redirect, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import cookieParser from "set-cookie-parser";

export const load: PageServerLoad = async ({ locals }) => {
  if (locals.isLoggedIn) {
    redirect(302, paths.budgets.list);
  }
};

export const actions: Actions = {
  login: async ({ fetch, request, cookies }) => {
    const res = await doRequest({
      method: "POST",
      path: `/v1/user/login`,
      request,
      fetch,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    res.headers.getSetCookie().forEach((str) => {
      const { name, value, ...options } = cookieParser.parseString(str);
      const sameSite = options.sameSite?.toLowerCase();
      const path = options.path;

      if (path && (sameSite === "strict" || sameSite === "lax")) {
        cookies.set(name, value, { ...options, path, sameSite });
      }
    });

    redirect(302, paths.budgets.list);
  },
};
