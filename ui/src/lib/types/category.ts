export type Category = {
  id: string;
  name: string;
};

export type CategoryGroup = {
  id: string;
  name: string;
  isIncome: boolean;
  categories: Array<Category>;
};
