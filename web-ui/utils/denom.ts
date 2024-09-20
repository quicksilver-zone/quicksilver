  export function formatQasset(denom: string): string {
    if (denom.substring(0, 1) == "Q" || denom.substring(0, 2) == "AQ"){
        return "q"+denom.substring(1)
      } 
      return denom
  }
  
  export function qDenomForDenom(denom: string|undefined, exponent: number|undefined): string {
    if (denom == null) { return ""}
    let prefix = ""
    switch (exponent) {
      case 3:
        prefix = "m"
        break;
      case 6:
        prefix = "u"
        break;
      case 9:
        prefix = "n"
        break;
      case 12:
        prefix = "p"
        break;
      case 15:
        prefix = "f"
        break;
      case 18:
        prefix = "a"
        break;
      case 21:
        prefix = "z"
        break;
      case 24:
        prefix = "y"
        break;
      default:
        prefix = ""
        break;
    }
    return prefix+"q"+denom
  }

  export function denomForQDenom(denom: string): string {
    return denom.substring(0, 1)+denom.substring(2)   
  }