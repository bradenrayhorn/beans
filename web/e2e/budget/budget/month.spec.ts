import { expect } from "@playwright/test";
import { createCategory, createCategoryGroup, test } from "../../test.js";

test("can navigate between months", async ({
  budget: { id },
  page,
  request,
}) => {
  const groupID = await createCategoryGroup(id, "Bills", request);
  await createCategory(id, groupID, "Electric", request);

  // go to budgets page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "budget" }).click();

  const categoryGroup = page
    .getByRole("list", { name: "Categories" })
    .filter({ hasText: "Bills" });
  const category = categoryGroup
    .getByRole("list")
    .getByRole("listitem")
    .filter({ hasText: "Electric" });
  const assigned = category
    .getByRole("group", {
      name: "Assigned",
    })
    .getByRole("definition");

  const formatMonth = (date: Date) =>
    `${date.getFullYear()}.${date.toISOString().substring(5, 7)}`;

  // month header is correct
  await expect(
    page.getByRole("heading", { name: formatMonth(new Date()) })
  ).toBeVisible();

  // fill out category
  await category.getByRole("button", { name: `Edit Electric` }).click();
  await page.getByLabel("Amount").fill("54");
  await page.getByRole("button", { name: "Save" }).click();
  expect(assigned).toHaveText("$54.00");

  // navigate to next month
  await page.getByRole("button", { name: "Next month" }).click();

  // month header should change
  const nextMonth = new Date();
  nextMonth.setMonth(nextMonth.getMonth() + 1);
  await expect(
    page.getByRole("heading", { name: formatMonth(nextMonth) })
  ).toBeVisible();

  // new month, assigned value should have changed
  expect(assigned).toHaveText("$0.00");

  // navigate to previous month
  await page.getByRole("button", { name: "Previous month" }).click();

  // month header should change
  await expect(
    page.getByRole("heading", { name: formatMonth(new Date()) })
  ).toBeVisible();

  // new month, assigned value should have changed
  await expect(assigned).toBeVisible();
  expect(assigned).toHaveText("$54.00");
  await expect(assigned).toBeVisible();
});
