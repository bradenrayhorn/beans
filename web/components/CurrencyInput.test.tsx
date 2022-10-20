import { FormControl, FormLabel } from "@chakra-ui/react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import "@testing-library/jest-dom";
import { PropsWithChildren } from "react";
import { FormProvider, useForm } from "react-hook-form";
import CurrencyInput from "./CurrencyInput";

const Form = ({ children }: PropsWithChildren) => {
  const form = useForm();
  return (
    <FormProvider {...form}>
      <form>{children}</form>
    </FormProvider>
  );
};

const setup = () => {
  const user = userEvent.setup();
  render(
    <Form>
      <FormControl>
        <FormLabel>Amount</FormLabel>
        <CurrencyInput name="a" />
      </FormControl>
    </Form>
  );
  return [user];
};

const firstInputCases = [
  ["12", "12.00"],
  ["-12", "-12.00"],
  ["0", "0.00"],
  ["-0", "0.00"],
  ["4.31", "4.31"],
  ["-4.31", "-4.31"],
  ["10035", "10,035.00"],
  ["-10035", "-10,035.00"],
  ["4.3", "4.30"],
  ["4.567", ""],
  ["9999999999", "9,999,999,999.00"],
  ["-9999999999", "-9,999,999,999.00"],
  ["9999999999.1", ""],
  ["-9999999999.1", ""],
];

describe("handles first input", () => {
  test.each(firstInputCases)("input %s gives %s", async (input, result) => {
    const [user] = setup();

    await user.type(screen.getByLabelText("Amount"), input);
    await user.tab();

    expect(screen.getByLabelText("Amount")).toHaveValue(result);
  });
});

test("reverts to previous value when given invalid input", async () => {
  const [user] = setup();

  await user.type(screen.getByLabelText("Amount"), "12");
  await user.tab();

  await user.type(screen.getByLabelText("Amount"), "93jfsa9");
  await user.tab();

  expect(screen.getByLabelText("Amount")).toHaveValue("12.00");
});

test("reverts to basic value when focused", async () => {
  const [user] = setup();

  await user.type(screen.getByLabelText("Amount"), "12");
  await user.tab();
  expect(screen.getByLabelText("Amount")).toHaveValue("12.00");

  await user.tab({ shift: true });
  expect(screen.getByLabelText("Amount")).toHaveValue("12");
});
