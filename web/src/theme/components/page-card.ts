import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

const PageCard = {
  baseStyle: (props: StyleFunctionProps) => ({
    borderRadius: 6,
    bg: mode("gray.50", "gray.700")(props),
    boxShadow: mode("base", "lg")(props),
  }),
};

export default PageCard;
