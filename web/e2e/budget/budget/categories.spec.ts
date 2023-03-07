import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createTransaction,
  test,
} from "../../test.js";

test("can edit categories", async ({ budget: { id }, page, request }) => {
  const groupID = await createCategoryGroup(id, "Bills", request);
  const categoryID = await createCategory(id, groupID, "Electric", request);
  const accountID = await createAccount(id, "Checking", request);
  const currentDate = new Date().toISOString().substring(0, 10);
  await createTransaction(
    id,
    categoryID,
    accountID,
    "-20",
    currentDate,
    request
  );

  // go to budget page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: /^budget$/ }).click();

  const toBudget = page.getByLabel("To Budget:");
  await expect(toBudget).toHaveText("$0.00");

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
  const activity = electricCategory
    .getByRole("group", {
      name: "Activity",
    })
    .getByRole("definition");
  const available = electricCategory
    .getByRole("group", {
      name: "Available",
    })
    .getByRole("definition");
  await expect(assigned).toHaveText("$0.00");
  await expect(activity).toHaveText("-$20.00");
  await expect(available).toHaveText("-$20.00");

  const drawer = page.getByRole("dialog", { name: "Edit Electric" });

  await electricCategory.getByRole("button", { name: "Edit Electric" }).click();
  await expect(drawer).toBeVisible();
  await page.getByLabel("Amount").fill("60.31");
  await page.getByRole("button", { name: "Save" }).click();
  await expect(drawer).toBeHidden();

  await expect(assigned).toHaveText("$60.31");
  await expect(activity).toHaveText("-$20.00");
  await expect(available).toHaveText("$40.31");

  await expect(toBudget).toHaveText("-$60.31");

  // navigate to next month
  await page.getByRole("button", { name: "Next month" }).click();

  // available should have carried over
  await expect(available).toHaveText("$40.31");
});
