<script lang="ts">
  import { melt, type ComboboxOptionProps } from "@melt-ui/svelte";
  import type { Category, CategoryGroup } from "$lib/types/category";
  import {
    ComboboxInput,
    ComboboxItem,
    ComboboxMenu,
    ComboboxNoResults,
    createComboboxCtx,
  } from "$lib/components/form/combobox";
  import { getTransactionFormCtx } from "../form-context";

  export let categoryGroups: CategoryGroup[];
  const { category } = getTransactionFormCtx();

  const toOption = (category: Category): ComboboxOptionProps<string> => ({
    value: category.id,
    label: category.name,
  });

  const {
    elements: { groupLabel, group },
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx($category ? toOption($category) : undefined);

  // sync to transaction ctx
  selected.subscribe((newValue) => {
    if ($category?.id !== newValue?.value) {
      category.update(() =>
        categoryGroups
          .flatMap((group) => group.categories)
          .find((category) => category.id === newValue?.value),
      );
    }
  });

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
