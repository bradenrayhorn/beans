import "../styles/globals.css";
import type { AppProps } from "next/app";
import { Button, ChakraProvider, extendTheme, useColorMode } from "@chakra-ui/react";
import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

const PageCard = {
  baseStyle: (props: StyleFunctionProps) => ({
    borderRadius: 6,
    bg: mode('gray.50', 'gray.700')(props),
  }),
};

const theme = extendTheme({
  components: {
    PageCard,
  }
});

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider theme={theme}>
      <Component {...pageProps} />
      <Tog />
    </ChakraProvider>
  );
}

const Tog = () => {
  const { toggleColorMode } = useColorMode();
  return (
    <Button onClick={toggleColorMode} position="fixed" bottom="2" right="2">
      T
    </Button>
  )
};

export default MyApp;
