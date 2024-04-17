import { error, fail, type NumericRange } from "@sveltejs/kit";
import type { ErrorStatus } from "sveltekit-superforms";

const defaultError = "Unknown error";

export const getError = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;

  if (res.status >= 400 && res.status <= 599) {
    error(res.status as NumericRange<400, 599>, msg);
  }

  error(500, msg);
};

export const getErrorForAction = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;

  if (res.status >= 400 && res.status < 500) {
    return fail(res.status as NumericRange<400, 599>, { message: msg });
  } else if (res.status <= 599) {
    error(res.status as NumericRange<400, 599>, msg);
  } else {
    error(500, msg);
  }
};

export const getErrorForForm = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;

  return { status: res.status as ErrorStatus, error: msg };
};
