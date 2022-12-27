import { ComponentMultiStyleConfig } from "@chakra-ui/theme";
import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

export default {
  parts: ["wrapper", "item"],
  baseStyle: (props: StyleFunctionProps) => ({
    wrapper: {
      w: "full",
      overflowX: "hidden",
      overflowY: "auto",
      maxHeight: 48,
      boxShadow: mode("sm", "dark-lg")(props),
    },
    item: {
      bg: "none",
      rounded: "none",
      display: "flex",
      justifyContent: "flex-start",
      p: 2,
      _checked: {
        bg: mode("gray.200", "whiteAlpha.200")(props),
      },
      _selected: {
        bg: mode("gray.100", "whiteAlpha.100")(props),
      },
    },
  }),
} as ComponentMultiStyleConfig;
