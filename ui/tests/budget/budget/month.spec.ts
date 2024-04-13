import { expect } from "@playwright/test";
import { createCategory, createCategoryGroup, test } from "../../setup.js";

test("can navigate between months", async ({
  budget: { id },
  page,
  request,
}) => {
  const groupID = await createCategoryGroup(id, "Bills", request);
  await createCategory(id, groupID, "Electric", request);

  // go to budgets page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "Budget" }).click();

  const billsCategoryGroup = page.getByRole("rowgroup", { name: "Bills" });
  const category = billsCategoryGroup
    .getByRole("row")
    .filter({ hasText: "Electric" });
  const assigned = category.getByRole("cell").nth(1);

  const formatMonth = (date: Date) =>
    `${date.toLocaleString("default", { month: "long" })} ${date.getFullYear()}`;

  // month header is correct
  await expect(page.getByText(formatMonth(new Date()))).toBeVisible();

  // fill out category
  await category.click();
  await page.getByRole("textbox", { name: "Assigned" }).fill("54");
  await page.getByRole("button", { name: "Save" }).click();
  await expect(assigned).toHaveText("$54.00");

  // navigate to next month
  await page.getByRole("link", { name: "Next month", exact: true }).click();

  // month header should change
  const nextMonth = new Date();
  nextMonth.setMonth(nextMonth.getMonth() + 1, 1);
  await expect(page.getByText(formatMonth(nextMonth))).toBeVisible();

  // new month, assigned value should have changed
  await expect(assigned).toHaveText("$0.00");

  // navigate to previous month
  await page.getByRole("link", { name: "Previous month" }).click();

  // month header should change
  await expect(page.getByText(formatMonth(new Date()))).toBeVisible();

  // back to first month, assigned value should have changed back
  await expect(assigned).toBeVisible();
  await expect(assigned).toHaveText("$54.00");
  await expect(assigned).toBeVisible();
});
