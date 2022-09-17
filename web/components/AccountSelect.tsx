import Select from "components/Select";
import { useAccounts } from "data/queries/account";

type Props = {
  name: string;
};

const AccountSelect = ({ name }: Props) => {
  const { accounts, isLoading } = useAccounts();

  return (
    <Select
      name={name}
      itemToString={(item) => item?.name ?? ""}
      itemToID={(item) => item?.id ?? ""}
      isLoading={isLoading}
      items={accounts}
    />
  );
};

export default AccountSelect;
