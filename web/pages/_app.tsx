import "../styles/globals.css";
import type { AppContext, AppProps } from "next/app";
import { Button, ChakraProvider, extendTheme, useColorMode } from "@chakra-ui/react";
import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";
import { dehydrate, Hydrate, QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactElement, ReactNode, useState } from "react";
import App from "next/app";
import { queries, queryKeys } from "constants/queries";
import nookies from 'nookies';
import Head from "next/head";
import { NextPage } from "next";
import { AuthProvider } from "components/AuthProvider";

const PageCard = {
  baseStyle: (props: StyleFunctionProps) => ({
    borderRadius: 6,
    bg: mode('gray.50', 'gray.700')(props),
  }),
};

const Sidebar = {
  baseStyle: (props: StyleFunctionProps) => ({
    bg: mode('gray.50', 'gray.700')(props),
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
  }),
};

const theme = extendTheme({
  components: {
    PageCard,
    Sidebar,
  }
});

export type NextPageWithLayout = NextPage & {
  getLayout?: (page: ReactElement) => ReactNode
}

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout
}

function MyApp({ Component, pageProps }: AppPropsWithLayout) {
  const [queryClient] = useState(() => new QueryClient());

  const getLayout = Component.getLayout ?? (page => page)

  return (
    <>
      <Head>
        <title>Beans</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <QueryClientProvider client={queryClient}>
        <Hydrate state={pageProps.dehydratedState}>
          <ChakraProvider theme={theme}>
            <AuthProvider>
              {getLayout(<Component {...pageProps} />)}
              <Tog />
            </AuthProvider>
          </ChakraProvider>
        </Hydrate>
      </QueryClientProvider>
    </>
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

// MyApp.getInitialProps = async function getInitialProps(context: AppContext) {
//   const appProps = await App.getInitialProps(context);

//   const cookies = nookies.get(context.ctx);
//   const sessionID = cookies["session_id"];
//   const queryClient = new QueryClient();
//   await queryClient.prefetchQuery([queryKeys.me], () => queries.me({ sessionID }));


//   return {
//     ...appProps,
//     pageProps: {
//       ...appProps.pageProps,
//       dehydratedState: dehydrate(queryClient),
//     }
//   }
// }

export default MyApp;
