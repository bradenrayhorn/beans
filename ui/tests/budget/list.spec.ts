import { expect } from "@playwright/test";
import { test } from "../setup.js";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
test("can add and go to budget", async ({ login: _, page }) => {
  await page.goto("/budget");

  // no budgets exist right now
  await expect(page.getByText("No existing budgets found.")).toBeVisible();

  // go to new budget page, fill out form
  await page.getByRole("link", { name: "New Budget" }).click();

  await page.getByLabel("Name").fill("Test budget");
  await page.getByRole("button", { name: "Save" }).click();

  // budget is now in list
  await expect(page.getByText("No existing budgets found.")).toBeHidden();

  // go to budget
  await page.getByRole("link", { name: "Test budget" }).click();
  await expect(page.getByText("beans")).toBeVisible();
});
