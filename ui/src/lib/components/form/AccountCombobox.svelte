<script lang="ts">
  import { type ComboboxOptionProps } from "@melt-ui/svelte";
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

  const toOption = (
    account: Account | RelatedAccount,
  ): ComboboxOptionProps<string> => ({
    value: account.id,
    label: account.name,
  });

  const {
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx(defaultAccount ? toOption(defaultAccount) : undefined);

  // Filter based on the input value
  $: filteredAccounts = $touchedInput
    ? accounts.filter((account) =>
        account.name.toLowerCase().includes($inputValue.toLowerCase()),
      )
    : accounts;
</script>

<input name="account_id" type="hidden" value={$selected?.value} />

<ComboboxInput label="Account" />

<ComboboxMenu>
  {#each filteredAccounts as account (account.id)}
    <ComboboxItem item={toOption(account)} />
  {:else}
    <ComboboxNoResults />
  {/each}
</ComboboxMenu>
