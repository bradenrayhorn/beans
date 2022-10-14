import { RowData, TValue } from "@tanstack/table-core";

declare module "@tanstack/table-core" {
  interface ColumnMeta<TData extends RowData, TValue> {
    isNumeric?: boolean;
  }
}
