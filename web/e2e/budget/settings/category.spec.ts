import { expect } from "@playwright/test";
import { test } from "../../test.js";

test("can add category", async ({ budget: { id }, page }) => {
  // go to settings page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "settings" }).click();

  // add category group
  const form = page.getByRole("form", { name: "Add category group" });
  await form.getByLabel("Name").fill("Bills");

  await page.getByRole("button", { name: "Add Group" }).click();
  await expect(form.getByLabel("Name")).toBeEmpty();

  // add category
  const billsCategory = page
    .getByRole("list", { name: "Categories" })
    .getByRole("listitem")
    .filter({ hasText: "Bills" });

  await expect(billsCategory).toBeVisible();
  await billsCategory.getByRole("button", { name: "Add category" }).click();

  const addModal = page.getByRole("dialog", { name: "Add Category" });
  await expect(addModal).toBeVisible();

  await addModal.getByLabel("Name").fill("Electric");
  await addModal.getByRole("button", { name: "Add" }).click();
  await expect(addModal).toBeHidden();

  // category was added to group
  await expect(
    billsCategory
      .getByRole("list")
      .getByRole("listitem")
      .filter({ hasText: "Electric" })
  ).toBeVisible();
});
