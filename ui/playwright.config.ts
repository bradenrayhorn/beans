import type { PlaywrightTestConfig } from "@playwright/test";
import { devices } from "@playwright/test";

const ciConfig: PlaywrightTestConfig = {
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

const localConfig: PlaywrightTestConfig = {
  ...ciConfig,
  webServer: [],
  use: {
    baseURL: "http://localhost:5173",
  },
};

const config = process.env.TEST_ENV === "local" ? localConfig : ciConfig;

export default config;
