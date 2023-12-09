import { expect } from "@playwright/test";
import { test } from "../../../setup";

test("can add payee", async ({ budget: { id }, page }) => {
  // go to settings page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "settings" }).click();

  await page.getByRole("link", { name: "Payees" }).click();

  await expect(page.getByText("No payees found.")).toBeVisible();

  // add new payee
  await page.getByRole("link", { name: "New Payee" }).click();

  await page.getByLabel("Name").fill("Workplace");

  await page.getByRole("button", { name: "Save" }).click();

  // check if payee is added
  await expect(
    page.getByRole("listitem").filter({ hasText: "Workplace" }),
  ).toBeVisible();
  await expect(page.getByText("No payees found.")).toBeHidden();
});
