import type { Amount } from "./amount";

export type Month = {
  id: string;
  date: string;
  budgetable: Amount;
  carryover: Amount;
  income: Amount;
  assigned: Amount;
  carriedOver: Amount;
  categories: Array<MonthCategory>;
};

export type MonthCategory = {
  id: string;
  categoryID: string;
  assigned: Amount;
  activity: Amount;
  available: Amount;
};
