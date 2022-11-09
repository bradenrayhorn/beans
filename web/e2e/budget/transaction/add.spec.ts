import { expect } from "@playwright/test";
import { createAccount, test } from "../../test.js";

test("can add transaction", async ({ budget: { id }, page }) => {
  await createAccount(id, "Checking", page.context().request);

  // go to transactions page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "transactions" }).click();

  // no transactions exist right away (1 header row)
  expect(await page.getByRole("row").count()).toBe(1);

  // open modal and add transaction
  await page.getByRole("button", { name: "Add" }).click();
  await expect(
    page.getByRole("dialog", { name: "Add Transaction" })
  ).toBeVisible();

  await page.getByLabel("Date").fill("2022-10-14");
  await page.getByLabel("Account").click();
  await page.getByRole("option").filter({ hasText: "Checking" }).click();
  await page.getByLabel("Amount").fill("10.78");
  await page.getByLabel("Notes").fill("Test notes");
  await page.getByRole("button", { name: "Add" }).click();

  await expect(
    page.getByRole("dialog", { name: "Add Transaction" })
  ).toBeHidden();

  // transaction should be added (1 header row and 1 data row)
  expect(await page.getByRole("row").count()).toBe(2);
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(0)).toHaveText("10/14/2022");
  await expect(cells.nth(1)).toHaveText("Checking");
  await expect(cells.nth(2)).toHaveText("Test notes");
  await expect(cells.nth(3)).toHaveText("$10.78");
});
