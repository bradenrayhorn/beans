import { useTransactions } from "@/data/queries/transaction";
import { Flex, Heading, useStyleConfig } from "@chakra-ui/react";
import { useEffect } from "react";
import { create } from "zustand";
import AddTransactionProvider from "./AddTransactionProvider";
import EditTransactionProvider from "./EditTransactionProvider";
import TransactionCell from "./TransactionCell";
import { columns } from "./TransactionColumns";
import TransactionRow from "./TransactionRow";

interface TransactionTableState {
  editID: string | null;
  isAdding: boolean;
  setEditID: (id: string | null) => void;
  setIsAdding: (is: boolean) => void;
}

export const useTransactionTableState = create<TransactionTableState>(
  (set) => ({
    editID: null,
    isAdding: false,
    setEditID: (id: string | null) => set({ editID: id, isAdding: false }),
    setIsAdding: (is: boolean) => set({ editID: null, isAdding: is }),
  })
);

export default function TransactionList() {
  const { transactions } = useTransactions();

  const {
    isAdding,
    setIsAdding,
    editID,
    setEditID,
  } = useTransactionTableState();

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.code === "Escape") {
        setEditID(null);
      }
    }

    document.addEventListener("keydown", handleKeyDown);

    return function cleanup() {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  const headerStyles = useStyleConfig("TableHeader");

  return (
    <Flex flexDir="column">
      <Flex direction="column" role="table" aria-rowcount={transactions.length}>
        <Flex pb={2} __css={headerStyles} role="row">
          {columns.map((column) => (
            <TransactionCell
              key={column.id}
              column={column}
              role="columnheader"
            >
              <Heading id={column.headerID} size="xs" textTransform="uppercase">
                {column.header}
              </Heading>
            </TransactionCell>
          ))}
        </Flex>

        <Flex direction="column">
          {isAdding && (
            <AddTransactionProvider onSuccess={() => setIsAdding(false)}>
              <TransactionRow
                isEditing={true}
                onEditBegin={() => {}}
                onEditEnd={() => {
                  setIsAdding(false);
                }}
              />
            </AddTransactionProvider>
          )}

          {transactions.map((data) => (
            <EditTransactionProvider
              key={data.id}
              transaction={data}
              isEditing={editID === data.id}
              onSuccess={() => {
                setEditID(null);
              }}
            >
              <TransactionRow
                transaction={data}
                isEditing={editID === data.id}
                onEditBegin={() => {
                  setEditID(data.id);
                }}
                onEditEnd={() => {
                  setEditID(null);
                }}
              />
            </EditTransactionProvider>
          ))}
        </Flex>
      </Flex>
    </Flex>
  );
}
