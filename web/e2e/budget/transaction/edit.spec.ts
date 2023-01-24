import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  test,
} from "../../test.js";

test("can edit transaction", async ({ budget: { id }, page }) => {
  await createAccount(id, "Checking", page.context().request);
  await createAccount(id, "Savings", page.context().request);
  const groupID = await createCategoryGroup(
    id,
    "Bills",
    page.context().request
  );
  await createCategory(id, groupID, "Electric", page.context().request);

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "transactions" }).click();

  // no transactions exist right away (1 header row)
  expect(await page.getByRole("row").count()).toBe(1);

  // open modal and add transaction
  await page.getByRole("button", { name: "Add" }).click();

  await page.getByLabel("Date").fill("2022-10-14");
  await page.getByLabel("Account").click();
  await page.getByRole("option").filter({ hasText: "Checking" }).click();
  await page.getByLabel("Amount").fill("10.78");
  await page.getByLabel("Notes").fill("Test notes");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // transaction should be added (1 header row and 1 data row)
  expect(await page.getByRole("row").count()).toBe(2);
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(1)).toHaveText("10/14/2022");
  await expect(cells.nth(2)).toHaveText("");
  await expect(cells.nth(3)).toHaveText("Checking");
  await expect(cells.nth(4)).toHaveText("Test notes");
  await expect(cells.nth(5)).toHaveText("$10.78");

  // select and edit transaction
  await page
    .getByRole("row")
    .nth(1)
    .getByRole("checkbox", { name: "Select transaction" })
    .focus();
  await page
    .getByRole("row")
    .nth(1)
    .getByRole("checkbox", { name: "Select transaction" })
    .press("Space");
  await page.locator("body").press("e");

  await page.getByLabel("Date").fill("2023-01-23");
  await page.getByLabel("Account").click();
  await page.getByRole("option").filter({ hasText: "Savings" }).click();
  await page.getByLabel("Category").click();
  await page.getByRole("option").filter({ hasText: "Electric" }).click();
  await page.getByLabel("Amount").fill("15.00");
  await page.getByLabel("Notes").fill("");
  await page.getByRole("button", { name: "Save" }).click();
  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  await expect(cells.nth(1)).toHaveText("01/23/2023");
  await expect(cells.nth(2)).toHaveText("Electric");
  await expect(cells.nth(3)).toHaveText("Savings");
  await expect(cells.nth(4)).toHaveText("");
  await expect(cells.nth(5)).toHaveText("$15.00");
});
