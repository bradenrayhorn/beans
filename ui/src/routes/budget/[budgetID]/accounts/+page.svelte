<script lang="ts">
  import { page } from "$app/stores";
  import PageTitle from "$lib/components/PageTitle.svelte";
  import { paths, withParameter } from "$lib/paths";
  import type { PageData } from "./$types";

  export let data: PageData;

  $: accounts = data.accounts.sort((a, b) =>
    a.offBudget === b.offBudget ? 0 : a.offBudget ? 1 : -1,
  );
</script>

<div class="p-4 grow md:p-12 md:max-w-md md:mx-auto">
  <PageTitle class="mb-8">Accounts</PageTitle>

  <div>
    <div class="flex flex-col gap-4" role="list">
      {#each accounts as account (account.id)}
        <div class="p-6 shadow-md rounded-md bg-base-100" role="listitem">
          <h2 class="font-bold text-lg mb-4">{account.name}</h2>
          <div>
            Balance: {account.balance.display}
          </div>
          {#if account.offBudget}
            <div class="italic">Off-Budget</div>
          {/if}
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
