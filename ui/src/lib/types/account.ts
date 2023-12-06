import type { Amount } from "./amount";

export type Account = {
  id: string;
  name: string;
  balance: Amount;
};

export type RelatedAccount = {
  id: string;
  name: string;
};
