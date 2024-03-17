<script lang="ts">
  import type { PageData } from "./$types";
  import { page, navigating } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import { selectedRows } from "../../../selected-state";
  import Form from "$lib/domain/ledger/transaction-form/Form.svelte";
  import IconBack from "~icons/mdi/ChevronLeft";
  import { enhance } from "$app/forms";
  import SubmitButton from "$lib/components/SubmitButton.svelte";

  export let data: PageData;

  $: transactionID = data.transaction.id;
  $: $selectedRows = { [transactionID]: true };

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="flex items-center">
  <a
    href={withParameter(paths.budget.ledger.base, $page.params)}
    class="btn btn-ghost btn-sm me-1"
    aria-label="Close form"
  >
    <IconBack />
  </a>

  <b>Edit Transaction</b>
</div>

<Form
  categoryGroups={data.categoryGroups}
  accounts={data.accounts}
  payees={data.payees}
  transaction={data.transaction}
/>

<form
  class="items-center mt-6"
  method="POST"
  action="?/delete"
  use:enhance={() => {
    isSubmitting = true;

    return async ({ update }) => {
      await update({ reset: false });

      isSubmitting = false;
    };
  }}
>
  <input class="hidden" type="text" name="ids[]" value={transactionID} />

  <SubmitButton class="btn btn-error btn-sm" {isLoading}
    >Delete Transaction</SubmitButton
  >
</form>
