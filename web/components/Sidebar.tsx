import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Button,
  Divider,
  Flex,
  Heading,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  useStyleConfig,
  VStack,
} from "@chakra-ui/react";
import { routes } from "constants/routes";
import { useBudget } from "data/queries/budget";
import Link from "next/link";
import { useUser } from "./AuthProvider";

const links = [
  { name: "home", pathname: routes.budget.index },
  { name: "budget", pathname: routes.budget.budget },
  { name: "accounts", pathname: routes.budget.accounts },
  { name: "settings", pathname: routes.budget.settings },
];

const Sidebar = () => {
  const user = useUser();
  const styles = useStyleConfig("Sidebar");
  const { budget } = useBudget();

  return (
    <Flex
      __css={styles}
      p={4}
      h="100vh"
      position="sticky"
      top="0"
      w={56}
      boxShadow="md"
      flexShrink="0"
    >
      <Flex direction="column">
        <Heading size="md">beans</Heading>
        <Divider my={3} />
        <Link passHref href={routes.budget.noneSelected}>
          <Button as="a" leftIcon={<ArrowBackIcon />} size="xs">
            {budget.name}
          </Button>
        </Link>
        <VStack align="flex-start" mt={6}>
          {links.map(({ name, pathname }) => (
            <Link
              key={pathname}
              href={{
                pathname,
                query: { budget: budget.id },
              }}
              passHref
            >
              <Button
                as="a"
                size="sm"
                w="full"
                justifyContent="flex-start"
                variant="ghost"
              >
                {name}
              </Button>
            </Link>
          ))}
        </VStack>
      </Flex>
      <Flex direction="column">
        <Divider my={3} />
        <Menu>
          <MenuButton as={Button} variant="ghost" size="sm" textAlign="left">
            {user?.username}
          </MenuButton>
          <MenuList>
            <MenuItem>Log out</MenuItem>
          </MenuList>
        </Menu>
      </Flex>
    </Flex>
  );
};

export default Sidebar;
