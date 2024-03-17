<script lang="ts">
  import type { PageData } from "./$types";
  import { page } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import { selectedRows } from "../../selected-state";
  import Form from "$lib/domain/ledger/transaction-form/Form.svelte";
  import IconBack from "~icons/mdi/ChevronLeft";

  export let data: PageData;

  $selectedRows = {};
</script>

<div class="flex items-center">
  <a
    href={withParameter(paths.budget.ledger.base, $page.params)}
    class="btn btn-ghost btn-sm me-1"
    aria-label="Close form"
  >
    <IconBack />
  </a>

  <b>Add Transaction</b>
</div>

<Form
  categoryGroups={data.categoryGroups}
  accounts={data.accounts}
  payees={data.payees}
/>
