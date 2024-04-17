import type { Config } from "tailwindcss";

const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    extend: {
      maxHeight: {
        select: "18rem",
      },
    },
    fontFamily: {
      sans: ["Inter", "ui-sans-serif"],
    },
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: [
      {
        beans: {
          primary: "#2F855A",
          "primary-content": "#ffffff",
          secondary: "#CBD5E0",
          accent: "#9AE6B4",
          neutral: "#374151",
          error: "#9B2C2C",
          success: "#48BB78",
          warning: "#ECC94B",
          "base-100": "#f9fafb",
          "base-200": "#f3f4f6",
          "base-300": "#e5e7eb",

          "--rounded-btn": "0.25rem",
        },
      },
      {
        beansDark: {
          primary: "#307050",
          "primary-content": "#ffffff",
          secondary: "#343b48",
          accent: "#9AE6B4",
          neutral: "#0E0E12",
          error: "#FC8181",
          success: "#48BB78",
          warning: "#ECC94B",
          "base-100": "#26262E",
          "base-200": "#202027",
          "base-300": "#1A1A20",

          "--rounded-btn": "0.25rem",
        },
      },
    ],
  },
} satisfies Config;

export default config;
