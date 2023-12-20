import { paths, withParameter } from "$lib/paths";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import dayjs from "dayjs";

export const load: PageLoad = async ({ params }) => {
  redirect(
    302,
    withParameter(paths.budget.budget.month, {
      ...params,
      month: dayjs().format("YYYY-MM"),
    }),
  );
};
