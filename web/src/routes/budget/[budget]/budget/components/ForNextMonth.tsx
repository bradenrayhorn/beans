import { amountToFraction, formatAmount } from "@/data/format/amount";
import { useMonth } from "@/data/queries/month";
import {
  Button,
  Heading,
  Popover,
  PopoverArrow,
  PopoverCloseButton,
  PopoverContent,
  PopoverTrigger,
  Skeleton,
  Tag,
} from "@chakra-ui/react";
import { useId, useRef } from "react";
import ForNextMonthForm from "./ForNextMonthForm";

export default function ForNextMonth() {
  const { month, isLoading: isMonthLoading } = useMonth();
  const formInputRef = useRef<HTMLInputElement>(null);

  const labelID = useId();

  return (
    <Skeleton isLoaded={!isMonthLoading} w="full">
      <Popover isLazy placement="bottom-start" initialFocusRef={formInputRef}>
        {({ onClose }) => (
          <>
            <PopoverTrigger>
              <Button
                borderRadius={4}
                bg="white"
                shadow="sm"
                px={4}
                py={2}
                justifyContent="space-between"
                alignItems="center"
                variant="ghost"
                w="full"
                _hover={{ bg: "gray.50" }}
              >
                <Heading size="xs" id={labelID}>
                  For Next Month
                </Heading>

                <Tag
                  aria-labelledby={labelID}
                  colorScheme={
                    amountToFraction(month?.carryover).compare(0) > 0
                      ? "gray"
                      : "red"
                  }
                >
                  {formatAmount(month?.carryover)}
                </Tag>
              </Button>
            </PopoverTrigger>

            <PopoverContent p={4} aria-label="Edit For Next Month">
              <PopoverArrow />
              <PopoverCloseButton />

              <ForNextMonthForm
                inputRef={formInputRef}
                month={month}
                onClose={onClose}
              />
            </PopoverContent>
          </>
        )}
      </Popover>
    </Skeleton>
  );
}
