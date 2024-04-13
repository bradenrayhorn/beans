import type { Action } from "svelte/action";

export const keybind: Action<
  HTMLElement,
  { key: string; action: () => void }
> = (_node, options) => {
  function onKeyDown(event: KeyboardEvent) {
    if (event.key === options.key) {
      event.preventDefault();

      options.action();
    }
  }

  window.addEventListener("keydown", onKeyDown);

  return {
    destroy() {
      window.removeEventListener("keydown", onKeyDown);
    },
  };
};
