
export const broadcastTx = async (client, address, gas, msgs) => {
    let fee = {
      amount: [],
      gas: `${gas}`,
    };
    const denom = process.env.REACT_APP_DENOM
    const msg = makeSendMsg(address, recipient, amount, denom)
  
    const result = await client.signAndBroadcast(
      address,
      msgs,
      fee,
    );
    return result
}