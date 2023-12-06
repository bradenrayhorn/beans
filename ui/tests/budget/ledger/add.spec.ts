import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  test,
} from "../../setup";

test("can add transaction", async ({ budget: { id }, page }) => {
  await createAccount(id, "Checking", page.context().request);
  const groupID = await createCategoryGroup(
    id,
    "Bills",
    page.context().request,
  );
  await createCategory(id, groupID, "Electric", page.context().request);

  // go to ledger page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // open modal and add transaction
  await page.getByRole("link", { name: "Add" }).click();

  await page.getByLabel("Date").locator("visible=true").fill("2022-10-14");
  await page
    .getByRole("combobox", { name: "Account" })
    .selectOption("Checking");
  await page
    .getByRole("combobox", { name: "Category" })
    .selectOption("Electric");
  await page.getByLabel("Amount").locator("visible=true").fill("10.78");
  await page.getByLabel("Notes").locator("visible=true").fill("Test notes");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // transaction should be added (1 header row and 1 data row)
  await expect(page.getByRole("row")).toHaveCount(2);
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(1)).toHaveText("2022-10-14");
  await expect(cells.nth(2)).toHaveText("Electric Checking");
  await expect(cells.nth(3)).toHaveText("Electric");
  await expect(cells.nth(4)).toHaveText("Checking");
  await expect(cells.nth(5)).toHaveText("Test notes");
  await expect(cells.nth(6)).toHaveText("$10.78");
});
