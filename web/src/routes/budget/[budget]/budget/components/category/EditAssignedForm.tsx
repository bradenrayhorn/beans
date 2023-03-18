import CurrencyInput from "@/components/CurrencyInput";
import FormButton from "@/components/FormButton";
import { getHTTPErrorResponseMessage } from "@/constants/queries";
import { MonthCategory } from "@/constants/types";
import { useMonthID } from "@/context/MonthProvider";
import { amountToFraction } from "@/data/format/amount";
import { useUpdateMonthCategory } from "@/data/queries/month";
import {
  Alert,
  AlertIcon,
  Flex,
  FormControl,
  FormLabel,
} from "@chakra-ui/react";
import { HTTPError } from "ky";
import { FormEvent, RefObject } from "react";
import { FormProvider, useForm } from "react-hook-form";

interface FormData {
  amount: string;
}

export default function EditAssignedForm({
  inputRef,
  monthCategory,
  onClose,
}: {
  inputRef: RefObject<HTMLInputElement>;
  monthCategory?: MonthCategory;
  onClose: () => void;
}) {
  const form = useForm<FormData>({
    defaultValues: {
      amount: amountToFraction(monthCategory?.assigned).toString(),
    },
  });

  const monthID = useMonthID();
  const { submit, error } = useUpdateMonthCategory({
    monthID,
    categoryID: monthCategory?.category_id ?? "",
  });

  const onSubmit = (event: FormEvent<HTMLFormElement>) =>
    form.handleSubmit((values) =>
      submit(values)
        .then(() => {
          onClose();
        })
        .catch((err) => {
          if (!(err instanceof HTTPError)) {
            console.error(err);
          }
        })
    )(event);

  return (
    <form onSubmit={onSubmit} autoComplete="off">
      <FormProvider {...form}>
        <FormControl size="sm">
          <FormLabel>Assigned</FormLabel>
          <CurrencyInput name="amount" allowNegative={false} ref={inputRef} />
        </FormControl>

        <Flex justifyContent="flex-end" mt={4} gap={4} alignItems="flex-end">
          {!!error && (
            <Alert status="error" size="xs">
              <AlertIcon />
              {getHTTPErrorResponseMessage(error)}
            </Alert>
          )}

          <FormButton colorScheme="green" size="sm" type="submit">
            Save
          </FormButton>
        </Flex>
      </FormProvider>
    </form>
  );
}
