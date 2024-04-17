import { submitForm } from "$lib/api/api";
import { getErrorForForm } from "$lib/api/fetch-error";
import { paths, withParameter } from "$lib/paths";
import { fail, redirect, type Actions } from "@sveltejs/kit";
import { message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { schema } from "./schema";

export const load = async () => {
  const form = await superValidate(zod(schema));
  return { form };
};

export const actions: Actions = {
  save: async ({ fetch, request, params }) => {
    const form = await superValidate(request, zod(schema));

    if (!form.valid) {
      return fail(400, { form });
    }

    const res = await submitForm({
      method: "POST",
      path: `/api/v1/accounts`,
      body: form.data,
      fetch,
      params,
    });

    if (!res.ok) {
      const { status, error } = await getErrorForForm(res);
      return message(form, error, { status });
    }

    redirect(302, withParameter(paths.budget.accounts.base, params));
  },
};
