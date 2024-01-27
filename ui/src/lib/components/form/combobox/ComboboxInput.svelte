<script lang="ts">
  import { melt } from "@melt-ui/svelte";
  import IconChevronDown from "~icons/mdi/ChevronDown";
  import IconChevronUp from "~icons/mdi/ChevronUp";
  import { getComboboxCtx } from "./combobox";

  export let label: string;

  const {
    elements: { label: labelEl, input },
    states: { open, selected, inputValue },
  } = getComboboxCtx();

  // Show value in input when it is box is not open
  $: if (!$open) {
    $inputValue = $selected?.value ? $selected.label ?? "" : "";
  } else {
    $inputValue = "";
  }
</script>

<div class="flex flex-col w-full">
  <!-- svelte-ignore a11y-label-has-associated-control -->
  <label use:melt={$labelEl} class="label label-text">{label}</label>

  <div class="relative w-full">
    <input
      use:melt={$input}
      class="input input-bordered cursor-pointer input-sm w-full"
    />

    <div
      class="absolute right-2 top-1/2 z-10 -translate-y-1/2 pointer-events-none"
    >
      {#if $open}
        <IconChevronUp />
      {:else}
        <IconChevronDown />
      {/if}
    </div>
  </div>
</div>
