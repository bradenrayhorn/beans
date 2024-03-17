import "@testing-library/jest-dom";
import { test, expect } from "vitest";
import { render, screen } from "@testing-library/svelte";
import userEvent from "@testing-library/user-event";
import AmountInput from "./AmountInput.svelte";

test("shows default value when passed", () => {
  render(AmountInput, { defaultAmount: "10.61" });

  expect(screen.getByLabelText("Amount")).toHaveValue("10.61");
});

test("shows helper when a positive value is entered", async () => {
  const user = userEvent.setup();
  render(AmountInput);

  await user.type(screen.getByLabelText("Amount"), "1500.581");

  expect(screen.getByText("You've received $1,500.58")).toBeInTheDocument();
});

test("shows helper when a negative value is entered", async () => {
  const user = userEvent.setup();
  render(AmountInput);

  await user.type(screen.getByLabelText("Amount"), "1500.581");

  expect(screen.getByText("You've received $1,500.58")).toBeInTheDocument();
});

test("shows no helper when 0, blank, or unknown", async () => {
  const user = userEvent.setup();
  render(AmountInput);

  await user.type(screen.getByLabelText("Amount"), "0");
  expect(screen.queryByText("You've")).not.toBeInTheDocument();

  await user.type(screen.getByLabelText("Amount"), "abcdef");
  expect(screen.queryByText("You've")).not.toBeInTheDocument();

  await user.clear(screen.getByLabelText("Amount"));
  expect(screen.queryByText("You've")).not.toBeInTheDocument();
});
