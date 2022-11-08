import { expect, test as base } from "@playwright/test";

export type UserFixture = {
  username: string;
  password: string;
};

type Fixtures = {
  user: UserFixture;
};

export const test = base.extend<Fixtures>({
  user: async ({ request }, use) => {
    const random = Math.random().toString(36);
    const username = `testuser-${random}`;
    const registerResponse = await request.post(`/api/v1/user/register`, {
      data: {
        username,
        password: "password",
      },
    });
    expect(registerResponse.ok()).toBeTruthy();

    await use({ username, password: "password" });
  },
});
