import { Button, Container, FormControl, FormLabel, Heading, Input, useToast, VStack } from '@chakra-ui/react'
import { useMutation } from '@tanstack/react-query'
import PageCard from 'components/PageCard'
import type { NextPage } from 'next'
import Head from 'next/head'
import { queries } from 'constants/queries'
import { useForm } from 'react-hook-form'

const Login: NextPage = () => {
  const toast = useToast();
  const { handleSubmit, register, formState: { isSubmitting } } = useForm();

  // TODO search react query + axios cancellation
  const mutation = useMutation(queries.login, {
    onSuccess: () => {
      console.log('login successful')
    },
  });

  const onSubmit = (event: any) =>
    handleSubmit((v) => mutation.mutateAsync(v))(event).catch(err => {
      toast({ title: 'Error', status: 'error' });
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
                <Input {...register('username')} />
              </FormControl>
              <FormControl>
                <FormLabel>Password</FormLabel>
                <Input {...register('password')} type="password" />
              </FormControl>
              <Button type="submit" w="full" isLoading={isSubmitting} isDisabled={isSubmitting}>
                Log in
              </Button>
            </VStack>
          </form>
        </PageCard>
      </Container>
    </main>
  )
}

export default Login
