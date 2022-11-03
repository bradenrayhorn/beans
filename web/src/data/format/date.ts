import dayjs from "dayjs";

export const formatDate = (date: string): string =>
  dayjs(date, "YYYY-MM-DD", true).format("MM/DD/YYYY");
