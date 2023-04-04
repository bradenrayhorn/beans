import { expect } from "@playwright/test";
import { test } from "./setup.js";

test("can login", async ({ register: { username, password }, page }) => {
  await page.goto("/login");

  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill(password);

  await page.getByRole("button", { name: "Log in" }).click();

  await expect(page).toHaveURL(/.*\/budget$/);
});

test("cannot login with invalid password", async ({
  register: { username },
  page,
}) => {
  await page.goto("/login");

  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill("a bad password");

  await page.getByRole("button", { name: "Log in" }).click();

  await expect(
    page.getByRole("alert").filter({ hasText: "Invalid username or password" })
  ).toBeVisible();

  await expect(page).toHaveURL(/.*\/login$/);
});
