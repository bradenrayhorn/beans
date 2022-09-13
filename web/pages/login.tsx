import {
  Button,
  Container,
  FormControl,
  FormLabel,
  Heading,
  Input,
  useToast,
  VStack,
} from "@chakra-ui/react";
import { useMutation } from "@tanstack/react-query";
import PageCard from "components/PageCard";
import type { NextPage } from "next";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import { queries } from "constants/queries";
import { routes } from "constants/routes";
import { useOnLogin } from "components/AuthProvider";

const Login: NextPage = () => {
  const toast = useToast();
  const router = useRouter();
  const {
    handleSubmit,
    register,
    formState: { isSubmitting },
  } = useForm();
  const onLogin = useOnLogin();

  // TODO search react query + axios cancellation
  const mutation = useMutation(queries.login, {
    onSuccess: (user) => {
      onLogin(user);
      router.push(routes.defaultAfterAuth);
    },
  });

  const onSubmit = (event: any) =>
    handleSubmit((v: any) => mutation.mutateAsync(v))(event).catch((err) => {
      toast({ title: "Error", status: "error" });
      console.error(err);
    });

  return (
    <main>
      <Container mt={4}>
        <Heading>Log in to beans</Heading>
        <PageCard>
          <form onSubmit={onSubmit}>
            <VStack p={8} mt={4} spacing={6}>
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
      </Container>
    </main>
  );
};

export default Login;
