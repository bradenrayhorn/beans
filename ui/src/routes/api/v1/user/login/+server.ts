import type { DataWrapped } from "$lib/api/requests/data-wrapped";
import { proxyToServer } from "$lib/api/server";
import type { RequestHandler } from "@sveltejs/kit";
import dayjs from "dayjs";

export const POST: RequestHandler = async ({ cookies, ...rest }) => {
  const result = await proxyToServer({ cookies, ...rest });

  const headers: HeadersInit = {};
  if (result.ok) {
    const response = result.clone();
    const sessionID = await response
      .json()
      .then((json: DataWrapped<{ sessionID: string }>) => json.data.sessionID);

    headers["set-cookie"] = cookies.serialize("si", sessionID, {
      path: "/",
      httpOnly: true,
      sameSite: "strict",
      expires: dayjs().add(1, "month").toDate(),
    });
  }

  return new Response(result.body, {
    status: result.status,
    statusText: result.statusText,
    headers,
  });
};
