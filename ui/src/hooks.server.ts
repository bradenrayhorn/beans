import { api } from "$lib/api/api";
import { getError } from "$lib/api/fetch-error";
import { paths } from "$lib/paths";
import { redirect, type Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
  if (event.route.id) {
    const res = await event.fetch(api(`/v1/user/me`));
    let isLoggedIn = true;

    if (res.status === 401) {
      isLoggedIn = false;
    } else if (!res.ok) {
      await getError(res);
    }

    if (res.ok) {
      event.locals.maybeUser = await res.json();
    }
    event.locals.isLoggedIn = isLoggedIn;

    if (event.route.id?.startsWith("/budget") && !isLoggedIn) {
      redirect(302, paths.login);
    }
  }

  return await resolve(event);
};
