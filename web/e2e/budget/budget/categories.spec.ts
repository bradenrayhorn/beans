import { expect, Locator } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createTransaction,
  test,
} from "../../setup.js";

const getAssigned = (locator: Locator) => locator.getByRole("cell").nth(1);
const getSpent = (locator: Locator) => locator.getByRole("cell").nth(2);
const getAvailable = (locator: Locator) => locator.getByRole("cell").nth(3);
const getReceived = (locator: Locator) => locator.getByRole("cell").nth(1);

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

  // STEP 1: Add income
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: /^transactions$/ }).click();

  await page.getByRole("button", { name: "Add" }).click();
  await page.getByLabel("Date").fill(currentDate);
  await page.getByLabel("Account").click();
  await page.getByRole("option").filter({ hasText: "Checking" }).click();
  await page.getByLabel("Category").click();
  await page.getByRole("option").filter({ hasText: "Income" }).click();
  await page.getByLabel("Amount").fill("100");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // STEP 2: View & edit budget
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: /^budget$/ }).click();

  // expand to budget panel
  page.getByRole("button", { name: "To Budget" }).click();

  const toBudget = page
    .getByRole("button", { name: "To Budget" })
    .getByRole("definition");
  await expect(toBudget).toHaveText("$100.00");

  const income = page.getByLabel("Income:");
  await expect(income).toHaveText("$100.00");

  const assignedThisMonth = page.getByLabel("Assigned this month:");
  await expect(assignedThisMonth).toHaveText("-$0.00");

  const incomeCategory = page
    .getByRole("table", { name: "Income" })
    .getByRole("row")
    .filter({ has: page.getByRole("cell").filter({ hasText: "Income" }) });
  await expect(incomeCategory).toBeVisible();

  await expect(getReceived(incomeCategory)).toHaveText("$100.00");

  const billsCategoryGroup = page.getByRole("rowgroup", { name: "Bills" });
  await expect(billsCategoryGroup).toBeVisible();

  const electricCategory = billsCategoryGroup
    .getByRole("row")
    .filter({ hasText: "Electric" });
  await expect(electricCategory).toBeVisible();

  const assigned = getAssigned(electricCategory);
  const spent = getSpent(electricCategory);
  const available = getAvailable(electricCategory);

  await expect(assigned).toHaveText("$0.00");
  await expect(spent).toHaveText("-$20.00");
  await expect(available).toHaveText("-$20.00");

  const editPopup = page.getByRole("dialog", { name: "Edit assigned" });

  await electricCategory.getByRole("button").click();
  await expect(editPopup).toBeVisible();
  await editPopup.getByLabel("Assigned").fill("60.31");
  await editPopup.getByRole("button", { name: "Save" }).click();
  await expect(editPopup).toBeHidden();

  await expect(assigned).toHaveText("$60.31");
  await expect(spent).toHaveText("-$20.00");
  await expect(available).toHaveText("$40.31");

  await expect(toBudget).toHaveText("$39.69");
  await expect(income).toHaveText("$100.00");
  await expect(assignedThisMonth).toHaveText("-$60.31");

  // navigate to next month
  await page.getByRole("button", { name: /^Next month$/i }).click();

  // available should have carried over
  await expect(available).toHaveText("$40.31");
});
