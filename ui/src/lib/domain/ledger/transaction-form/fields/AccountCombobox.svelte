<script lang="ts">
  import { melt, type ComboboxOptionProps } from "@melt-ui/svelte";
  import {
    ComboboxInput,
    ComboboxItem,
    ComboboxMenu,
    ComboboxNoResults,
    createComboboxCtx,
  } from "$lib/components/form/combobox";
  import type { Account } from "$lib/types/account";
  import { getTransactionFormCtx } from "../form-context";

  export let accounts: Account[];
  const { account, transferAccount } = getTransactionFormCtx();

  const toOption = (account: Account): ComboboxOptionProps<string> => ({
    value: account.id,
    label: account.name,
    disabled: $transferAccount?.id === account.id,
  });

  const {
    elements: { group, groupLabel },
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx($account ? toOption($account) : undefined);

  // sync to transaction ctx
  selected.subscribe((newValue) => {
    if ($account?.id !== newValue?.value) {
      account.update(() => accounts.find((a) => a.id === newValue?.value));
    }
  });

  // Filter based on the input value
  $: filteredAccounts = $touchedInput
    ? accounts.filter((account) =>
        account.name.toLowerCase().includes($inputValue.toLowerCase()),
      )
    : accounts;

  $: offBudget = filteredAccounts.filter((account) => account.offBudget);
  $: onBudget = filteredAccounts.filter((account) => !account.offBudget);
</script>

<input name="account_id" type="hidden" value={$selected?.value} />

<ComboboxInput label="Account" />

<ComboboxMenu>
  <div use:melt={$group("off-budget")}>
    <div use:melt={$groupLabel("off-budget")} class="font-bold py-2">
      Budgetable
    </div>

    {#each onBudget as account (account.id)}
      <ComboboxItem item={toOption(account)} />
    {/each}
  </div>

  <div use:melt={$group("off-budget")}>
    <div use:melt={$groupLabel("off-budget")} class="font-bold py-2">
      Off Budget
    </div>

    {#each offBudget as account (account.id)}
      <ComboboxItem item={toOption(account)} />
    {/each}
  </div>

  {#if filteredAccounts.length === 0}
    <ComboboxNoResults />
  {/if}
</ComboboxMenu>
