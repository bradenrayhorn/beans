import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createTransaction,
  test,
} from "../../test.js";

test("can view account", async ({ budget: { id }, page, request }) => {
  // create account and transaction
  const groupID = await createCategoryGroup(id, "Bills", request);
  const categoryID = await createCategory(id, groupID, "Electric", request);
  const accountID = await createAccount(id, "Checking", request);
  await createTransaction(
    id,
    categoryID,
    accountID,
    "20.43",
    new Date().toISOString().substring(0, 10),
    request
  );

  // go to accounts page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "accounts" }).click();

  const account = page
    .getByRole("listitem")
    .filter({ has: page.getByRole("heading", { name: "Checking" }) });
  await expect(account).toBeVisible();

  const assigned = account
    .getByRole("group", {
      name: "Balance",
    })
    .getByRole("definition");

  await expect(assigned).toHaveText("$20.43");
});
