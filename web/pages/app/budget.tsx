import AppLayout from "components/layouts/AppLayout";
import { NextPageWithLayout } from "pages/_app";

const Budget: NextPageWithLayout = () => {
  return <main>Budget</main>;
};

Budget.getLayout = (page) => <AppLayout>{page}</AppLayout>;

export default Budget;
