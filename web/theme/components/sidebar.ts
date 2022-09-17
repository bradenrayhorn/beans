import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

export default {
  baseStyle: (props: StyleFunctionProps) => ({
    bg: mode("gray.50", "gray.700")(props),
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-between",
  }),
};
