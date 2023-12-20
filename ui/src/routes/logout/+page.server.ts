import { paths } from "$lib/paths";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { logout } from "$lib/api/requests/user";

export const load: PageServerLoad = async ({ locals, fetch }) => {
  if (!locals.isLoggedIn) {
    redirect(302, paths.login);
  }

  await logout({ fetch });
};
