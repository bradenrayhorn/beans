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
  selectedRows: { [key: string]: boolean };
  setEditID: (id: string | null) => void;
  setIsAdding: (is: boolean) => void;
  isRowSelected: (id: string) => boolean;
  setRowSelection: (id: string, selected: boolean) => void;
  getSelectedRows: () => string[];
}

export const useTransactionTableState = create<TransactionTableState>(
  (set, get) => ({
    editID: null,
    isAdding: false,
    selectedRows: {},
    setEditID: (id: string | null) => set({ editID: id, isAdding: false }),
    setIsAdding: (is: boolean) => set({ editID: null, isAdding: is }),
    isRowSelected: (id: string) => !!get().selectedRows[id],
    setRowSelection: (id: string, selected: boolean) => {
      const current = get().selectedRows;
      current[id] = selected;
      set({ selectedRows: current });
    },
    getSelectedRows: () =>
      Object.entries(get().selectedRows)
        .filter(([, v]) => v)
        .map(([k]) => k),
  })
);

export default function TransactionList() {
  const { transactions } = useTransactions();

  const isAdding = useTransactionTableState((state) => state.isAdding);
  const setIsAdding = useTransactionTableState((state) => state.setIsAdding);
  const editID = useTransactionTableState((state) => state.editID);
  const setEditID = useTransactionTableState((state) => state.setEditID);
  const getSelectedRows = useTransactionTableState(
    (state) => state.getSelectedRows
  );

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.code === "Escape") {
        setEditID(null);
      }

      if (e.key.toLowerCase() === "e") {
        const selectedRows = getSelectedRows();
        if (selectedRows.length === 1) {
          setEditID(selectedRows[0]);
        }
      }
    }

    document.addEventListener("keydown", handleKeyDown);

    return function cleanup() {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [setEditID, getSelectedRows]);

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
