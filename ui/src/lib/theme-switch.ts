export function setTheme(isDark: boolean) {
  document.documentElement.setAttribute(
    "data-theme",
    isDark ? "beansDark" : "beans",
  );
  localStorage.setItem("dark", isDark.toString());
}

// self-contained function to be embedded in <head> of page
export function initTheme() {
  const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  const savedValue = localStorage.getItem("dark");
  const hasSavedValue = savedValue !== null;

  function setTheme(isDark: boolean) {
    document.documentElement.setAttribute(
      "data-theme",
      isDark ? "beansDark" : "beans",
    );
    localStorage.setItem("dark", isDark.toString());
  }

  setTheme(hasSavedValue ? savedValue === "true" : mediaQuery.matches);
}
