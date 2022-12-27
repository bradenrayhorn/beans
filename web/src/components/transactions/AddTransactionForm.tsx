import {
  Button,
  chakra,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import AccountSelect from "@/components/AccountSelect";
import CurrencyInput from "@/components/CurrencyInput";
import { Account, Category } from "@/constants/types";
import { useAddTransaction } from "@/data/queries/transaction";
import { PropsWithChildren } from "react";
import {
  FormProvider,
  useForm,
  useFormContext,
  useFormState,
} from "react-hook-form";
import CategorySelect from "../CategorySelect";

interface FormData {
  date: string;
  account: Account;
  category: Category;
  amount: string;
  notes: string;
}

const formID = "add-transaction-form";

export const AddTransactionFormFields = ({
  onClose,
}: {
  onClose: () => void;
}) => {
  const { register, handleSubmit } = useFormContext();
  const { submit } = useAddTransaction();

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit((values) =>
      submit({
        accountID: values.account.id,
        categoryID: values.category?.id,
        amount: values.amount,
        date: values.date,
        notes: values.notes ? values.notes : undefined,
      })
    )(event)
      .then(() => {
        onClose();
      })
      .catch((error) => {
        console.error(error);
      });

  return (
    <chakra.form
      id={formID}
      display="flex"
      h="full"
      flexDir="column"
      onSubmit={onSubmit}
    >
      <VStack spacing={6}>
        <FormControl isRequired>
          <FormLabel>Date</FormLabel>
          <Input {...register("date")} />
          <FormHelperText>Enter in format YYYY-MM-DD</FormHelperText>
        </FormControl>

        <FormControl isRequired>
          <FormLabel>Account</FormLabel>
          <AccountSelect name="account" />
        </FormControl>

        <FormControl>
          <FormLabel>Category</FormLabel>
          <CategorySelect name="category" />
        </FormControl>

        <FormControl isRequired>
          <FormLabel>Amount</FormLabel>
          <CurrencyInput name="amount" />
        </FormControl>

        <FormControl>
          <FormLabel>Notes</FormLabel>
          <Input {...register("notes")} maxLength={255} />
        </FormControl>
      </VStack>
    </chakra.form>
  );
};

export const AddTransactionFormSubmitButton = () => {
  const { isSubmitting } = useFormState();
  return (
    <Button
      type="submit"
      colorScheme="green"
      form={formID}
      disabled={isSubmitting}
      isLoading={isSubmitting}
    >
      Add
    </Button>
  );
};

export const AddTransactionForm = ({ children }: PropsWithChildren) => {
  const form = useForm<FormData>();
  return <FormProvider {...form}>{children}</FormProvider>;
};
