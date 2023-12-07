<script lang="ts">
  import { page } from "$app/stores";
  import PageTitle from "$lib/components/PageTitle.svelte";
  import { paths, withParameter } from "$lib/paths";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<div class="p-4 md:max-w-md md:mx-auto">
  <PageTitle class="mb-8">Accounts</PageTitle>

  <div>
    <div class="flex flex-col gap-4" role="list">
      {#each data.accounts as account (account.id)}
        <div class="p-4 shadow-md rounded-md bg-base-100" role="listitem">
          <h2 class="font-bold text-lg">{account.name}</h2>
          <div>
            Balance: {account.balance.display}
          </div>
        </div>
      {/each}
    </div>

    {#if data.accounts.length === 0}
      <i>No accounts found.</i>
    {/if}

    <div class="mt-8">
      <a
        class="btn btn-primary btn-sm w-full"
        href={withParameter(paths.budget.accounts.new, $page.params)}
        >Add New Account</a
      >
    </div>
  </div>
</div>
