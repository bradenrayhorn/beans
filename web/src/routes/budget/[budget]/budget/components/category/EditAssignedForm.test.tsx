import MonthProvider from "@/context/MonthProvider";
import { zeroAmount } from "@/data/format/amount";
import { render } from "@/test/render";
import { api, server } from "@/test/setup";
import { screen, waitFor } from "@testing-library/react";
import { rest } from "msw";
import { useRef } from "react";
import { expect } from "vitest";
import EditAssignedForm from "./EditAssignedForm";

describe("EditAssignedForm", async () => {
  it("can save", async () => {
    server.use(
      rest.post(api("/api/v1/months/1/categories"), (_, res, ctx) =>
        res(ctx.delay(50))
      )
    );
    const { user } = render(<Form />);

    await user.clear(screen.getByLabelText("Assigned"));
    await user.type(screen.getByLabelText("Assigned"), "61");

    const saveButton = screen.getByRole("button", { name: "Save" });
    await user.click(saveButton);

    expect(saveButton).toBeDisabled();

    await waitFor(() => expect(saveButton).toBeEnabled());
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("handles api error", async () => {
    const invalidError = "Internal error.";

    server.use(
      rest.post(api("/api/v1/months/1/categories"), (_, res, ctx) =>
        res(ctx.delay(50), ctx.status(400), ctx.json({ error: invalidError }))
      )
    );
    const { user } = render(<Form />);

    await user.clear(screen.getByLabelText("Assigned"));
    await user.type(screen.getByLabelText("Assigned"), "61");

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
      <EditAssignedForm
        inputRef={ref}
        monthCategory={{
          id: "1",
          category_id: "1",
          activity: zeroAmount,
          assigned: zeroAmount,
          available: zeroAmount,
        }}
        onClose={() => {}}
      />
    </MonthProvider>
  );
};
