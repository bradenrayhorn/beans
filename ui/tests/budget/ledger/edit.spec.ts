import { expect } from "@playwright/test";

import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createPayee,
  createTransaction,
  selectOption,
  test,
} from "../../setup";

test("can edit transaction", async ({ budget: { id }, page, request }) => {
  const account = await createAccount(id, request, { name: "Checking" });
  await createAccount(id, request, { name: "Savings" });
  const groupID = await createCategoryGroup(id, "Bills", request);
  const category = await createCategory(id, groupID, "Electric", request);
  await createCategory(id, groupID, "Home", request);
  await createPayee(id, "Workplace", request);

  await createTransaction(id, request, {
    date: "2022-10-11",
    accountID: account,
    categoryID: category,
    amount: "10.72",
  });

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  await page.getByRole("textbox", { name: "Date" }).fill("2023-01-23");
  await selectOption(page, "Account", "Savings");
  await selectOption(page, "Payee", "Workplace");
  await selectOption(page, "Category", "Home");
  await page.getByRole("textbox", { name: "Amount" }).fill("15");
  await page.getByRole("textbox", { name: "Notes" }).fill("hi there");
  await page.getByRole("button", { name: "Save" }).click();

  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(1)).toHaveText("2023-01-23");
  await expect(cells.nth(2)).toHaveText("Workplace");
  await expect(cells.nth(3)).toHaveText("Home");
  await expect(cells.nth(4)).toHaveText("Savings");
  await expect(cells.nth(5)).toHaveText("hi there");
  await expect(cells.nth(6)).toHaveText("$15.00");
});

test("can edit payee to blank", async ({ budget: { id }, page, request }) => {
  const account = await createAccount(id, request, { name: "Checking" });
  const groupID = await createCategoryGroup(id, "Bills", request);
  const category = await createCategory(id, groupID, "Electric", request);
  await createCategory(id, groupID, "Home", request);
  const payeeID = await createPayee(id, "Workplace", request);

  await createTransaction(id, request, {
    date: "2022-10-11",
    accountID: account,
    categoryID: category,
    payeeID,
    amount: "10.72",
  });

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  await selectOption(page, "Payee", "None");
  await page.getByRole("button", { name: "Save" }).click();

  // check that payee is now blank
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(2)).toHaveText("");
});

test("can edit transaction to off-budget account", async ({
  budget: { id },
  page,
  request,
}) => {
  const checkingAccount = await createAccount(id, request, {
    name: "Checking",
  });
  await createAccount(id, request, { name: "401k", offBudget: true });

  await createTransaction(id, request, {
    date: "2022-10-11",
    accountID: checkingAccount,
    amount: "10.72",
  });

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  // category is editable
  await expect(page.getByRole("combobox", { name: "Category" })).toBeEnabled();

  // select an on-budget account, category is editable
  await selectOption(page, "Account", "Checking");
  await expect(page.getByRole("combobox", { name: "Category" })).toBeEnabled();

  // select an off-budget account, category is not editable
  await selectOption(page, "Account", "401k");
  await expect(page.getByRole("textbox", { name: "Category" })).toBeDisabled();

  // save changes and verify
  await page.getByRole("button", { name: "Save" }).click();

  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(3)).toHaveText("Off-Budget");
});

test("can edit transfer transaction", async ({
  budget: { id },
  page,
  request,
}) => {
  const checkingAccount = await createAccount(id, request, {
    name: "Checking",
  });
  const savingsAccount = await createAccount(id, request, {
    name: "Savings",
  });
  await createAccount(id, request, {
    name: "T-Bills",
  });

  await createTransaction(id, request, {
    date: "2022-10-11",
    accountID: checkingAccount,
    transferAccountID: savingsAccount,
    amount: "10.72",
  });

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  const transactionOnCheckingAccount = page
    .getByRole("row")
    .filter({ has: page.getByRole("cell").nth(4).getByText("Checking") });
  await transactionOnCheckingAccount.getByRole("checkbox").focus();
  await transactionOnCheckingAccount.getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  // change account & amount
  await selectOption(page, "Account", "T-Bills");
  await page.getByRole("textbox", { name: "Amount" }).fill("15");

  // save changes
  await page.getByRole("button", { name: "Save" }).click();

  const tBills = page
    .getByRole("row")
    .filter({ has: page.getByRole("cell").nth(4).getByText("T-Bills") })
    .getByRole("cell");
  await expect(tBills.nth(2)).toHaveText("Savings");
  await expect(tBills.nth(3)).toHaveText("Transfer");
  await expect(tBills.nth(4)).toHaveText("T-Bills");
  await expect(tBills.nth(6)).toHaveText("$15.00");

  const savings = page
    .getByRole("row")
    .filter({ has: page.getByRole("cell").nth(4).getByText("Savings") })
    .getByRole("cell");
  await expect(savings.nth(2)).toHaveText("T-Bills");
  await expect(savings.nth(3)).toHaveText("Transfer");
  await expect(savings.nth(4)).toHaveText("Savings");
  await expect(savings.nth(6)).toHaveText("-$15.00");
});

test("can edit split", async ({ budget: { id }, page, request }) => {
  const checkingAccount = await createAccount(id, request, {
    name: "Checking",
  });
  const groupID = await createCategoryGroup(id, "Bills", request);
  const electricID = await createCategory(id, groupID, "Electric", request);
  await createCategory(id, groupID, "Water", request);

  await createTransaction(id, request, {
    date: "2022-10-11",
    accountID: checkingAccount,
    amount: "10.72",
    splits: [{ amount: "10.72", categoryID: electricID }],
  });

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  // change split info
  const parent = page.getByRole("group", { name: "Parent Transaction" });
  await parent.getByLabel("Amount").fill("10.78");

  const split = page.getByRole("group", { name: "Split 1" });
  await split.getByLabel("Amount").fill("10.78");
  await split.getByLabel("Notes").fill(":)");
  await selectOption(page, "Category", "Water");

  // save
  await page.getByRole("button", { name: "Save" }).click();
  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(3)).toHaveText("Split");
  await expect(cells.nth(6)).toHaveText("$10.78");

  // open form and check splits
  await page.getByRole("link", { name: "edit" }).click();

  await expect(split.getByLabel("Amount")).toHaveValue("10.78");
  await expect(split.getByLabel("Notes")).toHaveValue(":)");
  await expect(split.getByRole("combobox", { name: "Category" })).toHaveValue(
    "Water",
  );
});
