import { queries, queryKeys } from "@/constants/queries";
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
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

interface FormData {
  name: string;
}

function getHTTPErrorResponseMessage(error: unknown) {
  if (!error) {
    return "";
  }
  if (error instanceof HTTPError) {
    return error.message;
  }

  return "Unknown error.";
}

const CreateBudgetModal = ({ isOpen, onClose }: Props) => {
  const queryClient = useQueryClient();
  const {
    handleSubmit,
    register,
    formState: { isSubmitting },
    reset,
  } = useForm<FormData>();

  const {
    mutateAsync,
    error,
    reset: resetQuery,
  } = useMutation(queries.budget.create);
  const errorMessage = getHTTPErrorResponseMessage(error);

  useEffect(() => {
    if (isOpen) {
      reset();
      resetQuery();
    }
  }, [resetQuery, isOpen, reset]);

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit((v) => mutateAsync(v))(event)
      .then(() => {
        queryClient.invalidateQueries([queryKeys.budget.getAll]);
        onClose();
      })
      .catch((err) => {
        console.error(err);
      });

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>New Budget</ModalHeader>
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

export default CreateBudgetModal;
