import { useMonthID, useSetMonthID } from "@/context/MonthProvider";
import { currentDate, dateToString, parseDate } from "@/data/format/date";
import {
  useCreateMonth,
  useIsMonthLoading,
  useMonth,
} from "@/data/queries/month";
import { ChevronLeftIcon, ChevronRightIcon } from "@chakra-ui/icons";
import { Flex, IconButton, Skeleton, useToast } from "@chakra-ui/react";
import dayjs from "dayjs";
import MonthPicker from "./MonthPicker";

const MonthHeader = () => {
  const toast = useToast();

  const monthID = useMonthID();
  const { month } = useMonth({ monthID });

  const isMonthLoading = useIsMonthLoading();
  const setMonthID = useSetMonthID();
  const { mutate } = useCreateMonth({
    onSuccess: ({ data: { month_id } }) => {
      setMonthID(month_id);
    },
    onError: () => {
      toast({
        title: "Failed to switch month.",
        status: "error",
      });
    },
  });

  const selectMonth = (newDate: dayjs.Dayjs) => {
    mutate({ date: dateToString(newDate) });
  };

  const today = parseDate(month?.date ?? currentDate());

  return (
    <Skeleton isLoaded={!isMonthLoading}>
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
