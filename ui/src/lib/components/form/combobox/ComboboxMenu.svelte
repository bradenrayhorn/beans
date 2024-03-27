<script lang="ts">
  import { melt } from "@melt-ui/svelte";
  import { getComboboxCtx } from "./combobox";
  import { fly } from "svelte/transition";

  const {
    elements: { menu },
    states: { open },
  } = getComboboxCtx();
</script>

{#if $open}
  <ul
    class="z-10 flex max-h-select flex-col overflow-hidden rounded bg-base-100 shadow-md"
    use:melt={$menu}
    transition:fly={{ duration: 150, y: -5 }}
  >
    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
    <div
      class="flex max-h-full flex-col gap-0 overflow-y-auto p-2 text-sm"
      tabindex="0"
    >
      <slot />
    </div>
  </ul>
{/if}
