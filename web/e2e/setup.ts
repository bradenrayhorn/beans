import { expect, test as base } from "@playwright/test";

export type RegisterFixture = {
  username: string;
  password: string;
};

export type LoginFixture = {};

export type BudgetFixture = { name: string; id: string };

type Fixtures = {
  register: RegisterFixture;
  login: LoginFixture;
  budget: BudgetFixture;
};

const randomString = () => Math.random().toString(36);

export const test = base.extend<Fixtures>({
  register: async ({ request }, use) => {
    const username = `testuser-${randomString()}`;
    const registerResponse = await request.post(`/api/v1/user/register`, {
      data: {
        username,
        password: "password",
      },
    });
    expect(registerResponse.ok()).toBeTruthy();

    await use({ username, password: "password" });
  },
  login: async ({ request, register, page }, use) => {
    const loginResponse = await request.post(`/api/v1/user/login`, {
      data: {
        username: register.username,
        password: register.password,
      },
    });
    expect(loginResponse.ok()).toBeTruthy();
    const state = await request.storageState();
    page.context().addCookies(state.cookies);

    await use({});
  },
  budget: async ({ login: _, request }, use) => {
    const name = `budget-${randomString()}`;
    const response = await request.post(`/api/v1/budgets`, {
      data: {
        name,
      },
    });
    expect(response.ok()).toBeTruthy();
    const data = await response.json();

    await use({ name, id: data?.data?.id });
  },
});

export const createAccount = async (
  budgetID: string,
  name: string,
  request: any
): Promise<string> => {
  const response = await request.post(`/api/v1/accounts`, {
    data: {
      name,
    },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();
  const data = await response.json();
  return data.data?.id;
};

export const createCategoryGroup = async (
  budgetID: string,
  name: string,
  request: any
): Promise<string> => {
  const response = await request.post(`/api/v1/categories/groups`, {
    data: {
      name,
    },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();

  const data = await response.json();
  return data.data?.id;
};

export const createCategory = async (
  budgetID: string,
  groupID: string,
  name: string,
  request: any
): Promise<string> => {
  const response = await request.post(`/api/v1/categories`, {
    data: {
      group_id: groupID,
      name,
    },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();

  const data = await response.json();
  return data.data?.id;
};

export const createTransaction = async (
  budgetID: string,
  categoryID: string,
  accountID: string,
  amount: string,
  date: string,
  request: any
) => {
  const response = await request.post(`/api/v1/transactions`, {
    data: {
      date,
      category_id: categoryID,
      account_id: accountID,
      amount,
    },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();
};
