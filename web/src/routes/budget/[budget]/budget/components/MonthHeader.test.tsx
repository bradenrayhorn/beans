import MonthProvider from "@/context/MonthProvider";
import {
  currentDate,
  dateToString,
  formatBudgetMonth,
} from "@/data/format/date";
import { render } from "@/test/render";
import { api, server } from "@/test/setup";
import { screen, waitFor, within } from "@testing-library/react";
import dayjs from "dayjs";
import { rest } from "msw";
import { expect } from "vitest";
import MonthHeader from "./MonthHeader";

const defaultMonthHandler = rest.get(api("/api/v1/months/1"), (_, res, ctx) =>
  res(
    ctx.delay(20),
    ctx.json({
      data: {
        id: "1",
        date: dateToString(dayjs(currentDate()).startOf("month")),
      },
    })
  )
);

describe("MonthHeader", async () => {
  it("content hidden when loading", async () => {
    server.use(defaultMonthHandler);

    render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    const currentMonth = dayjs().format("YYYY.MM");

    // elements should be hidden
    expect(screen.queryByText(currentMonth)).not.toBeVisible();
    expect(
      screen.queryByRole("button", { name: "Next month" })
    ).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: "Previous month" })
    ).not.toBeInTheDocument();

    await isDoneLoading();

    // elements should be visible
    expect(screen.getByText(formatBudgetMonth(currentDate()))).toBeVisible();
    expect(
      screen.getByRole("button", { name: "Next month" })
    ).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Previous month" })
    ).toBeInTheDocument();
  });

  it("can navigate to next month", async () => {
    const nextMonth = dateToString(dayjs().add(1, "month"));
    server.use(
      defaultMonthHandler,
      rest.post(api("/api/v1/months"), (_, res, ctx) =>
        res(ctx.delay(50), ctx.json({ data: { month_id: "2" } }))
      ),
      rest.get(api("/api/v1/months/2"), (_, res, ctx) =>
        res(ctx.delay(20), ctx.json({ data: { id: "2", date: nextMonth } }))
      )
    );

    const { user } = render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    await isDoneLoading();

    await user.click(screen.getByRole("button", { name: "Next month" }));

    await isDoneLoading();

    expect(screen.getByText(formatBudgetMonth(nextMonth))).toBeVisible();
  });

  it("can navigate to previous month", async () => {
    const previousMonth = dateToString(dayjs().subtract(1, "month"));
    server.use(
      defaultMonthHandler,
      rest.post(api("/api/v1/months"), (_, res, ctx) => {
        return res(ctx.delay(50), ctx.json({ data: { month_id: "2" } }));
      }),
      rest.get(api("/api/v1/months/2"), (_, res, ctx) =>
        res(ctx.delay(20), ctx.json({ data: { id: "2", date: previousMonth } }))
      )
    );

    const { user } = render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    await isDoneLoading();

    await user.click(screen.getByRole("button", { name: "Previous month" }));

    await isDoneLoading();

    expect(screen.getByText(formatBudgetMonth(previousMonth))).toBeVisible();
  });
});

describe("MonthPicker", () => {
  it("can switch years", async () => {
    server.use(defaultMonthHandler);

    const { user } = render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    await isDoneLoading();

    await user.click(
      screen.getByRole("button", { name: formatBudgetMonth(currentDate()) })
    );

    expect(
      within(screen.getByRole("dialog")).queryByText(dayjs().format("YYYY"))
    ).toBeVisible();

    // navigate one year
    await user.click(screen.getByRole("button", { name: "Next year" }));

    expect(
      within(screen.getByRole("dialog")).queryByText(
        dayjs().add(1, "year").format("YYYY")
      )
    ).toBeVisible();

    // navigate back two years
    await user.click(screen.getByRole("button", { name: "Previous year" }));
    await user.click(screen.getByRole("button", { name: "Previous year" }));

    expect(
      within(screen.getByRole("dialog")).queryByText(
        dayjs().subtract(1, "year").format("YYYY")
      )
    ).toBeVisible();
  });

  it("forgets state when closed and reopened", async () => {
    server.use(defaultMonthHandler);

    const { user } = render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    await isDoneLoading();

    await user.click(
      screen.getByRole("button", { name: formatBudgetMonth(currentDate()) })
    );

    expect(
      within(screen.getByRole("dialog")).queryByText(dayjs().format("YYYY"))
    ).toBeVisible();

    // navigate one year
    await user.click(screen.getByRole("button", { name: "Next year" }));

    expect(
      within(screen.getByRole("dialog")).queryByText(
        dayjs().add(1, "year").format("YYYY")
      )
    ).toBeVisible();

    // close and reopen
    await user.click(
      screen.getByRole("button", { name: formatBudgetMonth(currentDate()) })
    );

    await waitFor(() =>
      expect(
        screen.queryByRole("button", { name: "Next year", hidden: true })
      ).not.toBeInTheDocument()
    );

    await user.click(
      screen.getByRole("button", { name: formatBudgetMonth(currentDate()) })
    );

    // should be back to current year
    expect(
      within(screen.getByRole("dialog")).queryByText(dayjs().format("YYYY"))
    ).toBeVisible();
  });

  it("can see active month", async () => {
    server.use(defaultMonthHandler);

    const { user } = render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    await isDoneLoading();

    await user.click(
      screen.getByRole("button", { name: formatBudgetMonth(currentDate()) })
    );

    const currentMonth = dayjs(currentDate()).format("MMM");

    expect(
      screen.getByRole("button", { name: currentMonth, current: true })
    ).toBeInTheDocument();

    // switch years and month should not be active
    await user.click(screen.getByRole("button", { name: "Next year" }));

    expect(
      screen.queryByRole("button", { name: currentMonth, current: true })
    ).not.toBeInTheDocument();
  });

  it("can select new month", async () => {
    const newDate = dateToString(dayjs().add(1, "year").startOf("year"));

    server.use(
      defaultMonthHandler,
      rest.post(api("/api/v1/months"), (_, res, ctx) => {
        return res(ctx.delay(50), ctx.json({ data: { month_id: "2" } }));
      }),
      rest.get(api("/api/v1/months/2"), (_, res, ctx) =>
        res(ctx.delay(20), ctx.json({ data: { id: "2", date: newDate } }))
      )
    );

    const { user } = render(
      <MonthProvider defaultMonthID="1">
        <MonthHeader />
      </MonthProvider>
    );

    await isDoneLoading();

    await user.click(
      screen.getByRole("button", { name: formatBudgetMonth(currentDate()) })
    );

    expect(
      within(screen.getByRole("dialog")).queryByText(dayjs().format("YYYY"))
    ).toBeVisible();

    // navigate one year
    await user.click(screen.getByRole("button", { name: "Next year" }));

    await user.click(screen.getByRole("button", { name: "Jan" }));

    // dialog should close
    await waitFor(() =>
      expect(
        screen.queryByRole("button", { name: "Next year", hidden: true })
      ).not.toBeInTheDocument()
    );

    // wait for new month to load
    await isDoneLoading();

    // new selection should be visible
    expect(screen.getByText(formatBudgetMonth(newDate))).toBeVisible();
  });
});

async function isDoneLoading() {
  await waitFor(() =>
    expect(
      screen.queryByRole("button", { name: "Next month" })
    ).toBeInTheDocument()
  );
}
