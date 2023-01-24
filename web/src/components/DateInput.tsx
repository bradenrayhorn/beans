import { parseDate } from "@/data/format/date";
import { Input, InputProps } from "@chakra-ui/react";
import { useState } from "react";
import { useController } from "react-hook-form";

type Props = {
  name: string;
};

export default function DateInput({
  name,
  ...inputProps
}: Props & Omit<InputProps, "onChange" | "onBlur" | "value" | "ref" | "name">) {
  const {
    field: {
      onChange: formOnChange,
      onBlur: formOnBlur,
      value: formValue,
      ref,
    },
  } = useController({ name });

  const [value, setValue] = useState<string>(formValue);

  return (
    <Input
      {...inputProps}
      onChange={(e) => {
        setValue(e.target.value);
      }}
      name={name}
      onBlur={() => {
        const date = parseDate(value);
        if (date.isValid()) {
          const newValue = date.format("YYYY-MM-DD");

          setValue(newValue);
          formOnChange(newValue);
        } else {
          setValue(formValue);
        }
        formOnBlur();
      }}
      value={value}
      ref={ref}
    />
  );
}
