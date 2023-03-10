import { Month } from "@/constants/types";
import {
  currentDate,
  formatBudgetMonth,
  formatDateAsYear,
  parseDate,
} from "@/data/format/date";
import { ChevronLeftIcon, ChevronRightIcon } from "@chakra-ui/icons";
import {
  Button,
  Flex,
  FocusLock,
  Heading,
  IconButton,
  Popover,
  PopoverContent,
  PopoverTrigger,
  SimpleGrid,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { useState } from "react";

export default function MonthPicker({
  month,
  selectMonth,
}: {
  month?: Month;
  selectMonth: (date: dayjs.Dayjs) => void;
}) {
  return (
    <Popover isLazy={true} autoFocus={true}>
      {({ onClose }) => (
        <>
          <PopoverTrigger>
            <Button variant="ghost">
              <Heading>
                {formatBudgetMonth(month?.date ?? currentDate())}
              </Heading>
            </Button>
          </PopoverTrigger>

          <PopoverContent>
            {month && (
              <FocusLock>
                <MonthPickerContent
                  month={month}
                  selectMonth={(date) => {
                    onClose();
                    selectMonth(date);
                  }}
                />
              </FocusLock>
            )}
          </PopoverContent>
        </>
      )}
    </Popover>
  );
}

function MonthPickerContent({
  month,
  selectMonth,
}: {
  month: Month;
  selectMonth: (date: dayjs.Dayjs) => void;
}) {
  const [year, setYear] = useState(() => parseDate(month.date));

  return (
    <>
      <Flex alignItems="center" justifyContent="center" mt={4} gap={2}>
        <IconButton
          aria-label="Previous year"
          icon={<ChevronLeftIcon />}
          size="sm"
          variant="ghost"
          onClick={() => {
            setYear((currentYear) => currentYear.add(-1, "year"));
          }}
        />
        <Heading size="md">{formatDateAsYear(year)}</Heading>
        <IconButton
          aria-label="Next year"
          icon={<ChevronRightIcon />}
          size="sm"
          variant="ghost"
          onClick={() => {
            setYear((currentYear) => currentYear.add(1, "year"));
          }}
        />
      </Flex>

      <SimpleGrid columns={3} spacing={1} m={4}>
        {[...Array(12)].map((_, x) => {
          const m = dayjs(`${x + 1} ${formatDateAsYear(year)}`, "M YYYY");
          const isCurrent = m.isSame(month.date, "day");

          return (
            <Button
              key={x}
              aria-current={isCurrent}
              variant={isCurrent ? "solid" : "ghost"}
              colorScheme={isCurrent ? "blue" : undefined}
              onClick={() => {
                selectMonth(m);
              }}
            >
              {m.format("MMM")}
            </Button>
          );
        })}
      </SimpleGrid>
    </>
  );
}
