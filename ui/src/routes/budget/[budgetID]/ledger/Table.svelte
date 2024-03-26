<script lang="ts">
  import type { Transaction } from "$lib/types/transaction";
  import { selectedRows } from "./selected-state";

  export let transactions: Transaction[];
  export let disableSelection: boolean;
</script>

<table class="w-full text-sm border-collapse table-fixed">
  <thead>
    <tr
      class="text-left uppercase text-base-content-light [&>th]:border-b-base-200 [&>th]:border-b"
    >
      <th class="invisible w-8 overflow-hidden">Checkbox</th>
      <th class="w-28">Date</th>
      <th>Payee</th>
      <th>Category</th>
      <th>Account</th>
      <th>Notes</th>
      <th class="text-right w-28">Amount</th>
      <th class="invisible w-0 overflow-hidden">Edit Links</th>
    </tr>
  </thead>

  <tbody>
    {#each transactions as transaction (transaction.id)}
      <tr class="text-base-content [&>td]:py-2 relative">
        <td
          ><input
            type="checkbox"
            bind:checked={$selectedRows[transaction.id]}
            disabled={disableSelection}
            class="checkbox checkbox-xs checkbox-primary rounded-none block"
          /></td
        >
        <td>{transaction.date}</td>

        {#if transaction.transferAccount}
          <td class="pr-2 truncate"
            >{transaction.transferAccount?.name ?? ""}</td
          >
        {:else}
          <td class="pr-2 truncate" title={transaction.payee?.name ?? ""}
            >{transaction.payee?.name ?? ""}</td
          >
        {/if}

        {#if transaction.variant === "split"}
          <td class="pr-2 truncate italic">Split</td>
        {:else if transaction.variant === "off_budget"}
          <td class="pr-2 truncate italic">Off-Budget</td>
        {:else if transaction.variant === "transfer"}
          <td class="pr-2 truncate italic">Transfer</td>
        {:else}
          <td class="pr-2 truncate" title={transaction.category?.name ?? ""}
            >{transaction.category?.name ?? ""}</td
          >
        {/if}

        <td class="pr-2 truncate" title={transaction.account.name}
          >{transaction.account.name}</td
        >
        <td class="pr-2 truncate" title={transaction.notes ?? ""}
          >{transaction.notes ?? ""}</td
        >
        <td class="text-right pr-0 truncate">{transaction.amount.display}</td>
      </tr>
    {/each}
  </tbody>
</table>
