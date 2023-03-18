import { Button, ButtonProps } from "@chakra-ui/react";
import { useFormState } from "react-hook-form";

export default function FormButton({ children, ...rest }: ButtonProps) {
  const { isSubmitting } = useFormState();

  return (
    <Button isDisabled={isSubmitting} isLoading={isSubmitting} {...rest}>
      {children}
    </Button>
  );
}
