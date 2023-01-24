import { Transaction } from "@/constants/types";
import { Button, Flex, useStyleConfig } from "@chakra-ui/react";
import { columns } from "./TransactionColumns";
import TransactionCell from "./TransactionCell";
import { useFormState } from "react-hook-form";

export default function TransactionRow({
  transaction,
  isEditing,
  onEditBegin,
  onEditEnd,
}: {
  transaction?: Transaction;
  isEditing: boolean;
  onEditBegin: () => void;
  onEditEnd: () => void;
}) {
  const styles = useStyleConfig("TableRow");
  const { isSubmitting } = useFormState();

  return (
    <Flex flexDir="column" w="full" pb={1} mt={1} __css={styles} role="row">
      <Flex>
        {columns.map((column) => (
          <TransactionCell
            key={column.id}
            column={column}
            isEditing={isEditing}
            onClick={() => onEditBegin()}
            role="cell"
          >
            {isEditing && !!column.input && column.input()}
            {!isEditing && transaction && column.cell(transaction)}
          </TransactionCell>
        ))}
      </Flex>

      {isEditing && (
        <Flex justifyContent="flex-end" py={2} gap={2}>
          <Button
            size="xs"
            variant="outline"
            onClick={() => onEditEnd()}
            disabled={isSubmitting}
          >
            Cancel
          </Button>

          <Button
            size="xs"
            colorScheme="green"
            type="submit"
            disabled={isSubmitting}
            isLoading={isSubmitting}
          >
            Save
          </Button>
        </Flex>
      )}
    </Flex>
  );
}
