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
  import { createEventDispatcher } from "svelte";

  export let categoryGroups: CategoryGroup[];
  export let defaultCategory: Category | undefined;
  export let value: Category | undefined;
  export let canSplit: boolean = false;

  const dispatch = createEventDispatcher();
  const onSplit = () => dispatch("split");

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

  // sync to transaction ctx
  selected.subscribe((newValue) => {
    if (value?.id !== newValue?.value) {
      value = categoryGroups
        .flatMap((group) => group.categories)
        .find((category) => category.id === newValue?.value);
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
  {/each}

  {#if filteredCategoryGroups.length > 0}
    {#if canSplit}
      <button class="btn btn-ghost btn-xs mt-4" on:click={onSplit}>Split</button
      >
    {/if}
  {:else}
    <ComboboxNoResults />
  {/if}
</ComboboxMenu>
