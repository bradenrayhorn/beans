import { getCategoryGroups } from "$lib/api/requests/category";
import { getMonth } from "$lib/api/requests/month";
import type { LayoutLoad } from "./$types";
import dayjs from "dayjs";

export const load: LayoutLoad = async ({ fetch, params }) => {
  const [month, categoryGroups] = await Promise.all([
    getMonth({
      fetch,
      date: dayjs(params.month).format("YYYY-MM-DD"),
      params,
    }),

    getCategoryGroups({
      fetch,
      params,
    }),
  ]);

  return { month, categoryGroups };
};
