import { paths } from "$lib/paths";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  if (!locals.isLoggedIn) {
    throw redirect(302, paths.login);
  } else {
    throw redirect(302, paths.budgets.list);
  }
};
