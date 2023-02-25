import { mode, StyleFunctionProps } from "@chakra-ui/theme-tools";

const TableRow = {
  baseStyle: (props: StyleFunctionProps) => ({
    borderBottom: "1px",
    borderColor: mode("gray.100", "gray.700")(props),
  }),
};

export default TableRow;
