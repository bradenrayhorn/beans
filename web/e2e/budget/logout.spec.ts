import { expect } from "@playwright/test";
import { test } from "../test.js";

test("can logout", async ({ register: { username }, budget: { id }, page }) => {
  await page.goto(`/budget/${id}`);
  await page.getByRole("button", { name: username }).click();

  await page.getByRole("menuitem", { name: "Log out" }).click();

  await expect(page).toHaveURL(/.*\/login$/);
});
