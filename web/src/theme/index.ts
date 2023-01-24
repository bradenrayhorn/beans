import { extendTheme } from "@chakra-ui/react";
import PageCard from "./components/page-card";
import ComponentSelect from "./components/component-select";
import Sidebar from "./components/sidebar";
import TableRow from "./components/table/row";
import TableHeader from "./components/table/header";

const overrides = {
  components: {
    ComponentSelect,
    PageCard,
    Sidebar,
    TableHeader,
    TableRow,
  },
  semanticTokens: {
    colors: {
      errorText: {
        default: "red.500",
        _dark: "red.300",
      },
    },
  },
};

export default extendTheme(overrides);
