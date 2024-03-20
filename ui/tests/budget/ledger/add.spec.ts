import { expect } from "@playwright/test";
import {
  createAccount,
  createCategory,
  createCategoryGroup,
  createPayee,
  selectOption,
  test,
} from "../../setup";

test("can add transaction", async ({ budget: { id }, page }) => {
  await createAccount(id, page.context().request, { name: "Checking" });
  const groupID = await createCategoryGroup(
    id,
    "Bills",
    page.context().request,
  );
  await createCategory(id, groupID, "Electric", page.context().request);
  await createPayee(id, "Workplace", page.context().request);

  // go to ledger page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // open modal and add transaction
  await page.getByRole("link", { name: "Add" }).click();

  await page.getByLabel("Date").locator("visible=true").fill("2022-10-14");
  await selectOption(page, "Payee", "Workplace");
  await selectOption(page, "Account", "Checking");
  await selectOption(page, "Category", "Electric");
  await page.getByLabel("Amount").locator("visible=true").fill("10.78");
  await page.getByLabel("Notes").locator("visible=true").fill("Test notes");
  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // transaction should be added (1 header row and 1 data row)
  await expect(page.getByRole("row")).toHaveCount(2);
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(1)).toHaveText("2022-10-14");
  await expect(cells.nth(2)).toHaveText("Workplace");
  await expect(cells.nth(3)).toHaveText("Electric");
  await expect(cells.nth(4)).toHaveText("Checking");
  await expect(cells.nth(5)).toHaveText("Test notes");
  await expect(cells.nth(6)).toHaveText("$10.78");
});

test("can add transaction with off-budget account", async ({
  budget: { id },
  page,
  request,
}) => {
  await createAccount(id, request, { name: "401k", offBudget: true });
  await createAccount(id, request, { name: "Checking", offBudget: false });

  // go to ledger page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // open modal and add transaction
  await page.getByRole("link", { name: "Add" }).click();

  await page.getByLabel("Date").locator("visible=true").fill("2022-10-14");
  await page.getByLabel("Amount").locator("visible=true").fill("10.78");

  // select an off-budget account
  await selectOption(page, "Account", "401k");

  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // transaction should be added (1 header row and 1 data row)
  await expect(page.getByRole("row")).toHaveCount(2);
  const cells = page.getByRole("row").nth(1).getByRole("cell");
  await expect(cells.nth(1)).toHaveText("2022-10-14");
  await expect(cells.nth(2)).toHaveText("");
  await expect(cells.nth(3)).toHaveText("Off-Budget");
  await expect(cells.nth(4)).toHaveText("401k");
  await expect(cells.nth(5)).toHaveText("");
  await expect(cells.nth(6)).toHaveText("$10.78");
});

test("can add transfer", async ({ budget: { id }, page, request }) => {
  await createAccount(id, request, { name: "Savings", offBudget: false });
  await createAccount(id, request, { name: "Checking", offBudget: false });

  // go to ledger page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // open modal and add transaction
  await page.getByRole("link", { name: "Add" }).click();

  await page.getByLabel("Date").locator("visible=true").fill("2022-10-14");
  await page.getByLabel("Amount").locator("visible=true").fill("10.78");

  // transfer from Savings -> Checking
  await selectOption(page, "Payee", "Savings");
  await selectOption(page, "Account", "Checking");

  await page.getByRole("button", { name: "Save" }).click();

  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // transactions should be added (1 header row and 2 data row)
  await expect(page.getByRole("row")).toHaveCount(3);

  const checking = page
    .getByRole("row")
    .filter({ has: page.getByRole("cell").nth(4).getByText("Checking") })
    .getByRole("cell");
  await expect(checking.nth(1)).toHaveText("2022-10-14");
  await expect(checking.nth(2)).toHaveText("Savings");
  await expect(checking.nth(3)).toHaveText("Transfer");
  await expect(checking.nth(4)).toHaveText("Checking");
  await expect(checking.nth(5)).toHaveText("");
  await expect(checking.nth(6)).toHaveText("$10.78");

  const savings = page
    .getByRole("row")
    .filter({ has: page.getByRole("cell").nth(4).getByText("Savings") })
    .getByRole("cell");
  await expect(savings.nth(1)).toHaveText("2022-10-14");
  await expect(savings.nth(2)).toHaveText("Checking");
  await expect(savings.nth(3)).toHaveText("Transfer");
  await expect(savings.nth(4)).toHaveText("Savings");
  await expect(savings.nth(5)).toHaveText("");
  await expect(savings.nth(6)).toHaveText("-$10.78");
});

test("can add split", async ({ budget: { id }, page, request }) => {
  await createAccount(id, request, { name: "Savings" });
  const groupID = await createCategoryGroup(id, "Bills", request);
  await createCategory(id, groupID, "Electric", request);

  // go to ledger page
  await page.goto(`/budget/${id}`);
  await page.getByRole("link", { name: "ledger" }).click();

  // open modal and add transaction
  await page.getByRole("link", { name: "Add" }).click();

  await page.getByLabel("Date").locator("visible=true").fill("2022-10-14");
  await page.getByLabel("Amount").locator("visible=true").fill("10.78");
  await selectOption(page, "Account", "Savings");

  // split category
  await page.getByRole("combobox", { name: "Category" }).click();
  await page.getByRole("button", { name: "Split" }).click();

  // add split details
  const split = page.getByRole("group", { name: "Split 1" });
  await split.getByLabel("Amount").fill("10.78");
  await split.getByLabel("Notes").fill(":)");
  await split.getByRole("combobox", { name: "Category" }).click();
  await page.getByRole("option", { name: "Electric" }).click();

  await page.getByRole("button", { name: "Save" }).click();
  await expect(page.getByRole("button", { name: "Save" })).toBeHidden();

  // transaction should be added (1 header row and 1 data row)
  await expect(page.getByRole("row")).toHaveCount(2);

  const checking = page
    .getByRole("row")
    .filter({ has: page.getByRole("cell").nth(4).getByText("Savings") })
    .getByRole("cell");
  await expect(checking.nth(1)).toHaveText("2022-10-14");
  await expect(checking.nth(2)).toHaveText("");
  await expect(checking.nth(3)).toHaveText("Split");
  await expect(checking.nth(4)).toHaveText("Savings");
  await expect(checking.nth(5)).toHaveText("");
  await expect(checking.nth(6)).toHaveText("$10.78");
});
