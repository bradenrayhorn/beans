<script lang="ts">
  import { paths } from "$lib/paths";
  import type { PageData } from "./$types";
  import IconAdd from "~icons/mdi/add";

  export let data: PageData;
</script>

<div class="max-w-md w-full mx-auto pt-8 px-4">
  <h1 class="mb-8">Select a budget to continue.</h1>

  <div class="flex flex-col gap-2">
    {#each data.budgets as budget}
      <p>
        <a href={`/budget/${budget.id}`} class="btn btn-ghost btn-sm"
          >{budget.name}</a
        >
      </p>
    {/each}
  </div>

  {#if data.budgets.length === 0}
    <p><i>No existing budgets found.</i></p>
  {/if}

  <a href={paths.budgets.new} class="mt-8 btn btn-sm">New Budget <IconAdd /></a>
</div>
