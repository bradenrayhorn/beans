import { expect } from "@playwright/test";
import { test } from "../../setup.js";

test("can carryover funds", async ({ budget: { id }, page }) => {
  // go to budgets page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: /^budget$/ }).click();

  // expand To Budget panel
  page.getByRole("button", { name: "To Budget" }).click();

  const toBudget = page
    .getByRole("button", { name: "To Budget" })
    .getByRole("definition");
  const fromLastMonth = page.getByLabel("From last month:");
  const forNextMonth = page.getByLabel("For next month:");
  const forNextMonthTag = page
    .getByRole("button", { name: "For Next Month" })
    .getByLabel("For Next Month");

  // assert initial values
  await expect(toBudget).toHaveText("$0.00");
  await expect(fromLastMonth).toHaveText("$0.00");
  await expect(forNextMonth).toHaveText("-$0.00");
  await expect(forNextMonthTag).toHaveText("$0.00");

  // carryover $50
  await page.getByRole("button", { name: "For Next Month" }).click();

  const dialog = page.getByRole("dialog");
  await dialog.getByLabel("Carryover").type("50");
  await dialog.getByRole("button", { name: "Save" }).click();
  await expect(dialog).toBeHidden();

  // asssert carryover saved
  await expect(toBudget).toHaveText("-$50.00");
  await expect(fromLastMonth).toHaveText("$0.00");
  await expect(forNextMonth).toHaveText("-$50.00");
  await expect(forNextMonthTag).toHaveText("$50.00");

  // navigate to next month
  await page.getByRole("button", { name: /^Next month$/i }).click();

  // expand To Budget panel
  await page.getByRole("button", { name: "To Budget" }).click();

  // assert values updated for new month
  await expect(toBudget).toHaveText("$50.00");
  await expect(fromLastMonth).toHaveText("$50.00");
  await expect(forNextMonth).toHaveText("-$0.00");
  await expect(forNextMonthTag).toHaveText("$0.00");
});
