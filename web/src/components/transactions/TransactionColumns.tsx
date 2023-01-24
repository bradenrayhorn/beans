import { Transaction } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import { formatDate } from "@/data/format/date";
import { Input, InputProps } from "@chakra-ui/react";
import { ReactNode } from "react";
import { useFormContext } from "react-hook-form";
import AccountSelect from "../AccountSelect";
import CategorySelect from "../CategorySelect";
import CurrencyInput from "../CurrencyInput";
import DateInput from "../DateInput";

export type Column<Data> = {
  id: string;
  headerID: string;
  header: string;
  width?: string | number;
  cell: (data: Data) => string;
  isNumeric?: boolean;
  input?: () => ReactNode;
};

const FormInput = ({ name, ...props }: InputProps & { name: string }) => {
  const { register } = useFormContext();

  return <Input {...register(name)} size="sm" {...props} />;
};

export const columns: Column<Transaction>[] = [
  {
    id: "date",
    headerID: "header-date",
    header: "Date",
    width: 28,
    cell: (data) => formatDate(data.date),
    input: () => (
      <DateInput name="date" aria-labelledby="header-date" size="sm" />
    ),
  },
  {
    id: "category",
    headerID: "header-category",
    header: "Category",
    cell: (data) => data.category?.name ?? "",
    input: () => (
      <CategorySelect
        name="category"
        inputProps={{ "aria-labelledby": "header-category" }}
      />
    ),
  },
  {
    id: "account",
    headerID: "header-account",
    header: "Account",
    cell: (data) => data.account.name,
    input: () => (
      <AccountSelect
        name="account"
        inputProps={{ "aria-labelledby": "header-account" }}
      />
    ),
  },
  {
    id: "notes",
    headerID: "header-notes",
    header: "Notes",
    cell: (data) => data.notes ?? "",
    input: () => <FormInput name="notes" aria-labelledby="header-notes" />,
  },
  {
    id: "amount",
    headerID: "header-amount",
    header: "Amount",
    width: 28,
    isNumeric: true,
    cell: (data) => formatAmount(data.amount),
    input: () => (
      <CurrencyInput name="amount" aria-labelledby="header-amount" />
    ),
  },
];
