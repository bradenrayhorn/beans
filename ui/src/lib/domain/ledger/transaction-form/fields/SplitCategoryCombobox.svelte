<script lang="ts">
  import type { Category, CategoryGroup } from "$lib/types/category";
  import CategoryCombobox from "./CategoryCombobox.svelte";

  export let categoryGroups: CategoryGroup[];
  export let value: Category | undefined = undefined;
  export let initialValue: Category | undefined = undefined;
</script>

<CategoryCombobox {categoryGroups} defaultCategory={initialValue} bind:value />
