import { doAction } from "$lib/api/api";
import { getErrorForAction } from "$lib/api/fetch-error";
import type { Actions } from "@sveltejs/kit";

export const actions: Actions = {
  save: async ({ fetch, request, params }) => {
    const data = await request.clone().formData();
    const res = await doAction({
      method: "POST",
      path: `/api/v1/months/${data.get("monthID")}/categories`,
      request,
      fetch,
      params,
    });

    if (!res.ok) {
      return await getErrorForAction(res);
    }

    return {};
  },
};
