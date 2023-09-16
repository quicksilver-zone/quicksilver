export function shortenAddress(address) {
    if (address) {
        const prefix = address.slice(0, 8);
        const suffix = address.slice(-4);
      
        return `${prefix}...${suffix}`;
    }
    return ''
}