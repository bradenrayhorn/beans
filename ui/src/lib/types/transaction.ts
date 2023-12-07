import type { RelatedAccount } from "./account";
import type { Amount } from "./amount";
import type { Category } from "./category";

export type Transaction = {
  id: string;
  account: RelatedAccount;
  category: Category | null;
  date: string;
  amount: Amount;
  notes: string | null;
};
