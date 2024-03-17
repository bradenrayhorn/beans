import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createTransaction,
  test,
} from "../../setup";

test("can delete transaction", async ({ budget: { id }, page, request }) => {
  const account = await createAccount(id, request, { name: "Checking" });
  const groupID = await createCategoryGroup(id, "Bills", request);
  const category = await createCategory(id, groupID, "Electric", request);
  await createCategory(id, groupID, "Home", request);

  await createTransaction(id, request, {
    date: "2022-10-11",
    accountID: account,
    categoryID: category,
    amount: "10.72",
  });

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  await expect(page.getByRole("row")).toHaveCount(2);

  // select and open edit form
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  // delete transaction
  await page.getByRole("button", { name: "Delete Transaction" }).click();

  // there should be none left (count is 1 because of header row)
  await expect(page.getByRole("row")).toHaveCount(1);
});
