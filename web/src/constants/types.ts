export interface User {
  id: string;
  username: string;
}

export interface Budget {
  id: string;
  name: string;
}

export interface Account {
  id: string;
  name: string;
}

export interface Amount {
  coefficient: number;
  exponent: number;
}

export interface Transaction {
  id: string;
  account: Account;
  date: string;
  amount: Amount;
  notes: string | null;
}

export interface CategoryGroup {
  id: string;
  name: string;
  categories: Array<Category>;
}

export interface Category {
  id: string;
  name: string;
}