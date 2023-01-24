import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

export default {
  baseStyle: (props: StyleFunctionProps) => ({
    borderBottom: "1px",
    borderColor: mode("gray.100", "gray.700")(props),
  }),
};
