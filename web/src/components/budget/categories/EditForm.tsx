import CurrencyInput from "@/components/CurrencyInput";
import { Amount } from "@/constants/types";
import { amountToFraction } from "@/data/format/amount";
import { useUpdateMonthCategory } from "@/data/queries/month";
import { Button, FormControl, FormLabel } from "@chakra-ui/react";
import { ReactNode } from "react";
import { FormProvider, useForm, useFormState } from "react-hook-form";

interface FormData {
  amount: string;
}

export const EditForm = ({
  children,
  categoryID,
  monthID,
  initialAmount,
  onSuccess,
}: {
  children: ReactNode;
  categoryID: string;
  monthID: string;
  initialAmount: Amount;
  onSuccess: () => void;
}) => {
  const { handleSubmit, ...methods } = useForm<FormData>({
    defaultValues: { amount: amountToFraction(initialAmount).toString() },
  });

  const { submit } = useUpdateMonthCategory({ monthID, categoryID });

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit((values) =>
      submit(values)
        .then(() => {
          onSuccess();
        })
        .catch((error) => {
          console.error(error);
        })
    )(event);

  return (
    <form onSubmit={onSubmit} autoComplete="off">
      <FormProvider {...methods} handleSubmit={handleSubmit}>
        {children}
      </FormProvider>
    </form>
  );
};

export const FormFields = () => {
  return (
    <>
      <FormControl>
        <FormLabel>Amount Assigned</FormLabel>
        <CurrencyInput name="amount" allowNegative={false} />
      </FormControl>
    </>
  );
};

export const FormSubmitButton = () => {
  const { isSubmitting } = useFormState();

  return (
    <Button
      type="submit"
      disabled={isSubmitting}
      isLoading={isSubmitting}
      colorScheme="green"
    >
      Save
    </Button>
  );
};
