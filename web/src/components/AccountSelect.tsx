import Select, { SelectProps, useAsyncSelect } from "@/components/Select";
import { Account } from "@/constants/types";
import { useAccounts } from "@/data/queries/account";

const AccountSelect = (
  props: Omit<
    SelectProps<Account>,
    "itemToString" | "itemToID" | "isLoading" | "items" | "isClearable"
  >
) => {
  const { selectProps } = useAsyncSelect();
  const { accounts, isLoading } = useAccounts();

  return (
    <Select
      itemToString={(item) => item?.name ?? ""}
      itemToID={(item) => item?.id ?? ""}
      isLoading={isLoading}
      items={accounts}
      {...props}
      {...selectProps}
    />
  );
};

export default AccountSelect;
