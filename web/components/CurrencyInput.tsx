import { Input, InputGroup, InputLeftElement } from "@chakra-ui/react";
import { useState } from "react";
import { useController } from "react-hook-form";

const validNumberRegex = /^-?[0-9]+(\.[0-9]{1,2})?$/;

type Props = {
  name: string;
};

const CurrencyInput = ({ name }: Props) => {
  const {
    field: { onChange, onBlur, value: formValue, ref },
  } = useController({ name });

  const [value, setValue] = useState(
    formValue?.toLocaleString(undefined, {
      minimumFractionDigits: 2,
    }) ?? ""
  );

  return (
    <InputGroup>
      <InputLeftElement pointerEvents="none" children="$" />
      <Input
        name={name}
        ref={ref}
        value={value}
        onFocus={() => {
          if (value) {
            setValue(`${+value.replace(/,/g, "")}`);
          }
        }}
        onChange={(e) => {
          setValue(e.target.value);
        }}
        onBlur={() => {
          if (value.trim() === "") {
            setValue("");
            onChange(undefined);
          } else if (
            !validNumberRegex.test(value) ||
            !Number.isSafeInteger(+value)
          ) {
            setValue(
              formValue?.toLocaleString(undefined, {
                minimumFractionDigits: 2,
              }) ?? ""
            );
          } else {
            onChange(+value);
            setValue(
              (+value).toLocaleString(undefined, { minimumFractionDigits: 2 })
            );
          }

          onBlur();
        }}
      />
    </InputGroup>
  );
};

export default CurrencyInput;
