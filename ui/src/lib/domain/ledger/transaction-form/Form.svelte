<script lang="ts">
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";
  import AmountInput from "./fields/AmountInput.svelte";
  import type { Account } from "$lib/types/account";
  import type { CategoryGroup } from "$lib/types/category";
  import type { Payee } from "$lib/types/payee";
  import type { Split, Transaction } from "$lib/types/transaction";
  import AccountCombobox from "./fields/AccountCombobox.svelte";
  import PayeeCombobox from "./fields/PayeeCombobox.svelte";
  import Placeholder from "./fields/Placeholder.svelte";
  import { createTransactionFormCtx } from "./form-context";
  import TransactionCategoryCombobox from "./fields/TransactionCategoryCombobox.svelte";
  import SplitFields from "./SplitFields.svelte";
  import { generateId } from "@melt-ui/svelte/internal/helpers";

  export let categoryGroups: Array<CategoryGroup>;
  export let accounts: Array<Account>;
  export let payees: Array<Payee>;
  export let transaction: Transaction | undefined = undefined;
  export let splits: Array<Split> = [];

  const { account, transferAccount, workingSplitIDs } =
    createTransactionFormCtx(transaction, splits.length);

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<FormError />

<form
  class="flex flex-col gap-2"
  method="POST"
  action="?/save"
  use:enhance={() => {
    isSubmitting = true;

    return async ({ update }) => {
      await update();

      isSubmitting = false;
    };
  }}
>
  <div class="flex flex-col gap-2" role="group" aria-label="Parent Transaction">
    <label>
      <span class="label label-text">Date</span>
      <input
        name="date"
        type="text"
        class="input input-sm input-bordered w-full"
        value={transaction?.date ?? ""}
      />
    </label>

    <PayeeCombobox
      {payees}
      {accounts}
      isDisabled={!!transaction?.transferAccount}
    />

    {#if $workingSplitIDs.length > 0}
      <Placeholder field="Category" value="Split" />
    {:else if $transferAccount && ($transferAccount.offBudget === $account?.offBudget || !$account)}
      <Placeholder field="Category" value="Transfer" />
    {:else if $account?.offBudget}
      <Placeholder field="Category" value="Off-Budget" />
    {:else}
      <TransactionCategoryCombobox {categoryGroups} />
    {/if}

    <AccountCombobox {accounts} />

    <label>
      <span class="label label-text">Notes</span>
      <input
        name="notes"
        type="text"
        class="input input-sm input-bordered w-full"
        value={transaction?.notes ?? ""}
      />
    </label>

    <AmountInput defaultAmount={transaction?.amount?.rawDisplay} />

    {#if $workingSplitIDs.length > 0}
      <b>Splits:</b>

      {#if splits.length === 0}
        <div class="flex space-between">
          <button
            class="btn btn-ghost btn-xs"
            type="button"
            on:click={() => {
              workingSplitIDs.update((cur) => [...cur, generateId()]);
            }}>Add</button
          >
          <button
            class="btn btn-ghost btn-xs"
            type="button"
            on:click={() => {
              workingSplitIDs.update((cur) => cur.slice(0, cur.length - 1));
            }}>Remove</button
          >
        </div>
      {/if}
    {/if}
  </div>

  {#each $workingSplitIDs as id, i (id)}
    <SplitFields {categoryGroups} index={i} initialSplit={splits[i]} />
  {/each}

  <div class="divider"></div>

  <div class="w-full flex flex-row justify-between">
    <SubmitButton class="btn btn-primary btn-sm" {isLoading}>Save</SubmitButton>
  </div>
</form>
