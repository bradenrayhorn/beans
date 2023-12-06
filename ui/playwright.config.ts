import type { PlaywrightTestConfig } from "@playwright/test";
import { devices } from "@playwright/test";

const config: PlaywrightTestConfig = {
  webServer: [
    {
      command: "npm run build && npm run preview",
      port: 4173,
      reuseExistingServer: !process.env.CI,
      env: {
        PUBLIC_BASE_API_URL: "",
      },
    },
    {
      command: "(cd ../ && make run)",
      url: "http://localhost:8000/health-check",
      reuseExistingServer: !process.env.CI,
    },
  ],
  use: {
    baseURL: "http://localhost:4173",
  },
  testDir: "tests",
  testMatch: /(.+\.)?(test|spec)\.[jt]s/,
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
      },
    },
  ],
};

export default config;
