import { paths } from "$lib/paths";
import { redirect } from "@sveltejs/kit";
import { doRequest } from "../api";
import { getError } from "../fetch-error";
import type { WithFetch } from "./fetch";

export const logout = async ({ fetch: _fetch }: WithFetch): Promise<void> => {
  const res = await doRequest({
    method: "POST",
    path: "/v1/user/logout",
    fetch: _fetch,
  });

  if (!res.ok) {
    return await getError(res);
  }

  redirect(302, paths.login);
};
