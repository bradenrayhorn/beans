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
  await page.getByRole("link", { name: /^budget$/ }).click();

  const billsCategoryGroup = page.getByRole("rowgroup", { name: "Bills" });
  const category = billsCategoryGroup
    .getByRole("row")
    .filter({ hasText: "Electric" });
  const assigned = category.getByRole("cell").nth(1);

  const formatMonth = (date: Date) =>
    `${date.getFullYear()}.${`${date.getMonth() + 1}`.padStart(2, "0")}`;

  // month header is correct
  await expect(
    page.getByRole("heading", { name: formatMonth(new Date()) })
  ).toBeVisible();

  // fill out category
  await assigned.getByRole("button").click();
  const editPopup = page.getByRole("dialog");
  await editPopup.getByLabel("Assigned").fill("54");
  await editPopup.getByRole("button", { name: "Save" }).click();
  await expect(editPopup).toBeHidden();
  expect(assigned).toHaveText("$54.00");

  // navigate to next month
  await page.getByRole("button", { name: /^Next month$/ }).click();

  // month header should change
  const nextMonth = new Date();
  nextMonth.setMonth(nextMonth.getMonth() + 1, 1);
  await expect(
    page.getByRole("heading", { name: formatMonth(nextMonth) })
  ).toBeVisible();

  // new month, assigned value should have changed
  expect(assigned).toHaveText("$0.00");

  // navigate to previous month
  await page.getByRole("button", { name: /^Previous month$/ }).click();

  // month header should change
  await expect(
    page.getByRole("heading", { name: formatMonth(new Date()) })
  ).toBeVisible();

  // new month, assigned value should have changed
  await expect(assigned).toBeVisible();
  expect(assigned).toHaveText("$54.00");
  await expect(assigned).toBeVisible();
});
