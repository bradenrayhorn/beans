import { ColorModeScript } from "@chakra-ui/react";
import { Html, Head, Main, NextScript } from "next/document";

// TODO fix color mode flash https://chakra-ui.com/docs/styled-system/color-mode#add-colormodemanager-optional-for-ssr

export default function Document() {
  return (
    <Html>
      <Head />
      <body>
        <ColorModeScript />
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
