import { env } from "$env/dynamic/public";

export function api(path: string): string {
  return `${env.PUBLIC_BASE_API_URL ?? ""}/api${path}`;
}

export const doRequest = async ({
  method,
  path,
  fetch: internalFetch,
  request,
  params,
}: {
  method: string;
  path: string;
  request?: Request;
  fetch: typeof fetch;
  params?: { [key: string]: string | undefined };
}): Promise<Response> => {
  let obj: null | { [key: string]: unknown } = null;

  if (request) {
    obj = {};
    const data = await request.formData();
    data.forEach((value, key) => {
      if (obj) {
        if (key.endsWith("[]")) {
          obj[key.slice(0, key.length - 2)] = data.getAll(key);
        } else {
          obj[key] = value.toString();
        }
      }
    });
  }

  const res = await internalFetch(api(path), {
    method,
    body: obj ? JSON.stringify(obj) : null,
    headers: { "Budget-ID": params?.budgetID ?? "" },
  });

  return res;
};
