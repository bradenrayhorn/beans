import { expect } from "@playwright/test";
import { test } from "../../../setup";

test("can add category", async ({ budget: { id }, page }) => {
  // go to settings page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "settings" }).click();

  await page.getByRole("link", { name: "Categories" }).click();

  // add category group
  await page.getByRole("link", { name: "Add new group" }).click();

  await page.getByLabel("Name").fill("Bills");

  await page.getByRole("button", { name: "Save" }).click();

  // add category
  await page.getByRole("link", { name: "Bills" }).click();
  await page.getByRole("link", { name: "Add new category" }).click();
  await page.getByLabel("Name").fill("Electric");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByText("Electric")).toBeVisible();
});
