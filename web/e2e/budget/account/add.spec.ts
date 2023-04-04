import { expect } from "@playwright/test";
import { test } from "../../setup.js";

test("can add and view account", async ({ budget: { id }, page }) => {
  // go to accounts page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "accounts" }).click();

  // no accounts exist right away
  await expect(page.getByText("No accounts found.")).toBeVisible();

  // open modal and add account
  await page.getByRole("button", { name: "Add" }).click();
  await expect(page.getByRole("dialog", { name: "Add Account" })).toBeVisible();

  await page.getByLabel("Name").fill("Checking account");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("dialog", { name: "Add Account" })).toBeHidden();

  // account should be added on page
  await expect(page.getByText("No accounts found.")).toBeHidden();

  const account = page
    .getByRole("listitem")
    .filter({ has: page.getByRole("heading", { name: "Checking account" }) });
  await expect(account).toBeVisible();
});
