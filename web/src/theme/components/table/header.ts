import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

export default {
  baseStyle: (props: StyleFunctionProps) => ({
    display: "flex",
    borderBottom: "1px",
    borderColor: mode("gray.100", "gray.700")(props),
    color: mode("gray.600", "gray.400")(props),
  }),
};
