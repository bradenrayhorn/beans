import { Container, Heading, VStack } from "@chakra-ui/react";
import LoginForm from "./LoginForm";

export default function LoginPage() {
  return (
    <main>
      <Container mt={4}>
        <VStack alignItems="flex-start" spacing={6}>
          <Heading>Log in to beans</Heading>
          <LoginForm />
        </VStack>
      </Container>
    </main>
  );
}
