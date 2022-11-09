import {
  Box,
  Text,
  Divider,
  Flex,
  Heading,
  Spinner,
  VStack,
  HStack,
  List,
  ListItem,
} from "@chakra-ui/react";
import AddCategoryButton from "@/components/settings/AddCategoryButton";
import AddCategoryGroup from "@/components/settings/AddCategoryGroup";
import { useCategories } from "@/data/queries/category";

export default function SettingsPage() {
  const { isLoading, categoryGroups } = useCategories();

  return (
    <Box as="main" w="full">
      <Flex justify="space-between" align="center" mb={8}>
        <Flex align="center">
          <Heading size="lg">Budget Settings</Heading>
        </Flex>
      </Flex>
      <VStack>
        <VStack align="flex-start" w="full">
          <Heading size="md">Categories</Heading>
          <VStack align="flex-start" w="full" pl="4">
            {isLoading ? (
              <Spinner />
            ) : (
              <>
                <List aria-label="Categories">
                  {categoryGroups.map(({ id, name, categories }) => (
                    <VStack as={ListItem} key={id} align="flex-start">
                      <HStack>
                        <Heading size="sm">{name}</Heading>
                        <AddCategoryButton groupID={id} />
                      </HStack>
                      <VStack as={List} align="flex-start" pl="4">
                        {categories.map(({ id, name }) => (
                          <Text key={id} as={ListItem}>
                            {name}
                          </Text>
                        ))}
                      </VStack>
                    </VStack>
                  ))}
                </List>

                <Divider />
                <AddCategoryGroup />
              </>
            )}
          </VStack>
        </VStack>
      </VStack>
    </Box>
  );
}
