export function formatAmount(amount: number): string {
  return amount.toLocaleString(undefined, {
    minimumFractionDigits: 2,
    style: "currency",
    currency: "USD",
  });
}
