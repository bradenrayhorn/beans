import {
  expect,
  test as base,
  type APIRequestContext,
  type Page,
} from "@playwright/test";

export type RegisterFixture = {
  username: string;
  password: string;
};

export type LoginFixture = Record<string, never>;

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
  request: APIRequestContext,
  { name, offBudget = false }: { name: string; offBudget?: boolean },
): Promise<string> => {
  const response = await request.post(`/api/v1/accounts`, {
    data: { name, offBudget },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();
  const data = await response.json();
  return data.data?.id;
};

export const createCategoryGroup = async (
  budgetID: string,
  name: string,
  request: APIRequestContext,
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
  request: APIRequestContext,
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
  request: APIRequestContext,
  {
    payeeID,
    categoryID,
    amount,
    date,
    accountID,
    transferAccountID,
    splits,
  }: {
    payeeID?: string;
    categoryID?: string;
    accountID: string;
    amount: string;
    date: string;
    transferAccountID?: string;
    splits?: Array<{ amount: string; categoryID: string }>;
  },
) => {
  const response = await request.post(`/api/v1/transactions`, {
    data: {
      date,
      payee_id: payeeID,
      category_id: categoryID,
      account_id: accountID,
      transferAccountID,
      amount,
      splits: splits?.map((split) => ({
        amount: split.amount,
        category_id: split.categoryID,
      })),
    },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();
};

export const createPayee = async (
  budgetID: string,
  name: string,
  request: APIRequestContext,
): Promise<string> => {
  const response = await request.post(`/api/v1/payees`, {
    data: {
      name,
    },
    headers: { "Budget-ID": budgetID },
  });

  expect(response.ok()).toBeTruthy();

  const data = await response.json();
  return data.data?.id;
};

export const selectOption = async (
  page: Page,
  label: string,
  option: string,
) => {
  await page.getByRole("combobox", { name: label }).click();
  await page.getByRole("option", { name: option }).click();
};
