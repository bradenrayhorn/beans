import { writable } from "svelte/store";

export const selectedRows = writable<{ [transactionID: string]: boolean }>({});
