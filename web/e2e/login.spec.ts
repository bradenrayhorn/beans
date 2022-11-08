import { expect } from "@playwright/test";
import { test } from "./test.js";

test("can login", async ({ user: { username, password }, page }) => {
  await page.goto("/login");

  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill(password);

  await page.getByRole("button", { name: "Log in" }).click();

  await expect(page).toHaveURL(/.*\/budget$/);
});
