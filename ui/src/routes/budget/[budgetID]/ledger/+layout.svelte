<script lang="ts">
  import { page } from "$app/stores";
  import PageTitle from "$lib/components/PageTitle.svelte";
  import { paths, withParameter } from "$lib/paths";
  import type { PageData } from "./$types";
  import MobileList from "./MobileList.svelte";
  import Table from "./Table.svelte";
  import IconAdd from "~icons/mdi/add";
  import { selectedRows } from "./selected-state";

  export let data: PageData;
  $: ({ transactions, transactionsByDate } = data);

  // un-select a row if it no longer exists
  $: Object.keys($selectedRows).forEach((rowID) => {
    if (!transactions.some((t) => t.id === rowID)) {
      selectedRows.update((current) => {
        delete current[rowID];
        return current;
      });
    }
  });

  $: isList =
    $page.url.pathname ===
    withParameter(paths.budget.ledger.base, $page.params);
  $: rowCount = Object.values($selectedRows).filter((x) => x).length;
  $: selectedRowID =
    rowCount !== 1
      ? undefined
      : Object.entries($selectedRows).find(([, isSelected]) => isSelected)?.[0];
</script>

<div class="flex h-full">
  <div class="grow flex-col flex">
    <PageTitle class="shrink-0 p-4">Ledger</PageTitle>

    {#if !isList}
      <div class="flex flex-col lg:hidden">
        <slot />
      </div>
    {/if}

    <div
      class="flex-col h-full"
      class:hidden={!isList}
      class:lg:flex={!isList}
      class:flex={isList}
    >
      <div class="shrink-0 flex w-full justify-end my-2 px-4">
        {#if isList}
          <div class="flex gap-4">
            {#if selectedRowID}
              <a
                href={withParameter(paths.budget.ledger.edit, {
                  ...$page.params,
                  transactionID: selectedRowID,
                })}
                class="btn btn-primary btn-xs btn-ghost hidden lg:inline-flex"
              >
                Edit
              </a>
            {/if}

            <a
              href={withParameter(paths.budget.ledger.new, $page.params)}
              class="btn btn-primary btn-xs btn-ghost"
            >
              Add <IconAdd />
            </a>
          </div>
        {:else}
          <span class="h-6"></span>
        {/if}
      </div>

      <div class="grow w-full bg-base-100 px-4 shadow">
        <div class="lg:hidden">
          <MobileList {transactionsByDate} />
        </div>

        <div class="hidden lg:block">
          <Table {transactions} disableSelection={!isList} />
        </div>
      </div>
    </div>
  </div>

  {#if !isList}
    <div class="hidden lg:block shrink-0 w-64 bg-base-100 h-full shadow-md p-4">
      <slot />
    </div>
  {/if}
</div>
