import { doAction } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import { paths } from "$lib/paths";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, fetch, cookies }) => {
  if (!locals.isLoggedIn) {
    redirect(302, paths.login);
  }

  const res = await doAction({
    method: "POST",
    path: "/api/v1/user/logout",
    fetch,
  });

  if (!res.ok) {
    return await getErrorForAction(res);
  }

  cookies.delete("si", { path: "/" });

  redirect(302, paths.login);
};
