import { AddIcon } from "@chakra-ui/icons";
import {
  Button,
  chakra,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  useDisclosure,
  VStack,
} from "@chakra-ui/react";
import AccountSelect from "components/AccountSelect";
import CurrencyInput from "components/CurrencyInput";
import { Account } from "constants/types";
import { useAddTransaction } from "data/queries/transaction";
import { FormProvider, useForm } from "react-hook-form";

interface FormData {
  date: string;
  account: Account;
  amount: string;
  notes: string;
}

const AddTransaction = () => {
  const { register, handleSubmit, ...methods } = useForm<FormData>({
    mode: "onSubmit",
  });
  const { isOpen, onOpen, onClose } = useDisclosure();

  const { submit } = useAddTransaction();

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit((values) =>
      submit({
        accountID: values.account.id,
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
    <>
      <Button size="sm" rightIcon={<AddIcon />} onClick={onOpen}>
        Add
      </Button>

      <Drawer isOpen={isOpen} placement="right" onClose={onClose}>
        <DrawerOverlay />
        <DrawerContent>
          <DrawerHeader>Add Transaction</DrawerHeader>

          <chakra.form
            display="flex"
            h="full"
            flexDir="column"
            onSubmit={onSubmit}
          >
            <FormProvider
              {...methods}
              handleSubmit={handleSubmit}
              register={register}
            >
              <DrawerBody>
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

                  <FormControl isRequired>
                    <FormLabel>Amount</FormLabel>
                    <CurrencyInput name="amount" />
                  </FormControl>

                  <FormControl>
                    <FormLabel>Notes</FormLabel>
                    <Input {...register("notes")} maxLength={255} />
                  </FormControl>
                </VStack>
              </DrawerBody>

              <DrawerFooter>
                <Button type="submit" colorScheme="green">
                  Add
                </Button>
              </DrawerFooter>
            </FormProvider>
          </chakra.form>
        </DrawerContent>
      </Drawer>
    </>
  );
};

export default AddTransaction;
