import {
  amountToFraction,
  formatAmount,
  formatFraction,
  zeroAmount,
} from "@/data/format/amount";
import { useMonth } from "@/data/queries/month";
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Flex,
  Heading,
  Skeleton,
  Tag,
  Text,
} from "@chakra-ui/react";
import { useId } from "react";

const colorScheme = (res: number): string => {
  if (res > 0) {
    return "blue";
  } else if (res < 0) {
    return "red";
  } else {
    return "green";
  }
};

export default function ToBudget() {
  const { month, isLoading: isMonthLoading } = useMonth();

  const budgetable = month?.budgetable ?? zeroAmount;

  const toBudgetId = useId();
  const breakdownId = useId();

  const breakdown = [
    {
      id: `${breakdownId}-income`,
      name: "Income",
      value: amountToFraction(month?.income ?? zeroAmount),
    },
    {
      id: `${breakdownId}-last-month`,
      name: "From last month",
      value: amountToFraction(month?.carried_over ?? zeroAmount),
    },
    {
      id: `${breakdownId}-assigned`,
      name: "Assigned this month",
      value: amountToFraction(month?.assigned ?? zeroAmount).neg(),
    },
    {
      id: `${breakdownId}-next-month`,
      name: "For next month",
      value: amountToFraction(month?.carryover ?? zeroAmount).neg(),
    },
  ];

  return (
    <Skeleton isLoaded={!isMonthLoading} w="full">
      <Accordion allowToggle reduceMotion w="full" variant="minimal">
        <AccordionItem boxShadow="sm">
          <Heading size="xs">
            <AccordionButton aria-label="To Budget">
              <Flex alignItems="center">
                <Box mr={1} fontSize="sm" as="b" role="term" id={toBudgetId}>
                  To Budget
                </Box>
                <AccordionIcon />
              </Flex>
              <Tag
                colorScheme={colorScheme(
                  amountToFraction(budgetable).compare(0)
                )}
                role="definition"
                aria-labelledby={toBudgetId}
              >
                {formatAmount(budgetable)}
              </Tag>
            </AccordionButton>
          </Heading>

          <AccordionPanel>
            {breakdown.map(({ id, name, value }) => (
              <Flex justifyContent="space-between" key={id}>
                <Text fontSize="sm" id={id} role="term">
                  {name}:
                </Text>
                <Text fontSize="sm" aria-labelledby={id} role="definition">
                  {formatFraction(value)}
                </Text>
              </Flex>
            ))}
          </AccordionPanel>
        </AccordionItem>
      </Accordion>
    </Skeleton>
  );
}
