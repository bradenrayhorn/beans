import type { RequestHandler } from "@sveltejs/kit";
import { env as privateEnv } from "$env/dynamic/private";

export const proxyToServer: RequestHandler = ({ request, fetch, cookies }) => {
  const url = new URL(request.url);
  request = new Request(
    `${privateEnv.UNPROXIED_SERVER_URL ?? "http://localhost:8000"}${url.pathname}${url.search}`,
    request,
  );

  request.headers.set("Authorization", cookies.get("si", {}) ?? "");

  return fetch(request);
};
