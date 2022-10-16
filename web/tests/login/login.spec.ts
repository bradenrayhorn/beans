import { test, expect } from "@playwright/test";

test("can login", async ({ page }) => {
  await page.route("**/api/v1/user/me", (route) =>
    route.fulfill({
      status: 401,
      body: JSON.stringify({}),
    })
  );
  await page.goto("/");

  await page.route("**/api/v1/user/login", (route) =>
    route.fulfill({
      status: 200,
      body: JSON.stringify({ id: "1", username: "user" }),
    })
  );

  await page.getByLabel("Username").fill("user");
  await page.getByLabel("Password").fill("password");
  await page.getByRole("button", { name: "Log in" }).click();

  await page.route("**/api/v1/budgets", (route) =>
    route.fulfill({
      status: 200,
      body: JSON.stringify({ data: [{ id: "1", name: "Test Budget" }] }),
    })
  );

  await expect(page).toHaveURL("/budget");
});
