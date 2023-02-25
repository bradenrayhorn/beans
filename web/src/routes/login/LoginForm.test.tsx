import { render } from "@/test/render";
import { api, server } from "@/test/setup";
import { screen, waitFor } from "@testing-library/react";
import { rest } from "msw";
import { expect } from "vitest";
import LoginForm from "./LoginForm";

describe("Login Form", async () => {
  it("can login", async () => {
    server.use(
      rest.post(api("/api/v1/user/login"), (_, res, ctx) =>
        res(ctx.delay(50), ctx.json({ id: "1", username: "user" }))
      )
    );
    const { user } = render(<LoginForm />, {
      routes: [{ path: "/budget", element: <>budget page</> }],
    });

    await user.type(screen.getByLabelText("Username"), "user");
    await user.type(screen.getByLabelText("Password"), "password");

    const loginButton = screen.getByRole("button", { name: /log in/i });
    await user.click(loginButton);

    expect(loginButton).toBeDisabled();

    await waitFor(() =>
      expect(screen.getByText("budget page")).toBeInTheDocument()
    );
  });

  it("handles api error", async () => {
    const invalidError = "Invalid username or password.";

    server.use(
      rest.post(api("/api/v1/user/login"), (_, res, ctx) =>
        res(ctx.delay(50), ctx.status(401), ctx.json({ error: invalidError }))
      )
    );
    const { user } = render(<LoginForm />);

    await user.type(screen.getByLabelText("Username"), "user");
    await user.type(screen.getByLabelText("Password"), "password");

    const loginButton = screen.getByRole("button", { name: /log in/i });
    await user.click(loginButton);

    expect(loginButton).toBeDisabled();

    await waitFor(() => expect(screen.getByRole("alert")).toBeInTheDocument());
    expect(screen.getByRole("alert")).toHaveTextContent(invalidError);
    expect(loginButton).toBeEnabled();
  });
});
