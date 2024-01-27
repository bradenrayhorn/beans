<script lang="ts">
  import { melt, type ComboboxOption } from "@melt-ui/svelte";
  import { getComboboxCtx } from "./combobox";
  import IconCheck from "~icons/mdi/check";

  export let item: ComboboxOption<string>;

  const {
    elements: { option },
    helpers: { isSelected },
  } = getComboboxCtx();
</script>

<li
  use:melt={$option(item)}
  class="relative cursor-pointer scroll-my-2 rounded py-1 hover:bg-primary/20 data-[highlighted]:bg-primary/30"
>
  {#if $isSelected(item.value)}
    <div class="text-xs absolute left-0 top-2 z-10 bold">
      <IconCheck />
    </div>
  {/if}

  <div class="pl-5">
    <span class="text-sm">{item.label}</span>
  </div>
</li>
