<script lang="ts">
  import { melt, type ComboboxOptionProps } from "@melt-ui/svelte";
  import {
    ComboboxInput,
    ComboboxItem,
    ComboboxMenu,
    createComboboxCtx,
  } from "$lib/components/form/combobox";
  import type { Payee } from "$lib/types/payee";
  import type { Account } from "$lib/types/account";
  import { getTransactionFormCtx } from "../form-context";

  const blankOption = { value: "", label: "None" };

  export let payees: Payee[];
  export let accounts: Account[];
  export let isDisabled: boolean = false;

  const {
    payee,
    account: selectedAccount,
    transferAccount,
  } = getTransactionFormCtx();

  const payeeToOption = (payee: Payee): ComboboxOptionProps<string> => ({
    value: `payee-${payee.id}`,
    label: payee.name,
  });
  const accountToOption = (account: Account): ComboboxOptionProps<string> => ({
    value: `account-${account.id}`,
    label: account.name,
    disabled: $selectedAccount?.id === account.id,
  });

  const defaultOption = $transferAccount
    ? accountToOption($transferAccount)
    : $payee
      ? payeeToOption($payee)
      : blankOption;

  const {
    elements: { group, groupLabel },
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx(defaultOption);

  // sync to transaction ctx
  selected.subscribe((newValue) => {
    if ($payee?.id !== newValue?.value?.replaceAll("payee-", "")) {
      payee.update(() =>
        payees.find(
          (payee) => payee.id === newValue?.value?.replaceAll("payee-", ""),
        ),
      );
    }
    if ($transferAccount?.id !== newValue?.value?.replaceAll("account-", "")) {
      transferAccount.update(() =>
        accounts.find(
          (account) =>
            account.id === newValue?.value?.replaceAll("account-", ""),
        ),
      );
    }
  });

  // Filter based on the input value
  $: filteredPayees = $touchedInput
    ? payees.filter((payee) =>
        payee.name.toLowerCase().includes($inputValue.toLowerCase()),
      )
    : payees;
  $: filteredAccounts = $touchedInput
    ? accounts.filter((account) =>
        account.name.toLowerCase().includes($inputValue.toLowerCase()),
      )
    : accounts;
</script>

{#if $transferAccount === undefined}
  <input
    name="payee_id"
    type="hidden"
    value={$selected?.value?.replaceAll("payee-", "")}
  />
{:else}
  <input name="transferAccountID" type="hidden" value={$transferAccount.id} />
{/if}

<ComboboxInput label="Payee" {isDisabled} />

<ComboboxMenu>
  <div use:melt={$group("payees")}>
    <div use:melt={$groupLabel("payees")} class="font-bold py-2">Payees</div>

    <ComboboxItem item={blankOption} />

    {#each filteredPayees as payee (payee.id)}
      <ComboboxItem item={payeeToOption(payee)} />
    {/each}
  </div>

  <div use:melt={$group("transfers")}>
    <div use:melt={$groupLabel("transfers")} class="font-bold py-2">
      Transfer
    </div>

    {#each filteredAccounts as account (account.id)}
      <ComboboxItem item={accountToOption(account)} />
    {/each}
  </div>
</ComboboxMenu>
