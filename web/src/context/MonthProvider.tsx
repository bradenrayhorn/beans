import { createContext, PropsWithChildren, useContext, useState } from "react";

const MonthContext = createContext({
  monthDate: "",
  setMonthDate: (_monthDate: string): void => {},
});

export const useMonthDate = () => useContext(MonthContext).monthDate;

export const useSetMonthDate = () => useContext(MonthContext).setMonthDate;

type Props = {
  defaultMonthDate: string;
} & PropsWithChildren;

export default function MonthProvider({ children, defaultMonthDate }: Props) {
  const [monthDate, setMonthDate] = useState(() => defaultMonthDate);

  return (
    <MonthContext.Provider value={{ monthDate, setMonthDate }}>
      {children}
    </MonthContext.Provider>
  );
}
