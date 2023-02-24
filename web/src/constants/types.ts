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
  category: Category | null;
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

export interface MonthCategory {
  id: string;
  assigned: Amount;
  activity: Amount;
  available: Amount;
  category_id: string;
}

export interface Month {
  id: string;
  date: string;
  categories: MonthCategory[];
}
