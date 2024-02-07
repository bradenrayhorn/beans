import { formatAmount } from "$lib/amount";
import Fraction from "fraction.js";

export type APIAmount = string;

export class Amount {
  private fraction: Fraction;
  public display: string;
  public rawDisplay: string;

  constructor(apiAmount: APIAmount) {
    this.fraction = new Fraction(apiAmount);

    this.rawDisplay = this.fraction.valueOf().toString();

    this.display = formatAmount(this.fraction.valueOf());
  }
}
