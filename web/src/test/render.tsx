import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, render } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  createMemoryRouter,
  RouteObject,
  RouterProvider,
} from "react-router-dom";
import { afterEach } from "vitest";

afterEach(() => {
  cleanup();
});

const queryClient = new QueryClient({
  // do not log as tests intentionally trigger errors, this pollutes the log
  logger: {
    log: () => {},
    warn: () => {},
    error: () => {},
  },
});

const customRender = (
  ui: React.ReactElement,
  options: { routes?: RouteObject[] } = {}
) => {
  const router = createMemoryRouter([
    {
      path: "/",
      element: ui,
    },

    ...(options.routes ?? []),
  ]);

  const user = userEvent.setup();

  return {
    user,
    ...render(<RouterProvider router={router} />, {
      wrapper: ({ children }) => (
        <QueryClientProvider client={queryClient}>
          {children}
        </QueryClientProvider>
      ),
      ...options,
    }),
  };
};

export { customRender as render };
