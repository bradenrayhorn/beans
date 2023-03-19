import { useSetMonthDate } from "@/context/MonthProvider";
import { currentDate, dateToString, parseDate } from "@/data/format/date";
import { useMonth } from "@/data/queries/month";
import { ChevronLeftIcon, ChevronRightIcon } from "@chakra-ui/icons";
import { Flex, IconButton, Skeleton } from "@chakra-ui/react";
import dayjs from "dayjs";
import MonthPicker from "./MonthPicker";

const MonthHeader = () => {
  const { month, isSuccess: isLoaded } = useMonth();
  const setMonthDate = useSetMonthDate();

  const selectMonth = (newDate: dayjs.Dayjs) => {
    setMonthDate(dateToString(newDate));
  };

  const today = parseDate(month?.date ?? currentDate());

  return (
    <Skeleton isLoaded={isLoaded}>
      <Flex alignItems="center" shrink="0">
        <IconButton
          aria-label="Previous month"
          icon={<ChevronLeftIcon />}
          variant="ghost"
          onClick={() => {
            selectMonth(today.subtract(1, "month"));
          }}
        />
        <MonthPicker month={month} selectMonth={selectMonth} />
        <IconButton
          aria-label="Next month"
          icon={<ChevronRightIcon />}
          variant="ghost"
          onClick={() => {
            selectMonth(today.add(1, "month"));
          }}
        />
      </Flex>
    </Skeleton>
  );
};

export default MonthHeader;
