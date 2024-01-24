  export function formatQasset(denom: string): string {
    if (denom.substring(0, 1) == "Q") {
        return "q"+denom.substring(1)
      } 
      return denom
  }