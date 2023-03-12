import { accordionAnatomy } from "@chakra-ui/anatomy";
import { createMultiStyleConfigHelpers } from "@chakra-ui/react";

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(accordionAnatomy.keys);

const minimal = definePartsStyle({
  container: {
    borderTopWidth: "0",
    borderRadius: 4,
    bg: "white",
    _last: { borderBottomWidth: 0 },
  },

  button: {
    justifyContent: "space-between",
    borderRadius: 4,
    _expanded: { borderBottomRadius: 0 },
  },

  panel: {
    display: "flex",
    flexDir: "column",
    gap: 1,
    py: 2,
    px: 4,
    borderTopColor: "gray.100",
    borderTopWidth: 1,
  },
});

const Accordion = defineMultiStyleConfig({
  variants: { minimal },
});

export default Accordion;
