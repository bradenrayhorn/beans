<script lang="ts">
  import type { PageData } from "./$types";
  import { page } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import { selectedRows } from "../../../selected-state";
  import TransactionForm from "../../TransactionForm.svelte";
  import IconBack from "~icons/mdi/ChevronLeft";

  export let data: PageData;

  $: transactionID = $page.params.transactionID as string;

  $: $selectedRows = { [transactionID]: true };
  $: transaction = data.transactions.find((t) => t.id === transactionID);
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

<TransactionForm
  categories={data.categories}
  accounts={data.accounts}
  {transaction}
/>
