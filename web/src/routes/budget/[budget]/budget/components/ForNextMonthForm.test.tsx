import MonthProvider from "@/context/MonthProvider";
import { zeroAmount } from "@/data/format/amount";
import { render } from "@/test/render";
import { api, server } from "@/test/setup";
import { screen, waitFor } from "@testing-library/react";
import { rest } from "msw";
import { useRef } from "react";
import { expect } from "vitest";
import ForNextMonthForm from "./ForNextMonthForm";

describe("ForNextMonthForm", async () => {
  it("can save", async () => {
    server.use(
      rest.put(api("/api/v1/months/1"), (_, res, ctx) => res(ctx.delay(50)))
    );
    const { user } = render(<Form />);

    expect(screen.getByLabelText("Carryover")).toHaveValue("1");
    await user.clear(screen.getByLabelText("Carryover"));
    await user.type(screen.getByLabelText("Carryover"), "61");

    const saveButton = screen.getByRole("button", { name: "Save" });
    await user.click(saveButton);

    expect(saveButton).toBeDisabled();

    await waitFor(() => expect(saveButton).toBeEnabled());
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("handles api error", async () => {
    const invalidError = "Internal error.";

    server.use(
      rest.put(api("/api/v1/months/1"), (_, res, ctx) =>
        res(ctx.delay(50), ctx.status(400), ctx.json({ error: invalidError }))
      )
    );
    const { user } = render(<Form />);

    await user.clear(screen.getByLabelText("Carryover"));
    await user.type(screen.getByLabelText("Carryover"), "61");

    const saveButton = screen.getByRole("button", { name: "Save" });
    await user.click(saveButton);

    expect(saveButton).toBeDisabled();

    await waitFor(() => expect(saveButton).toBeEnabled());

    expect(screen.getByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent(invalidError);
  });
});

const Form = () => {
  const ref = useRef<HTMLInputElement>(null);

  return (
    <MonthProvider defaultMonthID="1">
      <ForNextMonthForm
        inputRef={ref}
        month={{
          id: "1",
          date: "2023-02-01",
          budgetable: zeroAmount,
          income: zeroAmount,
          assigned: zeroAmount,
          carryover: { exponent: 0, coefficient: 1 },
          carried_over: zeroAmount,
          categories: [],
        }}
        onClose={() => {}}
      />
    </MonthProvider>
  );
};
