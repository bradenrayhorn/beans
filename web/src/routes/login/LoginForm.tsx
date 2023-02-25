import { getHTTPErrorResponseMessage, queries } from "@/constants/queries";
import { routes } from "@/constants/routes";
import { useOnLogin } from "@/context/AuthContext";
import PageCard from "@/components/PageCard";
import {
  Alert,
  AlertIcon,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { HTTPError } from "ky";

interface FormData {
  username: string;
  password: string;
}

export default function LoginForm() {
  const {
    handleSubmit,
    register,
    formState: { isSubmitting },
  } = useForm<FormData>();
  const onLogin = useOnLogin();
  const navigate = useNavigate();

  const mutation = useMutation(queries.login, {
    onSuccess: (user) => {
      onLogin(user);
      navigate(routes.defaultAfterAuth);
    },
  });

  const onSubmit = (event: any) =>
    handleSubmit((v) =>
      mutation.mutateAsync(v).catch((err) => {
        if (!(err instanceof HTTPError)) {
          console.error(err);
        }
      })
    )(event);

  return (
    <>
      {mutation.isError && (
        <Alert status="error">
          <AlertIcon />
          {getHTTPErrorResponseMessage(mutation.error)}
        </Alert>
      )}
      <PageCard w="full">
        <form onSubmit={onSubmit}>
          <VStack p={8} spacing={6}>
            <FormControl>
              <FormLabel>Username</FormLabel>
              <Input {...register("username")} />
            </FormControl>
            <FormControl>
              <FormLabel>Password</FormLabel>
              <Input {...register("password")} type="password" />
            </FormControl>
            <Button
              type="submit"
              w="full"
              isLoading={isSubmitting}
              isDisabled={isSubmitting}
            >
              Log in
            </Button>
          </VStack>
        </form>
      </PageCard>
    </>
  );
}
