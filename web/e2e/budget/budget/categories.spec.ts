import { expect } from "@playwright/test";
import { createCategory, createCategoryGroup, test } from "../../test.js";

test("can edit categories", async ({ budget: { id }, page, request }) => {
  const groupID = await createCategoryGroup(id, "Bills", request);
  await createCategory(id, groupID, "Electric", request);

  // go to budget page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: /^budget$/ }).click();

  const billsCategoryGroup = page
    .getByRole("list", { name: "Categories" })
    .filter({ hasText: "Bills" });

  await expect(billsCategoryGroup).toBeVisible();

  const electricCategory = billsCategoryGroup
    .getByRole("list")
    .getByRole("listitem")
    .filter({ hasText: "Electric" });

  await expect(electricCategory).toBeVisible();

  const assigned = electricCategory
    .getByRole("group", {
      name: "Assigned",
    })
    .getByRole("definition");
  await expect(assigned).toHaveText("$0.00");

  const drawer = page.getByRole("dialog", { name: "Edit Electric" });

  await electricCategory.getByRole("button", { name: "Edit Electric" }).click();
  await expect(drawer).toBeVisible();
  await page.getByLabel("Amount").fill("60.31");
  await page.getByRole("button", { name: "Save" }).click();
  await expect(drawer).toBeHidden();

  await expect(assigned).toHaveText("$60.31");
});
