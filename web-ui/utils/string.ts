  export function formatQasset(denom: string): string {
    if (denom.substring(0, 1) == "Q" || denom.substring(0, 2) == "AQ"){
        return "q"+denom.substring(1)
      } 
      return denom
  }