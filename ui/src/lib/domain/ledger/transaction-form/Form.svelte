<script lang="ts">
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";
  import AmountInput from "./fields/AmountInput.svelte";
  import type { Account } from "$lib/types/account";
  import type { CategoryGroup } from "$lib/types/category";
  import type { Payee } from "$lib/types/payee";
  import type { Transaction } from "$lib/types/transaction";
  import AccountCombobox from "./fields/AccountCombobox.svelte";
  import CategoryCombobox from "./fields/CategoryCombobox.svelte";
  import PayeeCombobox from "./fields/PayeeCombobox.svelte";
  import Placeholder from "./fields/Placeholder.svelte";
  import { createTransactionFormCtx } from "./form-context";

  export let categoryGroups: Array<CategoryGroup>;
  export let accounts: Array<Account>;
  export let payees: Array<Payee>;
  export let transaction: Transaction | undefined = undefined;

  const { account, transferAccount } = createTransactionFormCtx(transaction);

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

  {#if $transferAccount !== undefined}
    <Placeholder field="Category" value="Transfer" />
  {:else if $account?.offBudget}
    <Placeholder field="Category" value="Off-Budget" />
  {:else}
    <CategoryCombobox {categoryGroups} />
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

  <div class="divider"></div>

  <div class="w-full flex flex-row justify-between">
    <SubmitButton class="btn btn-primary btn-sm" {isLoading}>Save</SubmitButton>
  </div>
</form>
