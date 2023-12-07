import { getTransactions } from "$lib/api/requests/transaction";
import type { Transaction } from "$lib/types/transaction";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch, depends, params }) => {
  const transactions = await getTransactions({ fetch, params });

  depends("ledger:list");

  const transactionsByDate = transactions.reduce(
    (grouped, transaction) => {
      if (!grouped[transaction.date]) {
        grouped[transaction.date] = [];
      }

      grouped[transaction.date]?.push(transaction);

      return grouped;
    },
    {} as { [date: string]: Transaction[] },
  );

  return { transactions, transactionsByDate };
};
