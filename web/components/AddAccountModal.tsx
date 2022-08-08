import {
  Button,
  FormControl,
  FormLabel,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
} from "@chakra-ui/react";
import { useAddAccount } from "data/queries/account";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

interface FormData {
  name: string;
}

const AddAccountModal = ({ isOpen, onClose }: Props) => {
  const {
    handleSubmit,
    register,
    formState: { isSubmitting },
    reset,
  } = useForm<FormData>();

  const { submit, errorMessage, reset: resetQuery } = useAddAccount();

  useEffect(() => {
    if (isOpen) {
      reset();
      resetQuery();
    }
  }, [reset, resetQuery, isOpen]);

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit(submit)(event)
      .then(() => {
        onClose();
      })
      .catch((error) => {
        console.error(error);
      });

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Add Account</ModalHeader>
        <ModalCloseButton />
        <form onSubmit={onSubmit}>
          <ModalBody>
            <FormControl isRequired>
              <FormLabel>Name</FormLabel>
              <Input {...register("name")} />
            </FormControl>

            {errorMessage && (
              <Text color="errorText" mt={2}>
                {errorMessage}
              </Text>
            )}
          </ModalBody>

          <ModalFooter>
            <Button
              type="submit"
              disabled={isSubmitting}
              isLoading={isSubmitting}
              colorScheme="green"
            >
              Save
            </Button>
          </ModalFooter>
        </form>
      </ModalContent>
    </Modal>
  );
};

export default AddAccountModal;
