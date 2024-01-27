<script lang="ts">
  import { type ComboboxOptionProps } from "@melt-ui/svelte";
  import {
    ComboboxInput,
    ComboboxItem,
    ComboboxMenu,
    createComboboxCtx,
  } from "./combobox";
  import type { Payee } from "$lib/types/payee";

  const blankOption = { value: "", label: "None" };

  export let payees: Payee[];
  export let defaultPayee: Payee | null | undefined = undefined;

  const toOption = (payee: Payee): ComboboxOptionProps<string> => ({
    value: payee.id,
    label: payee.name,
  });

  const {
    states: { inputValue, touchedInput, selected },
  } = createComboboxCtx(defaultPayee ? toOption(defaultPayee) : blankOption);

  // Filter based on the input value
  $: filteredPayees = $touchedInput
    ? payees.filter((payee) =>
        payee.name.toLowerCase().includes($inputValue.toLowerCase()),
      )
    : payees;
</script>

<input name="payee_id" type="hidden" value={$selected?.value} />

<ComboboxInput label="Payee" />

<ComboboxMenu>
  <ComboboxItem item={blankOption} />

  {#each filteredPayees as payee (payee.id)}
    <ComboboxItem item={toOption(payee)} />
  {/each}
</ComboboxMenu>
