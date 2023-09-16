import { coins, coin } from "@cosmjs/amino";
const gov_1 = require("cosmjs-types/cosmos/gov/v1beta1/gov")
import axios from "axios"
import {
    calculateFee,
    GasPrice
} from '@cosmjs/stargate';
const amino_1 = require("@cosmjs/amino");
const math_1 = require("@cosmjs/math");
const utils_1 = require("@cosmjs/utils");
const signing_1 = require("cosmjs-types/cosmos/tx/signing/v1beta1/signing");
const service_1 = require("cosmjs-types/cosmos/tx/v1beta1/service");
const tx_1 = require("cosmjs-types/cosmos/tx/v1beta1/tx");
const long_1 = require("long")
const proto_signing_1 = require("@cosmjs/proto-signing");
const queryclient_1 = require("@cosmjs/stargate");
import { QueryClient } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import Timestamp from "timestamp-nano";



export const createRequestRedemptionMsg = (
    sender,
    amount,
    denom,
    destinationAddress,
    sourceChannel,
) => {
    const amt = `${amount}`
    const msgRequestRedemption = {
        fromAddress: sender,
        value: coin(amt, denom),
        destinationAddress: destinationAddress
    };
    const msg = {
        typeUrl: "/quicksilver.interchainstaking.v1.Msg/RequestRedemption",
        value: msgRequestRedemption,
    };
    return msg
}




