import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

const Sidebar = {
  baseStyle: (props: StyleFunctionProps) => ({
    bg: mode("gray.50", "gray.700")(props),
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-between",
  }),
};

export default Sidebar;
