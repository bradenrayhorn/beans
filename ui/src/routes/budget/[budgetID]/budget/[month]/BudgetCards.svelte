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

<h2 class="font-bold text-center w-full mt-6 mb-4">Expenses</h2>

<div class="flex flex-col w-full text-sm">
  {#each categoryGroups.filter((group) => !group.isIncome) as group (group.id)}
    <h3 class="font-bold text-center w-full mb-2">{group.name}</h3>

    {#each group.categories as category (category.id)}
      <a
        class="flex flex-col gap-2 w-full shrink-0 bg-base-100 shadow rounded p-4 mb-4"
        href={withParameter(paths.budget.budget.category, {
          ...$page.params,
          categoryID: category.id,
        })}
      >
        <div class="flex justify-between w-full items-center">
          <b>{category.name}</b>

          <span>
            Have {monthCategories[category.id]?.available?.display}
          </span>
        </div>

        <div class="flex justify-between w-full items-center">
          <span>
            Spent {monthCategories[category.id]?.activity?.display}
          </span>

          <span>
            +{monthCategories[category.id]?.assigned?.display}
          </span>
        </div>
      </a>
    {/each}
  {/each}
</div>

<h2 class="font-bold text-center w-full mt-6 mb-4">Income</h2>

<div class="flex flex-col w-full text-sm">
  {#each categoryGroups.filter((group) => group.isIncome) as group (group.id)}
    {#each group.categories as category (category.id)}
      <div
        class="flex flex-col gap-2 w-full shrink-0 bg-base-100 shadow rounded p-4 mb-4"
      >
        <div class="flex justify-between w-full items-center">
          <b>{category.name}</b>

          <span>
            {monthCategories[category.id]?.activity?.display}
          </span>
        </div>
      </div>
    {/each}
  {/each}
</div>
