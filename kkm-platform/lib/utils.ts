export function formatCurrency(amount: number): string {
  return `${amount.toFixed(2)} ₽`;
}

export function calculateTax(amount: number, rate: number): number {
  return amount * rate;
}
