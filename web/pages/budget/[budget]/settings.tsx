import BudgetLayout from "components/layouts/BudgetLayout";
import { NextPageWithLayout } from "pages/_app";

const Settings: NextPageWithLayout = () => {
  return <main>Budget settings</main>;
};

Settings.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default Settings;
