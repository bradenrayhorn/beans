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
} from "@chakra-ui/react";
import { useAddCategory } from "@/data/queries/category";
import { ReactNode } from "react";
import {
  FormProvider,
  useForm,
  useFormContext,
  useFormState,
} from "react-hook-form";

interface FormData {
  name: string;
}

const AddCategoryForm = ({
  children,
  groupID,
  onSuccess,
}: {
  children: ReactNode;
  groupID: string;
  onSuccess: () => void;
}) => {
  const { handleSubmit, ...methods } = useForm<FormData>();

  const { submit } = useAddCategory({ groupID });

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit(submit)(event)
      .then(() => {
        onSuccess();
      })
      .catch((error) => {
        console.error(error);
      });

  return (
    <form onSubmit={onSubmit}>
      <FormProvider {...methods} handleSubmit={handleSubmit}>
        {children}
      </FormProvider>
    </form>
  );
};

const FormFields = () => {
  const { register } = useFormContext();

  return (
    <>
      <FormControl isRequired>
        <FormLabel>Name</FormLabel>
        <Input {...register("name")} />
      </FormControl>
    </>
  );
};

const FormSubmitButton = () => {
  const { isSubmitting } = useFormState();

  return (
    <Button
      type="submit"
      disabled={isSubmitting}
      isLoading={isSubmitting}
      colorScheme="green"
    >
      Add
    </Button>
  );
};

type Props = {
  isOpen: boolean;
  onClose: () => void;
  groupID: string;
};

const AddCategoryModal = (props: Props) => {
  return (
    <Modal {...props}>
      <ModalOverlay />
      <ModalContent>
        <AddCategoryForm onSuccess={props.onClose} groupID={props.groupID}>
          <ModalHeader>Add Category</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <FormFields />
          </ModalBody>

          <ModalFooter>
            <FormSubmitButton />
          </ModalFooter>
        </AddCategoryForm>
      </ModalContent>
    </Modal>
  );
};

export default AddCategoryModal;
