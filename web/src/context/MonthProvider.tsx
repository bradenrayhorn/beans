import { createContext, PropsWithChildren, useContext, useState } from "react";

const MonthContext = createContext({
  monthID: "",
  setMonthID: (_monthID: string): void => {},
});

export const useMonthID = () => useContext(MonthContext).monthID;

export const useSetMonthID = () => useContext(MonthContext).setMonthID;

type Props = {
  defaultMonthID: string;
} & PropsWithChildren;

export default function MonthProvider({ children, defaultMonthID }: Props) {
  const [monthID, setMonthID] = useState(defaultMonthID);

  return (
    <MonthContext.Provider value={{ monthID, setMonthID }}>
      {children}
    </MonthContext.Provider>
  );
}
