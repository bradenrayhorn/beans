import { paths } from "$lib/paths";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  if (!locals.isLoggedIn) {
    redirect(302, paths.login);
  } else {
    redirect(302, paths.budgets.list);
  }
};
