<script lang="ts">
  import { page } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import type { CategoryGroup } from "$lib/types/category";
  import type { Month, MonthCategory } from "$lib/types/month";

  export let month: Month;

  $: monthCategories = month.categories.reduce(
    (map, monthCategory) => {
      map[monthCategory.categoryID] = monthCategory;
      return map;
    },
    {} as { [categoryID: string]: MonthCategory },
  );

  export let categoryGroups: Array<CategoryGroup>;
</script>

<div role="table" class="flex flex-col w-full text-sm" aria-label="Expenses">
  <div role="row" class="flex w-full bg-base-200 font-bold px-4 py-1">
    <div role="columnheader" class="flex-1">Category</div>
    <div role="columnheader" class="text-right flex-1">Assigned</div>
    <div role="columnheader" class="text-right flex-1">Spent</div>
    <div role="columnheader" class="text-right flex-1">Available</div>
  </div>

  {#each categoryGroups.filter((group) => !group.isIncome) as group (group.id)}
    <div role="rowgroup" aria-label={group.name}>
      <div aria-hidden="true" class="font-semibold text-lg mb-1 mt-8 px-4">
        {group.name}
      </div>

      <div role="rowgroup">
        {#each group.categories as category (category.id)}
          <a
            role="row"
            class="flex items-center w-full px-4 py-2 hover:bg-base-200 border-accent"
            class:border-y-2={$page.params.categoryID === category.id}
            href={withParameter(paths.budget.budget.category, {
              ...$page.params,
              categoryID: category.id,
            })}
          >
            <div role="cell" class="flex-1">
              {category.name}
            </div>

            <div role="cell" class="text-right flex-1">
              {monthCategories[category.id]?.assigned?.display}
            </div>

            <div role="cell" class="text-right flex-1">
              {monthCategories[category.id]?.activity?.display}
            </div>

            <div role="cell" class="text-right flex-1">
              {monthCategories[category.id]?.available?.display}
            </div>
          </a>
        {/each}
      </div>
    </div>
  {/each}
</div>

<div
  role="table"
  class="flex flex-col w-full text-sm mt-12"
  aria-label="Income"
>
  <div
    role="row"
    class="flex w-full bg-base-200 font-bold px-4 py-1 [&>div]:flex-1"
  >
    <div role="columnheader">Income</div>
    <div role="columnheader" class="text-right">Received</div>
  </div>

  <div role="rowgroup" class="px-4 mt-2">
    {#each categoryGroups.filter((group) => group.isIncome) as group (group.id)}
      {#each group.categories as category (category.id)}
        <div role="row" class="flex items-center w-full mt-1 [&>div]:flex-1">
          <div role="cell">
            {category.name}
          </div>

          <div role="cell" class="text-right">
            {monthCategories[category.id]?.activity?.display}
          </div>
        </div>
      {/each}
    {/each}
  </div>
</div>
