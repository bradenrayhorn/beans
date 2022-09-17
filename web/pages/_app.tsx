import { MoonIcon, SunIcon } from "@chakra-ui/icons";
import { ChakraProvider, IconButton, useColorMode } from "@chakra-ui/react";
import {
  Hydrate,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { AuthProvider } from "components/AuthProvider";
import { buildQueries } from "constants/queries";
import {
  forceUnproctedRoutes,
  routes,
  unprotectedRoutes,
} from "constants/routes";
import { User } from "constants/types";
import ky from "ky";
import { NextPage } from "next";
import type { AppContext, AppProps } from "next/app";
import App from "next/app";
import Head from "next/head";
import { ReactElement, ReactNode, useState } from "react";
import theme from "theme";
import "../styles/globals.css";

export type NextPageWithLayout = NextPage & {
  getLayout?: (page: ReactElement) => ReactNode;
};

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout;
  initialUser?: User;
};

function MyApp({ Component, pageProps, initialUser }: AppPropsWithLayout) {
  const [queryClient] = useState(() => new QueryClient());

  const getLayout = Component.getLayout ?? ((page) => page);

  return (
    <>
      <Head>
        <title>Beans</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <QueryClientProvider client={queryClient}>
        <Hydrate state={pageProps.dehydratedState}>
          <ChakraProvider theme={theme}>
            <AuthProvider initialUser={initialUser}>
              {getLayout(<Component {...pageProps} />)}
              <ColorModeToggle />
            </AuthProvider>
          </ChakraProvider>
        </Hydrate>
      </QueryClientProvider>
    </>
  );
}

const ColorModeToggle = () => {
  const { colorMode, toggleColorMode } = useColorMode();
  return (
    <IconButton
      onClick={toggleColorMode}
      position="fixed"
      bottom="2"
      right="2"
      aria-label="Toggle color mode"
      variant="ghost"
      icon={colorMode === "dark" ? <SunIcon /> : <MoonIcon />}
    />
  );
};

MyApp.getInitialProps = async function getInitialProps(context: AppContext) {
  const appProps = await App.getInitialProps(context);

  const req = context.ctx.req;
  if (!req || (req.url && req.url.startsWith("/_next/data"))) {
    return appProps;
  }

  const client = ky.extend({ prefixUrl: "http://localhost:8000" });
  let user = undefined;
  try {
    user = await buildQueries(client).me({
      cookie: context.ctx.req?.headers?.cookie,
    });
  } catch {}

  if (
    !user &&
    !unprotectedRoutes.some((route) => context.router.pathname.match(route))
  ) {
    context.ctx.res?.writeHead?.(302, { Location: routes.login });
    context.ctx.res?.end?.();
    return appProps;
  }

  if (
    user &&
    forceUnproctedRoutes.some((route) => context.router.pathname.match(route))
  ) {
    context.ctx.res?.writeHead?.(302, { Location: routes.defaultAfterAuth });
    context.ctx.res?.end?.();
    return appProps;
  }

  return {
    ...appProps,
    initialUser: user,
  };
};

export default MyApp;
