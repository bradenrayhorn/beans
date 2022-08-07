import BudgetLayout from "components/layouts/BudgetLayout";
import { NextPageWithLayout } from "pages/_app";

const App: NextPageWithLayout = () => {
  return <main>Welcome to Beans.</main>;
};

App.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default App;
