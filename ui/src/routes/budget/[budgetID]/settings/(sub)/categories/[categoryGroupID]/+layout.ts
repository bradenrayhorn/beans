import { getCategoryGroups } from "$lib/api/requests/category";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch, params }) => {
  const categoryGroups = await getCategoryGroups({ fetch, params });

  return { categoryGroups };
};
