import { AddIcon } from "@chakra-ui/icons";
import { IconButton, useDisclosure } from "@chakra-ui/react";
import AddCategoryModal from "./AddCategoryModal";

const AddCategoryButton = ({ groupID }: { groupID: string }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton
        size="xs"
        onClick={() => onOpen()}
        icon={<AddIcon />}
        aria-label="Add category"
        variant="ghost"
      />
      <AddCategoryModal isOpen={isOpen} onClose={onClose} groupID={groupID} />
    </>
  );
};

export default AddCategoryButton;
