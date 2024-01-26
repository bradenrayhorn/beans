<script lang="ts">
  import {
    createCombobox,
    melt,
    type ComboboxOptionProps,
    type ComboboxOption,
  } from "@melt-ui/svelte";
  import type { Category, CategoryGroup } from "$lib/types/category";
  import { fly } from "svelte/transition";
  import IconCheck from "~icons/mdi/check";

  export let categoryGroups: CategoryGroup[];
  export let defaultCategory: Category | null | undefined = undefined;

  const toOption = (category: Category): ComboboxOptionProps<string> => ({
    value: category.id,
    label: category.name,
  });

  const {
    elements: { label, input, menu, option, groupLabel, group },
    states: { open, inputValue, touchedInput, selected },
    helpers: { isSelected },
  } = createCombobox<string>({
    defaultSelected: defaultCategory
      ? toOption(defaultCategory)
      : (undefined as unknown as ComboboxOption<string>),
  });

  // Show value in input when it is box is not open
  $: if (!$open) {
    $inputValue = $selected?.label ?? "";
  }

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

<div class="flex flex-col w-full">
  <!-- svelte-ignore a11y-label-has-associated-control -->
  <label use:melt={$label} class="label label-text">Category</label>

  <div class="relative w-full">
    <input
      use:melt={$input}
      on:focus={(e) => e.currentTarget.select()}
      class="select select-bordered select-sm w-full"
    />
  </div>
</div>

{#if $open}
  <ul
    class="z-10 flex max-h-select flex-col overflow-hidden rounded-sm bg-white shadow"
    use:melt={$menu}
    transition:fly={{ duration: 150, y: -5 }}
  >
    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
    <div
      class="flex max-h-full flex-col gap-0 overflow-y-auto p-2 text-sm"
      tabindex="0"
    >
      {#each filteredCategoryGroups as categoryGroup (categoryGroup.id)}
        <div use:melt={$group(categoryGroup.id)}>
          <div use:melt={$groupLabel(categoryGroup.id)} class="font-bold py-2">
            {categoryGroup.name}
          </div>

          {#each categoryGroup.categories as category (category.id)}
            <li
              use:melt={$option(toOption(category))}
              class="relative cursor-pointer scroll-my-2 rounded py-1 hover:bg-primary/20 data-[highlighted]:bg-primary/30"
            >
              {#if $isSelected(category.id)}
                <div class="text-xs absolute left-0 top-2 z-10 bold">
                  <IconCheck />
                </div>
              {/if}

              <div class="pl-5">
                <span class="text-sm">{category.name}</span>
              </div>
            </li>
          {/each}
        </div>
      {:else}
        <div>No results found</div>
      {/each}
    </div>
  </ul>
{/if}
