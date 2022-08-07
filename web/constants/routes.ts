export const routes = {
  // auth
  login: "/login",
  defaultAfterAuth: "/budget",

  // budget
  budget: {
    noneSelected: "/budget",
    index: "/budget/[budget]",
  },
};

// all other routes require auth
export const unprotectedRoutes = [routes.login];

// these routes cannot be access if user is authed
export const forceUnproctedRoutes = [routes.login];
