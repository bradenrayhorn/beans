import type { Account } from "$lib/types/account";
import type { Category } from "$lib/types/category";
import type { Payee } from "$lib/types/payee";
import type { Transaction } from "$lib/types/transaction";
import { generateId } from "@melt-ui/svelte/internal/helpers";
import { getContext, setContext } from "svelte";
import { writable } from "svelte/store";

const NAME = "transaction-form";

export type TransactionFormCtx = ReturnType<typeof createTransactionFormCtx>;

export function createTransactionFormCtx(
  transaction?: Transaction,
  splitCount?: number,
) {
  const transactionForm = {
    account: writable<Account | undefined>(transaction?.account),
    category: writable<Category | undefined>(
      transaction?.category ?? undefined,
    ),
    payee: writable<Payee | undefined>(transaction?.payee ?? undefined),
    transferAccount: writable<Account | undefined>(
      transaction?.transferAccount ?? undefined,
    ),
    workingSplitIDs: writable<Array<string>>(
      Array.from({ length: splitCount ?? 0 }, () => generateId()),
    ),
  };

  setContext(NAME, transactionForm);

  return transactionForm;
}

export function getTransactionFormCtx() {
  return getContext<TransactionFormCtx>(NAME);
}
