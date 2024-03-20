<script lang="ts">
  import type { Category, CategoryGroup } from "$lib/types/category";
  import type { Split } from "$lib/types/transaction";
  import { generateId } from "@melt-ui/svelte/internal/helpers";
  import SplitCategoryCombobox from "./fields/SplitCategoryCombobox.svelte";

  export let categoryGroups: CategoryGroup[];
  export let initialSplit: Split | undefined;
  export let index: number;

  let amount = initialSplit?.amount?.rawDisplay ?? "";
  let notes = initialSplit?.notes ?? "";
  let category: Category | undefined = initialSplit?.category;

  $: splitValue = JSON.stringify({ amount, notes, category_id: category?.id });

  let id = generateId();
</script>

<input name="splits[json]" value={splitValue} type="hidden" />

<div class="flex flex-col gap-1" role="group" aria-labelledby={id}>
  <b {id}>Split {index + 1}</b>
  <label>
    <span class="label label-text">Amount</span>
    <input
      type="text"
      class="input input-sm input-bordered w-full"
      bind:value={amount}
    />
  </label>
  <SplitCategoryCombobox
    {categoryGroups}
    initialValue={category}
    bind:value={category}
  />
  <label>
    <span class="label label-text">Notes</span>
    <input
      type="text"
      class="input input-sm input-bordered w-full"
      bind:value={notes}
    />
  </label>
</div>
