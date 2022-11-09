import {
  Button,
  FormControl,
  FormLabel,
  HStack,
  Input,
} from "@chakra-ui/react";
import { useAddCategoryGroup } from "@/data/queries/category";
import { useForm } from "react-hook-form";

interface FormData {
  name: string;
}
const AddCategoryGroup = () => {
  const {
    handleSubmit,
    register,
    formState: { isSubmitting },
    reset: resetForm,
  } = useForm<FormData>();

  const { submit, reset: resetQuery } = useAddCategoryGroup();

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) =>
    handleSubmit(submit)(event)
      .then(() => {
        resetForm();
        resetQuery();
      })
      .catch((error) => {
        console.error(error);
      });

  return (
    <form onSubmit={onSubmit} aria-label="Add category group">
      <HStack align="flex-end">
        <FormControl isRequired>
          <FormLabel>Name</FormLabel>
          <Input {...register("name")} />
        </FormControl>
        <Button type="submit" disabled={isSubmitting} isLoading={isSubmitting}>
          Add Group
        </Button>
      </HStack>
    </form>
  );
};

export default AddCategoryGroup;
