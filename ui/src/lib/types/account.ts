import type { Amount } from "./amount";

export type AccountWithBalance = {
  id: string;
  name: string;
  balance: Amount;
  offBudget: boolean;
};

export type Account = {
  id: string;
  name: string;
  offBudget: boolean;
};
