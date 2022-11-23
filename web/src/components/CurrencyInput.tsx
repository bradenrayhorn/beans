import { Input, InputGroup, InputLeftElement } from "@chakra-ui/react";
import { useController, UseControllerProps } from "react-hook-form";
import { NumericFormat } from "react-number-format";

const maxInput = 9999999999;
const minInput = -maxInput;

type Props = {
  allowNegative?: boolean;
};

const CurrencyInput = ({
  allowNegative,
  ...props
}: UseControllerProps & Props) => {
  const {
    field: { onChange, onBlur, value: formValue, ref },
  } = useController(props);

  return (
    <InputGroup>
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
        name={props.name}
        pl={10}
      />
    </InputGroup>
  );
};

export default CurrencyInput;
