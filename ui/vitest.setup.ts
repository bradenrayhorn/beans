import { vi } from "vitest";
import "vitest-dom/extend-expect";
import "@testing-library/jest-dom";
import { readable } from "svelte/store";

vi.mock("$app/stores", () => {
  return {
    navigating: readable(),
    updated: readable(),
    page: readable({ form: undefined }),
  };
});

vi.mock("esm-env", () => {
  return {
    BROWSER: true,
    DEV: true,
  };
});
