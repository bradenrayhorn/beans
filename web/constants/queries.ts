import ky from "ky";

const queryKeys = {
  login: "login",
  me: 'me',
};

interface MeResponse {
  id: string,
  username: string,
}

const queries = {
  login: ({ username, password }: { username: string, password: string }) => ky.post('api/v1/user/login', { json: { username, password } }).json(),

  me: (): Promise<MeResponse> => ky.get('api/v1/user/me').json(),
};

export { queries, queryKeys };
