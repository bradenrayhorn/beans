import CurrencyInput from "@/components/CurrencyInput";
import FormButton from "@/components/FormButton";
import { getHTTPErrorResponseMessage } from "@/constants/queries";
import { Month } from "@/constants/types";
import { amountToFraction } from "@/data/format/amount";
import { useUpdateMonth } from "@/data/queries/month";
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
  carryover: string;
}

export default function ForNextMonthForm({
  inputRef,
  month,
  onClose,
}: {
  inputRef: RefObject<HTMLInputElement>;
  month?: Month;
  onClose: () => void;
}) {
  const form = useForm<FormData>({
    defaultValues: {
      carryover: amountToFraction(month?.carryover).toString(),
    },
  });

  const { submit, error } = useUpdateMonth({ monthID: month?.id ?? "" });

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
          <FormLabel>Carryover</FormLabel>
          <CurrencyInput
            name="carryover"
            allowNegative={false}
            ref={inputRef}
          />
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
