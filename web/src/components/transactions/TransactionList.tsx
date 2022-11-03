import {
  Flex,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Transaction } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import { formatDate } from "@/data/format/date";
import { useTransactions } from "@/data/queries/transaction";

const columnHelper = createColumnHelper<Transaction>();
const columns = [
  columnHelper.accessor("date", {
    header: "Date",
    cell: (info) => formatDate(info.getValue()),
  }),
  columnHelper.accessor("account.name", { header: "Account" }),
  columnHelper.accessor("notes", { header: "Notes" }),
  columnHelper.accessor("amount", {
    header: "Amount",
    cell: (info) => formatAmount(info.getValue()),
    meta: { isNumeric: true },
  }),
];

const TransactionList = () => {
  const { transactions } = useTransactions();

  const table = useReactTable({
    data: transactions,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Flex flexDir="column">
      <TableContainer>
        <Table>
          <Thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <Tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <Th
                    key={header.id}
                    isNumeric={header.column.columnDef.meta?.isNumeric ?? false}
                  >
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                  </Th>
                ))}
              </Tr>
            ))}
          </Thead>
          <Tbody>
            {table.getRowModel().rows.map((row) => (
              <Tr key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <Td
                    key={cell.id}
                    isNumeric={cell.column.columnDef.meta?.isNumeric ?? false}
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </Td>
                ))}
              </Tr>
            ))}
          </Tbody>
        </Table>
      </TableContainer>
    </Flex>
  );
};

export default TransactionList;
