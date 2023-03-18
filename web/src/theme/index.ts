import { extendTheme } from "@chakra-ui/react";
import Accordion from "./components/accordion";
import Alert from "./components/alert";
import ComponentSelect from "./components/component-select";
import PageCard from "./components/page-card";
import Sidebar from "./components/sidebar";
import TableHeader from "./components/table/header";
import TableRow from "./components/table/row";

const overrides = {
  components: {
    Accordion,
    Alert,
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
