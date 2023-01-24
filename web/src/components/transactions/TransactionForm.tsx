import { Account, Category } from "@/constants/types";
import { chakra } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { FormProvider, SubmitHandler, useForm } from "react-hook-form";

export interface FormData {
  date: string;
  account: Account;
  category: Category | null;
  amount: string;
  notes: string;
}

export const TransactionFormProvider = ({
  defaultValues,
  onSubmit,
  children,
}: PropsWithChildren & {
  defaultValues: FormData;
  onSubmit: SubmitHandler<FormData>;
}) => {
  const form = useForm<FormData>({ defaultValues });

  return (
    <FormProvider {...form}>
      <chakra.form
        display="flex"
        w="full"
        onSubmit={form.handleSubmit(onSubmit)}
      >
        {children}
      </chakra.form>
    </FormProvider>
  );
};
