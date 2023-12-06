<script lang="ts">
  import { page } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<div class="text-sm breadcrumbs mb-8">
  <ul>
    <li>
      <a href={withParameter(paths.budget.settings.base, $page.params)}>
        Settings
      </a>
    </li>
    <li>Categories</li>
  </ul>
</div>

<div class="flex flex-col gap-12">
  {#each data.categoryGroups as categoryGroup (categoryGroup.id)}
    <a
      href={withParameter(paths.budget.settings.categories.group, {
        ...$page.params,
        categoryGroupID: categoryGroup.id,
      })}
      class="shadow bg-base-100 rounded-md p-4"
    >
      <h2 class="font-bold">{categoryGroup.name}</h2>
    </a>
  {/each}

  <a
    class="btn btn-primary md:w-64"
    href={withParameter(paths.budget.settings.categories.new, $page.params)}
  >
    Add New Group
  </a>
</div>
