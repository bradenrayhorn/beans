import type { Account } from "./account";
import type { Amount } from "./amount";
import type { Category } from "./category";
import type { Payee } from "./payee";

export type TransactionVariant =
  | "standard"
  | "off_budget"
  | "transfer"
  | "split";

export type Transaction = {
  id: string;
  account: Account;
  payee: Payee | null;
  category: Category | null;
  date: string;
  amount: Amount;
  notes: string | null;
  variant: TransactionVariant;
  transferAccount: Account | null;
};

export type Split = {
  id: string;
  category: Category;
  amount: Amount;
  notes: string | null;
};
