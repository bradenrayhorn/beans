import { Button, Divider, Flex, Heading, Menu, MenuButton, MenuItem, MenuList, useStyleConfig } from "@chakra-ui/react";
import { useUser } from "./AuthProvider";

const Sidebar = () => {
  const user = useUser();
  const styles = useStyleConfig("Sidebar");

  return (
    <Flex __css={styles} p={4} h="full" w={56} boxShadow="md">
      <Flex direction="column">
        <Heading size="md">beans</Heading>
        <Divider my={3} />
      </Flex>
      <Flex direction="column">
        <Divider my={3} />
        <Menu>
          <MenuButton as={Button} variant="ghost" textAlign="left">{user?.username}</MenuButton>
          <MenuList>
            <MenuItem>Log out</MenuItem>
          </MenuList>
        </Menu>
      </Flex>
    </Flex>
  )
};

export default Sidebar;

