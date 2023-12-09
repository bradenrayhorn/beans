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
    <li>Payees</li>
  </ul>
</div>

<div class="shadow bg-base-100 rounded-md py-6">
  <div class="flex justify-between items-center px-4">
    <h2 class="font-bold mb-4">Payees</h2>

    <a
      href={withParameter(paths.budget.settings.payees.new, $page.params)}
      class="btn btn-sm btn-primary"
    >
      New Payee
    </a>
  </div>

  <div class="text-sm px-4">
    {#if data.payees.length === 0}
      <i>No payees found.</i>
    {/if}

    <ul class="flex flex-col gap-3">
      {#each data.payees as payee (payee.id)}
        <li>
          {payee.name}
        </li>
      {/each}
    </ul>
  </div>
</div>
