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
import type { GetServerSideProps, NextPage } from "next";
import { useForm } from "react-hook-form";
import { useRouter } from "next/router";
import ky from "ky";
import { buildQueries, queries } from "constants/queries";

export const getServerSideProps: GetServerSideProps = async (context) => {
  const client = ky.extend({ prefixUrl: "http://localhost:8000" });

  try {
    await buildQueries(client).me({ cookie: context.req.headers.cookie });
    // user is logged in if get me was success
    return {
      redirect: {
        permanent: false,
        destination: "/app",
      },
    };
  } catch {}

  return { props: {} };
};

const Login: NextPage = () => {
  const toast = useToast();
  const router = useRouter();
  const {
    handleSubmit,
    register,
    formState: { isSubmitting },
  } = useForm();

  // TODO search react query + axios cancellation
  const mutation = useMutation(queries.login, {
    onSuccess: () => {
      router.push("/app");
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
