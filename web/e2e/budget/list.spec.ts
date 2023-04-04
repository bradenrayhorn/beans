import { expect } from "@playwright/test";
import { test } from "../setup.js";

test("can add and go to budget", async ({ login: _, page }) => {
  await page.goto("/budget");

  // no budgets exist right now
  await expect(page.getByText("No existing budgets found.")).toBeVisible();

  // open modal, fill out form
  await page.getByRole("button", { name: "New Budget" }).click();
  await expect(page.getByRole("dialog", { name: "New Budget" })).toBeVisible();

  await page.getByLabel("Name").fill("Test budget");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("dialog", { name: "New Budget" })).toBeHidden();

  // budget is now in list
  await expect(page.getByText("No existing budgets found.")).toBeHidden();

  // go to budget
  await page.getByRole("link", { name: "Test budget" }).click();
  await expect(page.getByRole("heading", { name: "beans" })).toBeVisible();
});
