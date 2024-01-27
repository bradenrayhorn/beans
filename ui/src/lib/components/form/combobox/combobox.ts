import {
  createCombobox as createComboboxMelt,
  type ComboboxOptionProps,
  type ComboboxOption,
} from "@melt-ui/svelte";
import { getContext, setContext } from "svelte";

const NAME = "combobox";

export type ComboboxCtx = ReturnType<typeof createComboboxCtx>;

export function createComboboxCtx(
  defaultSelected: undefined | ComboboxOptionProps<string>,
) {
  const combobox = createComboboxMelt<string>({
    defaultSelected: defaultSelected
      ? defaultSelected
      : (undefined as unknown as ComboboxOption<string>),
  });

  setContext(NAME, combobox);

  return combobox;
}

export function getComboboxCtx() {
  return getContext<ComboboxCtx>(NAME);
}
