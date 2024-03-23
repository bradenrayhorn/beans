import { api } from "$lib/api/api";
import { getError } from "$lib/api/fetch-error";
import { paths } from "$lib/paths";
import { redirect, type Handle, type HandleFetch } from "@sveltejs/kit";
import { env as privateEnv } from "$env/dynamic/private";

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

export const handleFetch: HandleFetch = ({ event, request, fetch }) => {
  const url = new URL(request.url);

  if (privateEnv.UNPROXIED_BASE_API_URL) {
    request = new Request(
      `${privateEnv.UNPROXIED_BASE_API_URL ?? ""}${url.pathname}${url.search}`,
      request,
    );

    request.headers.set("cookie", event.request.headers.get("cookie") ?? "");
  }

  return fetch(request);
};
