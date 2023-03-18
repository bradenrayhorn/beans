import { alertAnatomy } from "@chakra-ui/anatomy";
import { createMultiStyleConfigHelpers } from "@chakra-ui/react";

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(alertAnatomy.keys);

const xs = definePartsStyle({
  container: {
    fontSize: "xs",
    p: 2,
    borderRadius: 4,
  },

  icon: {
    boxSize: 3,
  },
});

const Alert = defineMultiStyleConfig({
  sizes: { xs },
});

export default Alert;
