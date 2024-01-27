<script lang="ts">
  import { formatAmount } from "$lib/amount";
  import { generateId } from "@melt-ui/svelte/internal/helpers";
  import IconSpent from "~icons/mdi/chevron-double-down";
  import IconReceived from "~icons/mdi/chevron-double-up";

  export let defaultAmount: string | undefined | null;

  let id = generateId();
  let value = defaultAmount ?? "";
  $: parsedAmount = +value;
  $: formattedAmount = formatAmount(parsedAmount);
</script>

<div>
  <label for={id} class="label label-text">Amount</label>
  <input
    {id}
    name="amount"
    type="text"
    class="input input-sm input-bordered w-full"
    autocomplete="off"
    bind:value
  />

  <div class="text-xs min-h-6 mt-1 flex items-center gap-1">
    {#if parsedAmount > 0}
      <div class="p-0.5 rounded-full bg-success text-success-content text-sm">
        <IconReceived />
      </div>
      <span>You've received {formattedAmount}</span>
    {:else if parsedAmount < 0}
      <div class="p-0.5 rounded-full bg-warning text-warning-content text-sm">
        <IconSpent />
      </div>
      <span>You've spent {formattedAmount}</span>
    {:else}
      <span />
    {/if}
  </div>
</div>
