import { Amount, Category } from "@/constants/types";
import { EditIcon } from "@chakra-ui/icons";
import {
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  IconButton,
  useDisclosure,
} from "@chakra-ui/react";
import { EditForm, FormFields, FormSubmitButton } from "./EditForm";

export default function EditButton({
  category,
  monthID,
  amount,
}: {
  category: Category;
  monthID: string;
  amount: Amount;
}) {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton
        aria-label={`Edit ${category.name}`}
        icon={<EditIcon />}
        variant="ghost"
        onClick={() => onOpen()}
      />

      <Drawer isOpen={isOpen} placement="right" onClose={onClose}>
        <DrawerOverlay />
        <DrawerContent>
          <DrawerHeader>Edit {category.name}</DrawerHeader>

          <EditForm
            categoryID={category.id}
            monthID={monthID}
            initialAmount={amount}
            onSuccess={() => onClose()}
          >
            <DrawerBody>
              <FormFields />
            </DrawerBody>

            <DrawerFooter>
              <FormSubmitButton />
            </DrawerFooter>
          </EditForm>
        </DrawerContent>
      </Drawer>
    </>
  );
}
