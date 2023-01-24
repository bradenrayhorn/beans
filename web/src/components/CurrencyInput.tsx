import { Input, InputGroup, InputLeftElement } from "@chakra-ui/react";
import { useController, UseControllerProps } from "react-hook-form";
import { NumericFormat } from "react-number-format";

const maxInput = 9999999999;
const minInput = -maxInput;

type Props = {
  name: string;
  allowNegative?: boolean;
  "aria-labelledby"?: string;
};

const CurrencyInput = ({
  name,
  allowNegative,
  "aria-labelledby": ariaLabelledBy,
  ...props
}: UseControllerProps & Props) => {
  const {
    field: { onChange, onBlur, value: formValue, ref },
  } = useController({ ...props, name });

  return (
    <InputGroup size="sm">
      <InputLeftElement pointerEvents="none">$</InputLeftElement>
      <NumericFormat
        decimalScale={2}
        thousandSeparator=","
        allowNegative={allowNegative}
        isAllowed={({ value }) => +value < maxInput && +value > minInput}
        value={formValue}
        onChange={onChange}
        onBlur={onBlur}
        getInputRef={ref}
        customInput={Input}
        name={name}
        pl={6}
        aria-labelledby={ariaLabelledBy}
      />
    </InputGroup>
  );
};

export default CurrencyInput;
