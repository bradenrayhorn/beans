import { error, fail } from "@sveltejs/kit";

const defaultError = "Unknown error";

export const getError = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;

  throw error(res.status, msg);
};

export const getErrorForAction = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;

  if (res.status >= 400 && res.status < 500) {
    return fail(res.status, { message: msg });
  }

  throw error(res.status, msg);
};
