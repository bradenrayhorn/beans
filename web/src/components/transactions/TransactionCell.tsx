import { Transaction } from "@/constants/types";
import { Flex, FlexProps } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Column } from "./TransactionColumns";

export default function TransactionCell({
  column,
  isEditing,
  children,
  ...props
}: {
  column: Column<Transaction>;
  isEditing?: boolean;
} & PropsWithChildren &
  FlexProps) {
  return (
    <Flex
      key={column.id}
      width={column.width}
      shrink={!!column.width ? 0 : undefined}
      flex={!column.width ? 1 : undefined}
      justifyContent={column.isNumeric ? "flex-end" : ""}
      pr={isEditing ? 1 : 0}
      _last={{ pr: 0 }}
      {...props}
    >
      {children}
    </Flex>
  );
}
