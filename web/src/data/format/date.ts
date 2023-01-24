import dayjs from "dayjs";

export const formatDate = (date: string): string =>
  parseDate(date).format("MM/DD/YYYY");

export const formatBudgetMonth = (date: string): string =>
  parseDate(date).format("YYYY.MM");

export const parseDate = (date: string): dayjs.Dayjs =>
  dayjs(date, "YYYY-MM-DD", true);

export const currentDate = () => dayjs().format("YYYY-MM-DD");
