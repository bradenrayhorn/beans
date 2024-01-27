<script lang="ts">
  import { melt, type ComboboxOptionProps } from "@melt-ui/svelte";
  import type { Category, CategoryGroup } from "$lib/types/category";
  import {
    ComboboxInput,
    ComboboxItem,
    ComboboxMenu,
    ComboboxNoResults,
    createComboboxCtx,
  } from "./combobox";

  export let categoryGroups: CategoryGroup[];
  export let defaultCategory: Category | null | undefined = undefined;

  const toOption = (category: Category): ComboboxOptionProps<string> => ({
    value: category.id,
    label: category.name,
  });

  const {
    elements: { groupLabel, group },
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx(
    defaultCategory ? toOption(defaultCategory) : undefined,
  );

  // Filter based on the input value
  $: filteredCategoryGroups = $touchedInput
    ? categoryGroups
        .map((group) => ({
          ...group,
          categories: group.categories.filter((category) =>
            category.name.toLowerCase().includes($inputValue.toLowerCase()),
          ),
        }))
        .filter((group) => group.categories.length > 0)
    : categoryGroups;
</script>

<input name="category_id" type="hidden" value={$selected?.value} />

<ComboboxInput label="Category" />

<ComboboxMenu>
  {#each filteredCategoryGroups as categoryGroup (categoryGroup.id)}
    <div use:melt={$group(categoryGroup.id)}>
      <div use:melt={$groupLabel(categoryGroup.id)} class="font-bold py-2">
        {categoryGroup.name}
      </div>

      {#each categoryGroup.categories as category (category.id)}
        <ComboboxItem item={toOption(category)} />
      {/each}
    </div>
  {:else}
    <ComboboxNoResults />
  {/each}
</ComboboxMenu>
