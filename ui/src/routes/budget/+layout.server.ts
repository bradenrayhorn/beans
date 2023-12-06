import type { User } from "../../app";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
  const user: User = locals.maybeUser;

  return {
    user,
  };
};
