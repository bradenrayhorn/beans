import { expect } from "@playwright/test";
import { test } from "../../setup.js";

test("can carryover funds", async ({ budget: { id }, page }) => {
  // go to budgets page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "Budget", exact: true }).click();

  // expand To Budget panel
  page.getByRole("button", { name: "To Budget" }).click();

  const toBudget = page.getByRole("dialog", { name: "To budget breakdown" });

  // assert initial values
  await expect(
    page.getByRole("button", { name: "To Budget: $0.00" }),
  ).toBeVisible();
  await expect(toBudget).toContainText("From last month $0.00");
  await expect(toBudget).toContainText("For next month -$0.00");

  // edit carryover - close popup and open carryover form
  await page.getByRole("link", { name: "Budget", exact: true }).click();
  await page.getByRole("link", { name: "Save for next month" }).click();

  // carryover $50
  await page.getByRole("textbox", { name: "Carryover" }).clear();
  await page.getByRole("textbox", { name: "Carryover" }).type("50");
  await page.getByRole("button", { name: "Save" }).click();

  // close carryover form
  await page.getByRole("link", { name: "Budget", exact: true }).click();

  // assert carryover saved
  await expect(
    page.getByRole("button", { name: "To Budget: -$50.00" }),
  ).toBeVisible();
  await page.getByRole("button", { name: "To Budget" }).click();
  await expect(toBudget).toBeVisible();
  await expect(toBudget).toContainText("From last month $0.00");
  await expect(toBudget).toContainText("For next month -$50.00");

  // navigate to next month
  await page.getByRole("link", { name: "Next month", exact: true }).click();

  // assert values updated for new month
  await expect(
    page.getByRole("button", { name: "To Budget: $50.00" }),
  ).toBeVisible();
  await page.getByRole("button", { name: "To Budget" }).click();
  await expect(toBudget).toContainText("From last month $50.00");
  await expect(toBudget).toContainText("For next month -$0.00");
});
