import BudgetLayout from "components/layouts/BudgetLayout";
import { NextPageWithLayout } from "pages/_app";

const Budget: NextPageWithLayout = () => {
  return <main>Budget</main>;
};

Budget.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default Budget;
