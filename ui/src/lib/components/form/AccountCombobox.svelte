<script lang="ts">
  import { melt, type ComboboxOptionProps } from "@melt-ui/svelte";
  import {
    ComboboxInput,
    ComboboxItem,
    ComboboxMenu,
    ComboboxNoResults,
    createComboboxCtx,
  } from "./combobox";
  import type { Account, RelatedAccount } from "$lib/types/account";

  export let accounts: Account[];
  export let defaultAccount: RelatedAccount | null | undefined = undefined;
  export let account: Account | undefined = undefined;

  const toOption = (
    account: Account | RelatedAccount,
  ): ComboboxOptionProps<string> => ({
    value: account.id,
    label: account.name,
  });

  const {
    elements: { group, groupLabel },
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx(defaultAccount ? toOption(defaultAccount) : undefined);

  $: account = accounts.find((a) => a.id === $selected?.value);

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
