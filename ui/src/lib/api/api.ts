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
  mapFormData,
}: {
  method: string;
  path: string;
  request?: Request;
  fetch: typeof fetch;
  params?: { [key: string]: string | undefined };
  mapFormData?: (obj: { [key: string]: unknown }) => { [key: string]: unknown };
}): Promise<Response> => {
  let obj: null | { [key: string]: unknown } = null;

  if (request) {
    obj = {};
    const data = await request.formData();
    data.forEach((value, key) => {
      if (obj) {
        if (key.endsWith("[]")) {
          obj[key.slice(0, key.length - 2)] = data.getAll(key);
        } else if (key.endsWith("[json]")) {
          obj[key.slice(0, key.length - 6)] = data
            .getAll(key)
            .map((value) => JSON.parse(value.toString()));
        } else {
          obj[key] = value.toString();
        }
      }
    });

    if (mapFormData) {
      obj = mapFormData(obj);
    }
  }

  const res = await internalFetch(api(path), {
    method,
    body: obj ? JSON.stringify(obj) : null,
    headers: { "Budget-ID": params?.budgetID ?? "" },
  });

  return res;
};
