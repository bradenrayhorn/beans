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

test("can edit transaction", async ({ budget: { id }, page }) => {
  const account = await createAccount(id, "Checking", page.context().request);
  await createAccount(id, "Savings", page.context().request);
  const groupID = await createCategoryGroup(
    id,
    "Bills",
    page.context().request,
  );
  const category = await createCategory(
    id,
    groupID,
    "Electric",
    page.context().request,
  );
  await createCategory(id, groupID, "Home", page.context().request);
  await createPayee(id, "Workplace", page.context().request);

  await createTransaction(
    id,
    null,
    category,
    account,
    "10.72",
    "2022-10-11",
    page.context().request,
  );

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  await page.getByRole("textbox", { name: "Date" }).fill("2023-01-23");
  await page
    .getByRole("combobox", { name: "Account" })
    .selectOption("Checking");
  await page.getByRole("combobox", { name: "Payee" }).selectOption("Workplace");
  await selectOption(page, "Category", "Home");
  await page.getByRole("textbox", { name: "Amount" }).fill("15");
  await page.getByRole("textbox", { name: "Notes" }).fill("hi there");
  await page.getByRole("button", { name: "Save" }).click();

  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(1)).toHaveText("2023-01-23");
  await expect(cells.nth(2)).toHaveText("Workplace");
  await expect(cells.nth(3)).toHaveText("Home");
  await expect(cells.nth(4)).toHaveText("Checking");
  await expect(cells.nth(5)).toHaveText("hi there");
  await expect(cells.nth(6)).toHaveText("$15.00");
});

test("can edit payee to blank", async ({ budget: { id }, page }) => {
  const account = await createAccount(id, "Checking", page.context().request);
  await createAccount(id, "Savings", page.context().request);
  const groupID = await createCategoryGroup(
    id,
    "Bills",
    page.context().request,
  );
  const category = await createCategory(
    id,
    groupID,
    "Electric",
    page.context().request,
  );
  await createCategory(id, groupID, "Home", page.context().request);
  const payeeID = await createPayee(id, "Workplace", page.context().request);

  await createTransaction(
    id,
    payeeID,
    category,
    account,
    "10.72",
    "2022-10-11",
    page.context().request,
  );

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // select and edit transaction
  await page.getByRole("row").nth(1).getByRole("checkbox").focus();
  await page.getByRole("row").nth(1).getByRole("checkbox").press("Space");
  await page.getByRole("link", { name: "edit" }).click();

  await page.getByRole("combobox", { name: "Payee" }).selectOption("");
  await page.getByRole("button", { name: "Save" }).click();

  // check that payee is now blank
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(2)).toHaveText("");
});
