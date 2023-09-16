import { MsgSend } from "@hoangdv2429/quicksilverjs/dist/codegen/cosmos/bank/v1beta1/tx"
import { Coin } from "@hoangdv2429/quicksilverjs/dist/codegen/cosmos/base/v1beta1/coin"

export const staking = (zone, sender) => {
    const msgSend = MsgSend.fromJSON({
        fromAddress: sender,
        toAddress: zone.depositAddress,
        amount: new Coin()
    })
}

