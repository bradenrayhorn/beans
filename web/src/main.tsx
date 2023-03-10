import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { ChakraProvider } from "@chakra-ui/react";
import "./index.css";
import theme from "./theme";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import customParseFormat from "dayjs/plugin/customParseFormat";
import dayjs from "dayjs";

dayjs.extend(customParseFormat);

const client = new QueryClient();

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <ChakraProvider theme={theme}>
      <QueryClientProvider client={client}>
        <App />
      </QueryClientProvider>
    </ChakraProvider>
  </React.StrictMode>
);
