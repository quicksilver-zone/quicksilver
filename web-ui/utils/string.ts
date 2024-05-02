  export function formatQasset(denom: string): string {
    if (denom.startsWith("Q") || denom.startsWith("AQ")) {
        return "q" + denom.substring(1);
    } 
    return denom;
        return "q" + denom.substring(1)
      } 
      return denom
  }