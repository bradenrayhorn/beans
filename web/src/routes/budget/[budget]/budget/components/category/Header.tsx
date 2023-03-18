import { Flex, Heading } from "@chakra-ui/react";

type Column = {
  name: string;
  isNumeric: boolean;
};

export default function Header({ columns }: { columns: Column[] }) {
  return (
    <Flex w="full" bg="gray.100" px="4" py={2} role="row">
      {columns.map((column) => (
        <Heading
          key={column.name}
          size="xs"
          flex={1}
          textAlign={column.isNumeric ? "right" : "left"}
          role="columnheader"
        >
          {column.name}
        </Heading>
      ))}
    </Flex>
  );
}
