import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { paths, withParameter } from "$lib/paths";

export const load: PageServerLoad = async ({ params }) => {
  throw redirect(302, withParameter(paths.budget.budget.base, params));
};
