import { test } from "../../../setup";
import { expect } from "@playwright/test";

test("can logout", async ({ budget: { id }, page }) => {
  await page.goto(`/budget/${id}/settings`);

  await page.getByRole("link", { name: "General" }).click();

  await page.getByRole("link", { name: "Log out" }).click();

  await expect(page).toHaveURL(/.*\/login$/);
});
