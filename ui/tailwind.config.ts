import type { Config } from "tailwindcss";

const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    extend: {
      maxHeight: {
        select: "18rem",
      },
    },
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    darkTheme: "beansDark",
    themes: [
      {
        beansLight: {
          primary: "#4b6bfb",
          secondary: "#7b92b2",
          accent: "#67cba0",
          neutral: "#181a2a",
          "base-100": "#ffffff",
          "base-200": "#F5F5F5",
          info: "#3abff8",
          success: "#36d399",
          warning: "#fbbd23",
          error: "#f87272",

          "--rounded-btn": "0.25rem",
        },
      },
      {
        beansDark: {
          primary: "#4b6bfb",
          secondary: "#7b92b2",
          accent: "#67cba0",
          neutral: "#181a2a",
          "base-100": "#ffffff",
          "base-200": "#F5F5F5",
          info: "#3abff8",
          success: "#36d399",
          warning: "#fbbd23",
          error: "#f87272",

          "--rounded-btn": "0.25rem",
        },
      },
    ],
  },
} satisfies Config;

export default config;
