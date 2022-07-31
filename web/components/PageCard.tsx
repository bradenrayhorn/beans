import { Box, useStyleConfig } from "@chakra-ui/react";

const PageCard = (props: any) => {
  const styles = useStyleConfig('PageCard');

  return <Box __css={styles} {...props} />
};

export default PageCard;
