<script lang="ts">
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";
  import type { Account } from "$lib/types/account";
  import type { Category } from "$lib/types/category";
  import type { Payee } from "$lib/types/payee";
  import type { Transaction } from "$lib/types/transaction";

  export let categories: Array<Category>;
  export let accounts: Array<Account>;
  export let payees: Array<Payee>;
  export let transaction: Transaction | undefined = undefined;

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

  <label>
    <span class="label label-text">Payee</span>
    <select name="payee_id" class="select select-bordered select-sm w-full">
      <option value> </option>
      {#each payees as payee (payee.id)}
        <option value={payee.id} selected={payee.id === transaction?.payee?.id}>
          {payee.name}
        </option>
      {/each}
    </select>
  </label>

  <label>
    <span class="label label-text">Category</span>
    <select name="category_id" class="select select-bordered select-sm w-full">
      {#each categories as category (category.id)}
        <option
          value={category.id}
          selected={category.id === transaction?.category?.id}
        >
          {category.name}
        </option>
      {/each}
    </select>
  </label>

  <label>
    <span class="label label-text">Account</span>
    <select name="account_id" class="select select-bordered select-sm w-full">
      {#each accounts as account (account.id)}
        <option
          value={account.id}
          selected={account.id === transaction?.account?.id}
        >
          {account.name}
        </option>
      {/each}
    </select>
  </label>

  <label>
    <span class="label label-text">Notes</span>
    <input
      name="notes"
      type="text"
      class="input input-sm input-bordered w-full"
      value={transaction?.notes ?? ""}
    />
  </label>

  <label>
    <span class="label label-text">Amount</span>
    <input
      name="amount"
      type="text"
      class="input input-sm input-bordered w-full"
      value={transaction?.amount?.rawDisplay ?? ""}
    />
  </label>

  <div class="w-full flex flex-row justify-between">
    <SubmitButton class="btn btn-success btn-sm" {isLoading}>Save</SubmitButton>
  </div>
</form>
