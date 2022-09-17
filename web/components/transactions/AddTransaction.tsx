import { AddIcon } from "@chakra-ui/icons";
import {
  Button,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import {
  AddTransactionForm,
  AddTransactionFormFields,
  AddTransactionFormSubmitButton,
} from "components/transactions/AddTransactionForm";

const AddTransaction = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Button size="sm" rightIcon={<AddIcon />} onClick={onOpen}>
        Add
      </Button>

      <Drawer isOpen={isOpen} placement="right" onClose={onClose}>
        <DrawerOverlay />
        <DrawerContent>
          <DrawerHeader>Add Transaction</DrawerHeader>

          <AddTransactionForm>
            <DrawerBody>
              <AddTransactionFormFields onClose={onClose} />
            </DrawerBody>

            <DrawerFooter>
              <AddTransactionFormSubmitButton />
            </DrawerFooter>
          </AddTransactionForm>
        </DrawerContent>
      </Drawer>
    </>
  );
};

export default AddTransaction;
