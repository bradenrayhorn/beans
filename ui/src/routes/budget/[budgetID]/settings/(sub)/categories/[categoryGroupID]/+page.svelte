<script lang="ts">
  import { paths, withParameter } from "$lib/paths";
  import { page } from "$app/stores";
  import type { PageData } from "./$types";

  export let data: PageData;

  $: categoryGroup = data.categoryGroups.find(
    (it) => it.id === $page.params.categoryGroupID,
  );
</script>

<div class="text-sm breadcrumbs mb-8">
  <ul>
    <li>
      <a href={withParameter(paths.budget.settings.base, $page.params)}>
        Settings
      </a>
    </li>
    <li>
      <a
        href={withParameter(
          paths.budget.settings.categories.base,
          $page.params,
        )}
      >
        Categories
      </a>
    </li>
    <li>{categoryGroup?.name}</li>
  </ul>
</div>

<div class="flex flex-col gap-6">
  {#each categoryGroup?.categories ?? [] as category (category.id)}
    <div class="shadow bg-base-100 rounded-md p-4">
      <h2 class="font-bold">{category.name}</h2>
    </div>
  {/each}

  <a
    class="btn btn-primary md:w-64"
    href={withParameter(
      paths.budget.settings.categories.newSubGroup,
      $page.params,
    )}
  >
    Add New Category
  </a>
</div>
