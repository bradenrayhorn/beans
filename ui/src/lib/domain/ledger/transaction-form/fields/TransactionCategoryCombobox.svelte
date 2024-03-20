<script lang="ts">
  import type { CategoryGroup } from "$lib/types/category";
  import { generateId } from "@melt-ui/svelte/internal/helpers";
  import { getTransactionFormCtx } from "../form-context";
  import CategoryCombobox from "./CategoryCombobox.svelte";

  export let categoryGroups: CategoryGroup[];
  const { category, workingSplitIDs } = getTransactionFormCtx();

  const onSplit = () => {
    workingSplitIDs.update(() => [generateId()]);
  };
</script>

<CategoryCombobox
  {categoryGroups}
  defaultCategory={$category}
  bind:value={$category}
  canSplit={true}
  on:split={onSplit}
/>

{#if $category}
  <input name="category_id" type="hidden" value={$category.id} />
{/if}
