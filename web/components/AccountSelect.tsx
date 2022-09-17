import Select, { useAsyncSelect } from "components/Select";
import { useAccounts } from "data/queries/account";

type Props = {
  name: string;
};

const AccountSelect = ({ name }: Props) => {
  const { isOpen, selectProps } = useAsyncSelect();
  const { accounts, isLoading } = useAccounts({ enabled: isOpen });

  return (
    <Select
      name={name}
      itemToString={(item) => item?.name ?? ""}
      itemToID={(item) => item?.id ?? ""}
      isLoading={isLoading}
      items={accounts}
      {...selectProps}
    />
  );
};

export default AccountSelect;
