import { MonthCategory } from "@/constants/types";
import { formatAmount } from "@/data/format/amount";
import {
  Button,
  Popover,
  PopoverArrow,
  PopoverCloseButton,
  PopoverContent,
  PopoverTrigger,
  useDisclosure,
} from "@chakra-ui/react";
import { useRef } from "react";
import EditAssignedForm from "./EditAssignedForm";

export default function AssignedCell({
  monthCategory,
}: {
  monthCategory?: MonthCategory;
}) {
  const { onOpen, onClose, isOpen } = useDisclosure();
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <Popover
      isOpen={isOpen}
      onOpen={onOpen}
      onClose={onClose}
      placement="right"
      initialFocusRef={inputRef}
      isLazy={true}
    >
      <PopoverTrigger>
        <Button
          size="sm"
          variant="link"
          color="chakra-body-text"
          fontWeight="normal"
        >
          {formatAmount(monthCategory?.assigned)}
        </Button>
      </PopoverTrigger>

      <PopoverContent p={5} aria-label="Edit assigned">
        <PopoverArrow />
        <PopoverCloseButton />

        <EditAssignedForm
          onClose={onClose}
          inputRef={inputRef}
          monthCategory={monthCategory}
        />
      </PopoverContent>
    </Popover>
  );
}
