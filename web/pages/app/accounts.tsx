import AppLayout from "components/layouts/AppLayout";
import { NextPageWithLayout } from "pages/_app";

const Accounts: NextPageWithLayout = () => {
  return <main>Accounts</main>;
};

Accounts.getLayout = (page) => <AppLayout>{page}</AppLayout>;

export default Accounts;
