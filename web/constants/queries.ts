import ky from "ky";
import { KyInstance } from "ky/distribution/types/ky";
import { User } from "constants/types";

const queryKeys = {
  login: "login",
  me: "me",
};

const buildQueries = (client: KyInstance) => {
  return {
    login: ({ username, password }: { username: string; password: string }) =>
      client.post("api/v1/user/login", { json: { username, password } }),

    me: ({ cookie }: { cookie?: string } = {}): Promise<User> =>
      client.get("api/v1/user/me", { headers: { cookie } }).json(),
  };
};

export const queries = buildQueries(ky.extend({ prefixUrl: "/" }));

export { buildQueries, queryKeys };
