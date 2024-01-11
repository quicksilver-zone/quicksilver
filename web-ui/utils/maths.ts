import BigNumber from 'bignumber.js';

export const isGreaterThanZero = (
  val: number | string | undefined,
) => {
  return new BigNumber(val || 0).gt(0);
};

export const shiftDigits = (
  num: string | number,
  places: number,
  decimalPlaces?: number,
) => {
  return new BigNumber(num)
    .shiftedBy(places)
    .decimalPlaces(decimalPlaces || 6)
    .toString();
};

export const toNumber = (
  val: string,
  decimals: number = 6,
) => {
  return new BigNumber(val)
    .decimalPlaces(decimals)
    .toNumber();
};

export const sum = (...args: string[]) => {
  return args
    .reduce(
      (prev, cur) => prev.plus(cur),
      new BigNumber(0),
    )
    .toString();
};

export function abbreviateNumber(value: number): string {
  
  if (value < 1000) {
    return Number(value.toFixed(1)).toString();
  }

  const suffixes = ["", "k", "M", "B", "T"];

  const suffixNum = Math.floor(Math.log10(value) / 3);
 
  let shortValue = value / Math.pow(1000, suffixNum);

  shortValue = Math.round(shortValue * 10) / 10;
  
  
  let newValue = shortValue % 1 === 0 ? shortValue.toString() : shortValue.toFixed(1);

  return newValue + suffixes[suffixNum];
}