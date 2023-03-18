import { Amount } from "@/constants/types";
import Fraction from "fraction.js";

export const amountToFraction = (amount?: Amount): Fraction => {
  amount = amount ?? zeroAmount;
  return new Fraction(amount.coefficient).mul(
    new Fraction(10).pow(amount.exponent)
  );
};

export const zeroAmount: Amount = { exponent: 0, coefficient: 0 };

export const formatAmount = (amount?: Amount): string =>
  formatFraction(amountToFraction(amount));

export const formatFraction = (fraction: Fraction): string =>
  fraction.valueOf().toLocaleString(undefined, {
    minimumFractionDigits: 2,
    style: "currency",
    currency: "USD",
  });
