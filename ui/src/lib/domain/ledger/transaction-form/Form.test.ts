import type { Account } from "$lib/types/account";
import { Amount } from "$lib/types/amount";
import type { Category, CategoryGroup } from "$lib/types/category";
import type { Payee } from "$lib/types/payee";
import type { Split, Transaction } from "$lib/types/transaction";
import "@testing-library/jest-dom";
import { render, screen, within } from "@testing-library/svelte";
import userEvent from "@testing-library/user-event";
import { describe, expect, test } from "vitest";
import Form from "./Form.svelte";

const accounts: Account[] = [
  { id: "1", name: "Checking", offBudget: false },
  { id: "2", name: "Savings", offBudget: false },
  { id: "3", name: "401k", offBudget: true },
  { id: "4", name: "IRA", offBudget: true },
];
const categoryGroups: CategoryGroup[] = [
  {
    id: "1",
    name: "Bills",
    isIncome: false,
    categories: [{ id: "1", name: "Electric" }],
  },
];
const payees: Payee[] = [];

test("selecting off-budget account", async () => {
  const user = userEvent.setup();
  render(Form, { accounts, categoryGroups, payees });

  // select off-budget account, category should be disabled
  await user.click(screen.getByLabelText("Account"));
  await user.click(screen.getByRole("option", { name: "401k" }));

  expect(screen.getByLabelText("Category")).toBeDisabled();
  expect(screen.getByLabelText("Category")).toHaveValue("Off-Budget");

  // select another account, category should be enabled again
  await user.click(screen.getByLabelText("Account"));
  await user.click(screen.getByRole("option", { name: "Checking" }));

  expect(screen.getByLabelText("Category")).toBeEnabled();
});

describe("transfers", () => {
  test("disables category option", async () => {
    const user = userEvent.setup();
    render(Form, { accounts, categoryGroups, payees });

    // select account in payee dropdown, category should be disabled
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "Savings" }));

    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");

    // clear payee, category should be enabled again
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "None" }));

    expect(screen.getByLabelText("Category")).toBeEnabled();
  });

  test("cannot transfer to same account", async () => {
    const user = userEvent.setup();
    render(Form, { accounts, categoryGroups, payees });

    // select Savings as payee
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "Savings" }));

    // Savings should not be selectable as an account
    await user.click(screen.getByLabelText("Account"));
    expect(screen.getByRole("option", { name: "Savings" })).toHaveAttribute(
      "aria-disabled",
      "true",
    );

    // select Checking as account
    await user.click(screen.getByRole("option", { name: "Checking" }));

    // Checking should not be selectable as the payee
    await user.click(screen.getByLabelText("Payee"));
    expect(screen.getByRole("option", { name: "Checking" })).toHaveAttribute(
      "aria-disabled",
      "true",
    );
  });

  test("on-on budget", async () => {
    const user = userEvent.setup();
    render(Form, { accounts, categoryGroups, payees });

    // select Savings as payee
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "Savings" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");

    // select Checking as account
    await user.click(screen.getByLabelText("Account"));
    await user.click(screen.getByRole("option", { name: "Checking" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");
  });

  test("on-off budget", async () => {
    const user = userEvent.setup();
    render(Form, { accounts, categoryGroups, payees });

    // select Savings as payee
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "Savings" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");

    // select 401k as account
    await user.click(screen.getByLabelText("Account"));
    await user.click(screen.getByRole("option", { name: "401k" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Off-Budget");
  });

  test("off-off budget", async () => {
    const user = userEvent.setup();
    render(Form, { accounts, categoryGroups, payees });

    // select 401k as payee
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "401k" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");

    // select 401k as account
    await user.click(screen.getByLabelText("Account"));
    await user.click(screen.getByRole("option", { name: "IRA" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");
  });

  test("off-on budget", async () => {
    const user = userEvent.setup();
    render(Form, { accounts, categoryGroups, payees });

    // select 401k as payee
    await user.click(screen.getByLabelText("Payee"));
    await user.click(screen.getByRole("option", { name: "401k" }));

    // category is disabled
    expect(screen.getByLabelText("Category")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");

    // select Checking as account
    await user.click(screen.getByLabelText("Account"));
    await user.click(screen.getByRole("option", { name: "Checking" }));

    // can pick a category
    expect(screen.getByLabelText("Category")).toBeEnabled();
  });
});

test("can split transaction", async () => {
  const user = userEvent.setup();
  render(Form, { accounts, categoryGroups, payees });

  const parent = screen.getByRole("group", { name: "Parent Transaction" });

  // choose to split
  await user.click(within(parent).getByLabelText("Category"));
  await user.click(screen.getByRole("button", { name: "Split" }));

  // category is now split
  expect(within(parent).getByLabelText("Category")).toHaveValue("Split");
  expect(within(parent).getByLabelText("Category")).toBeDisabled();

  // split exists
  expect(screen.getByRole("group", { name: "Split 1" })).toBeInTheDocument();

  // add second split
  await user.click(screen.getByRole("button", { name: "Add" }));
  expect(screen.getByRole("group", { name: "Split 2" })).toBeInTheDocument();

  // remove both splits
  await user.click(screen.getByRole("button", { name: "Remove" }));
  await user.click(screen.getByRole("button", { name: "Remove" }));

  // splits are gone, category is editable
  expect(
    screen.queryByRole("group", { name: "Split 1" }),
  ).not.toBeInTheDocument();
  expect(
    screen.queryByRole("group", { name: "Split 2" }),
  ).not.toBeInTheDocument();
  expect(within(parent).getByLabelText("Category")).toHaveValue("");
  expect(within(parent).getByLabelText("Category")).toBeEnabled();
});

describe("loads existing transaction", () => {
  test("standard", async () => {
    const transaction: Transaction = {
      id: "1",
      account: { id: "1", name: "Checking", offBudget: false },
      category: { id: "1", name: "Bills" },
      payee: { id: "1", name: "Utility Co." },
      date: "2024-03-01",
      amount: new Amount("12.50"),
      notes: "Test notes",
      variant: "standard",
      transferAccount: null,
    };
    render(Form, { accounts, categoryGroups, payees, transaction });

    expect(screen.getByLabelText("Date")).toHaveValue("2024-03-01");
    expect(screen.getByLabelText("Payee")).toHaveValue("Utility Co.");
    expect(screen.getByLabelText("Account")).toHaveValue("Checking");
    expect(screen.getByLabelText("Category")).toHaveValue("Bills");
    expect(screen.getByLabelText("Notes")).toHaveValue("Test notes");
    expect(screen.getByLabelText("Amount")).toHaveValue("12.5");
  });

  test("off-budget", async () => {
    const transaction: Transaction = {
      id: "1",
      account: { id: "1", name: "Checking", offBudget: true },
      category: { id: "1", name: "Bills" },
      payee: { id: "1", name: "Utility Co." },
      date: "2024-03-01",
      amount: new Amount("12.50"),
      notes: "Test notes",
      variant: "standard",
      transferAccount: null,
    };
    render(Form, { accounts, categoryGroups, payees, transaction });

    expect(screen.getByLabelText("Date")).toHaveValue("2024-03-01");
    expect(screen.getByLabelText("Payee")).toHaveValue("Utility Co.");
    expect(screen.getByLabelText("Account")).toHaveValue("Checking");
    expect(screen.getByLabelText("Category")).toHaveValue("Off-Budget");
    expect(screen.getByLabelText("Notes")).toHaveValue("Test notes");
    expect(screen.getByLabelText("Amount")).toHaveValue("12.5");

    expect(screen.getByLabelText("Category")).toBeDisabled();
  });

  test("transfer", async () => {
    const transaction: Transaction = {
      id: "1",
      account: { id: "1", name: "Checking", offBudget: false },
      category: null,
      payee: null,
      date: "2024-03-01",
      amount: new Amount("12.50"),
      notes: "Test notes",
      variant: "transfer",
      transferAccount: { id: "2", name: "Savings", offBudget: false },
    };
    render(Form, { accounts, categoryGroups, payees, transaction });

    expect(screen.getByLabelText("Date")).toHaveValue("2024-03-01");
    expect(screen.getByLabelText("Payee")).toHaveValue("Savings");
    expect(screen.getByLabelText("Account")).toHaveValue("Checking");
    expect(screen.getByLabelText("Category")).toHaveValue("Transfer");
    expect(screen.getByLabelText("Notes")).toHaveValue("Test notes");
    expect(screen.getByLabelText("Amount")).toHaveValue("12.5");

    expect(screen.getByLabelText("Payee")).toBeDisabled();
    expect(screen.getByLabelText("Category")).toBeDisabled();
  });

  test("split", async () => {
    const movieCategory: Category = { id: "1", name: "Movie" };
    const transaction: Transaction = {
      id: "1",
      account: { id: "1", name: "Checking", offBudget: false },
      category: null,
      payee: null,
      date: "2024-03-01",
      amount: new Amount("12.50"),
      notes: "",
      variant: "split",
      transferAccount: null,
    };
    const splits: Array<Split> = [
      {
        id: "1",
        amount: new Amount("7.50"),
        category: movieCategory,
        notes: ":)",
      },
    ];
    render(Form, { accounts, categoryGroups, payees, transaction, splits });

    const parent = screen.getByRole("group", { name: "Parent Transaction" });
    const split = screen.getByRole("group", { name: "Split 1" });

    // category should have loaded as split
    expect(within(parent).getByLabelText("Category")).toHaveValue("Split");

    // split loaded
    expect(within(split).getByLabelText("Amount")).toHaveValue("7.5");
    expect(within(split).getByLabelText("Category")).toHaveValue("Movie");
    expect(within(split).getByLabelText("Notes")).toHaveValue(":)");

    // cannot add or remove splits on edit
    expect(
      screen.queryByRole("button", { name: "Add" }),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: "Remove" }),
    ).not.toBeInTheDocument();
  });
});
