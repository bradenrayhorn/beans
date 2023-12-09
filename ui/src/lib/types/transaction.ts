import type { RelatedAccount } from "./account";
import type { Amount } from "./amount";
import type { Category } from "./category";
import type { Payee } from "./payee";

export type Transaction = {
  id: string;
  account: RelatedAccount;
  payee: Payee | null;
  category: Category | null;
  date: string;
  amount: Amount;
  notes: string | null;
};
