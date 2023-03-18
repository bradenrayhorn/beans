import {
  Input,
  InputGroup,
  InputLeftElement,
  useMergeRefs,
} from "@chakra-ui/react";
import { forwardRef } from "react";
import { useController, UseControllerProps } from "react-hook-form";
import { NumericFormat } from "react-number-format";

const maxInput = 9999999999;
const minInput = -maxInput;

type Props = {
  name: string;
  allowNegative?: boolean;
  "aria-labelledby"?: string;
};

const CurrencyInput = forwardRef(
  (
    {
      name,
      allowNegative,
      "aria-labelledby": ariaLabelledBy,
      ...props
    }: UseControllerProps & Props,
    ref
  ) => {
    const {
      field: { onChange, onBlur, value: formValue, ref: formRef },
    } = useController({ ...props, name });

    const refs = useMergeRefs(ref, formRef);

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
          getInputRef={refs}
          customInput={Input}
          name={name}
          pl={6}
          aria-labelledby={ariaLabelledBy}
        />
      </InputGroup>
    );
  }
);

export default CurrencyInput;
