import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createTransaction,
  test,
} from "../../setup";

test("can view account", async ({ budget: { id }, page, request }) => {
  // create account and transaction
  const groupID = await createCategoryGroup(id, "Bills", request);
  const categoryID = await createCategory(id, groupID, "Electric", request);
  const accountID = await createAccount(id, request, { name: "Checking" });
  createAccount(id, request, { name: "401k", offBudget: true });
  await createTransaction(id, request, {
    date: new Date().toISOString().substring(0, 10),
    accountID,
    categoryID,
    amount: "20.43",
  });

  // go to accounts page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "accounts" }).click();

  const account = page
    .getByRole("listitem")
    .filter({ has: page.getByRole("heading", { name: "Checking" }) });
  await expect(account).toBeVisible();

  await expect(account).toContainText("Balance: $20.43");

  const offBudgetAccount = page
    .getByRole("listitem")
    .filter({ has: page.getByRole("heading", { name: "401k" }) });
  await expect(offBudgetAccount).toBeVisible();

  await expect(offBudgetAccount).toContainText("Off-Budget");
});
