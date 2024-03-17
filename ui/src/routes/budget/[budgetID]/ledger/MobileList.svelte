<script lang="ts">
  import { page } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import type { Transaction } from "$lib/types/transaction";

  export let transactionsByDate: { [date: string]: Transaction[] };
</script>

<div class="pb-8">
  {#each Object.entries(transactionsByDate) as [date, transactions] (date)}
    <div class="divider text-center w-full py-8">
      <b>{date}</b>
    </div>

    <div class="flex flex-col gap-4">
      {#each transactions as transaction (transaction.id)}
        <a
          href={withParameter(paths.budget.ledger.edit, {
            ...$page.params,
            transactionID: transaction.id,
          })}
        >
          <div class="flex justify-between">
            {#if transaction.variant === "transfer"}
              <div>
                {transaction.transferAccount?.name ?? ""}
              </div>
            {:else}
              <div>
                {transaction.payee?.name ?? ""}
              </div>
            {/if}
            <div>
              {transaction.amount.display}
            </div>
          </div>

          <div class="flex justify-between">
            {#if transaction.variant === "off_budget"}
              <div class="italic">Off-Budget</div>
            {:else if transaction.variant === "transfer"}
              <div class="italic">Transfer</div>
            {:else}
              <div>
                {transaction.category?.name ?? ""}
              </div>
            {/if}
            <div>
              {transaction.account.name ?? ""}
            </div>
          </div>
        </a>
      {/each}
    </div>
  {/each}
</div>
