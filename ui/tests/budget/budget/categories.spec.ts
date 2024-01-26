import { expect, type Locator } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createTransaction,
  selectOption,
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
    null,
    categoryID,
    accountID,
    "-20",
    currentDate,
    request,
  );

  const toBudgetButton = page.getByRole("button", { name: "To Budget" });
  const toBudget = page.getByRole("dialog", { name: "To budget breakdown" });

  const incomeCategory = page
    .getByRole("table", { name: "Income" })
    .getByRole("row")
    .filter({ has: page.getByRole("cell").filter({ hasText: "Income" }) });

  // STEP 1: Add income
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  await page.getByRole("link", { name: "Add" }).click();
  await page.getByRole("textbox", { name: "Date" }).fill(currentDate);
  await page
    .getByRole("combobox", { name: "Account" })
    .selectOption("Checking");
  await selectOption(page, "Category", "Income");
  await page.getByRole("textbox", { name: "Amount" }).fill("100");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // STEP 2: Check income has reflected in budget
  await page.getByRole("link", { name: "Budget", exact: true }).click();

  // expand "to budget" panel
  await page.getByRole("button", { name: "To Budget" }).click();
  await toBudgetButton.click();

  // assert initial values
  await expect(
    page.getByRole("button", { name: "To Budget: $100.00" }),
  ).toBeVisible();
  await expect(toBudget).toContainText("Income $100.00");
  await expect(toBudget).toContainText("Assigned this month -$0.00");

  await expect(incomeCategory).toBeVisible();
  await expect(getReceived(incomeCategory)).toHaveText("$100.00");

  // STEP 3: Check "Electric" transaction has reflected in budget

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

  // STEP 4: Assign $60.31 to "Electric"

  await electricCategory.click();
  await page.getByRole("textbox", { name: "Assigned" }).fill("60.31");
  await page.getByRole("button", { name: "Save" }).click();
  await page.getByRole("link", { name: "Close form" }).click();

  // STEP 5: Check assign has reflected in budget

  // check budget table
  await expect(assigned).toHaveText("$60.31");
  await expect(spent).toHaveText("-$20.00");
  await expect(available).toHaveText("$40.31");

  // check to budget breakdown
  await toBudgetButton.click();
  await expect(
    page.getByRole("button", { name: "To Budget: $39.69" }),
  ).toBeVisible();
  await expect(toBudget).toContainText("Income $100.00");
  await expect(toBudget).toContainText("Assigned this month -$60.31");

  // STEP 6: Navigate to next month and see if available dollars carried over

  await page.getByRole("link", { name: "Next month", exact: true }).click();
  await expect(available).toHaveText("$40.31");
});
