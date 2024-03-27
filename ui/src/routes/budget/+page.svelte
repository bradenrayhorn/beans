<script lang="ts">
  import { paths } from "$lib/paths";
  import type { PageData } from "./$types";
  import IconAdd from "~icons/mdi/add";

  export let data: PageData;
</script>

<div class="w-full min-h-screen bg-base-300 flex flex-col justify-center">
  <div
    class="max-w-md w-full mx-auto p-14 bg-base-100 rounded-md shadow-md grow sm:grow-0"
  >
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

    <a href={paths.budgets.new} class="mt-8 btn btn-sm btn-primary"
      >Create Budget <IconAdd /></a
    >
  </div>
</div>
