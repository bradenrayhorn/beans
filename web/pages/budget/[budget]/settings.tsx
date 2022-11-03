import {
  Box,
  Text,
  Divider,
  Flex,
  Heading,
  Spinner,
  VStack,
  HStack,
} from "@chakra-ui/react";
import BudgetLayout from "components/layouts/BudgetLayout";
import AddCategoryButton from "components/settings/AddCategoryButton";
import AddCategoryGroup from "components/settings/AddCategoryGroup";
import { useCategories } from "data/queries/category";
import { NextPageWithLayout } from "pages/_app";

const Settings: NextPageWithLayout = () => {
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
          <VStack align="flex-start" w="full" spacing="8" pl="4">
            {isLoading ? (
              <Spinner />
            ) : (
              <>
                {categoryGroups.map(({ id, name, categories }) => (
                  <VStack key={id} align="flex-start">
                    <HStack>
                      <Heading size="sm">{name}</Heading>
                      <AddCategoryButton groupID={id} />
                    </HStack>
                    <VStack align="flex-start" pl="4">
                      {categories.map(({ id, name }) => (
                        <Text key={id}>{name}</Text>
                      ))}
                    </VStack>
                  </VStack>
                ))}

                <Divider />
                <AddCategoryGroup />
              </>
            )}
          </VStack>
        </VStack>
      </VStack>
    </Box>
  );
};

Settings.getLayout = (page) => <BudgetLayout>{page}</BudgetLayout>;

export default Settings;
