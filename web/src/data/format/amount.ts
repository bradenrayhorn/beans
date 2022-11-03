import { Amount } from "@/constants/types";
import Fraction from "fraction.js";

export const amountToFraction = (amount: Amount): Fraction => {
  return new Fraction(amount.coefficient).mul(
    new Fraction(10).pow(amount.exponent)
  );
};

export const formatAmount = (amount: Amount): string =>
  amountToFraction(amount).valueOf().toLocaleString(undefined, {
    minimumFractionDigits: 2,
    style: "currency",
    currency: "USD",
  });
