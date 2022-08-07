import BudgetLayout from "components/layouts/BudgetLayout";
import { NextPageWithLayout } from "pages/_app";

const Accounts: NextPageWithLayout = () => {
  return <main>Accounts</main>;
};

Accounts.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default Accounts;
