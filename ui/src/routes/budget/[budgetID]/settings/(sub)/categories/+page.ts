import { getCategoryGroups } from "$lib/api/requests/category";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
  const categoryGroups = await getCategoryGroups({ fetch, params });

  return { categoryGroups };
};
