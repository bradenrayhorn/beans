import { Month } from "@/constants/types";
import { useSetMonthID } from "@/context/MonthProvider";
import { formatBudgetMonth, parseDate } from "@/data/format/date";
import { useCreateMonth } from "@/data/queries/month";
import { ChevronLeftIcon, ChevronRightIcon } from "@chakra-ui/icons";
import { Flex, Heading, IconButton, useToast } from "@chakra-ui/react";

const MonthHeader = ({ month }: { month: Month }) => {
  const toast = useToast();

  const setMonthID = useSetMonthID();
  const { isLoading, mutate } = useCreateMonth({
    onSuccess: ({ data: { month_id } }) => {
      setMonthID(month_id);
    },
    onError: () => {
      toast({
        title: "Failed to create month.",
        status: "error",
      });
    },
  });

  return (
    <Flex alignItems="center" shrink="0" gap={2}>
      <IconButton
        aria-label="Previous month"
        icon={<ChevronLeftIcon />}
        variant="ghost"
        disabled={isLoading}
        isLoading={isLoading}
        onClick={() => {
          mutate({
            date: parseDate(month.date)
              .subtract(1, "month")
              .format("YYYY-MM-DD"),
          });
        }}
      />
      <Heading>{formatBudgetMonth(month.date)}</Heading>
      <IconButton
        aria-label="Next month"
        icon={<ChevronRightIcon />}
        variant="ghost"
        disabled={isLoading}
        isLoading={isLoading}
        onClick={() => {
          mutate({
            date: parseDate(month.date).add(1, "month").format("YYYY-MM-DD"),
          });
        }}
      />
    </Flex>
  );
};

export default MonthHeader;
