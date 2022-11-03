import {
  createBrowserRouter,
  Navigate,
  Outlet,
  RouterProvider,
} from "react-router-dom";
import { routes } from "./constants/routes";
import { AuthProvider } from "./context/AuthContext";
import BudgetProvider from "./context/BudgetProvider";
import MustBeAuthed from "./routes/authed";
import BudgetsPage from "./routes/budget";
import BudgetHomePage from "./routes/budget/[budget]";
import AccountsPage from "./routes/budget/[budget]/accounts";
import BudgetPage from "./routes/budget/[budget]/budget";
import SettingsPage from "./routes/budget/[budget]/settings";
import TransactionsPage from "./routes/budget/[budget]/transactions";
import LoginPage from "./routes/login";
import MustBeNotAuthed from "./routes/not-authed";
import Root from "./routes/root";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      {
        path: "/",
        element: <Navigate to={routes.login} />,
      },
      {
        element: <MustBeNotAuthed />,
        children: [{ path: "login", element: <LoginPage /> }],
      },
      {
        element: <MustBeAuthed />,
        children: [
          {
            path: "/budget",
            element: <BudgetsPage />,
          },
          {
            path: "/budget/:budget",
            element: (
              <BudgetProvider>
                <Outlet />
              </BudgetProvider>
            ),
            children: [
              {
                path: "",
                element: <BudgetHomePage />,
              },
              {
                path: "accounts",
                element: <AccountsPage />,
              },
              {
                path: "budget",
                element: <BudgetPage />,
              },
              {
                path: "settings",
                element: <SettingsPage />,
              },
              {
                path: "transactions",
                element: <TransactionsPage />,
              },
            ],
          },
        ],
      },
    ],
  },
]);

function App() {
  return (
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  );
}

export default App;
