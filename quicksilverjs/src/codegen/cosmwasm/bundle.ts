import * as _94 from "./wasm/v1/genesis";
import * as _95 from "./wasm/v1/ibc";
import * as _96 from "./wasm/v1/proposal";
import * as _97 from "./wasm/v1/query";
import * as _98 from "./wasm/v1/tx";
import * as _99 from "./wasm/v1/types";
import * as _299 from "./wasm/v1/tx.amino";
import * as _300 from "./wasm/v1/tx.registry";
import * as _301 from "./wasm/v1/query.lcd";
import * as _302 from "./wasm/v1/query.rpc.Query";
import * as _303 from "./wasm/v1/tx.rpc.msg";
import * as _357 from "./lcd";
import * as _358 from "./rpc.query";
import * as _359 from "./rpc.tx";
export namespace cosmwasm {
  export namespace wasm {
    export const v1 = {
      ..._94,
      ..._95,
      ..._96,
      ..._97,
      ..._98,
      ..._99,
      ..._299,
      ..._300,
      ..._301,
      ..._302,
      ..._303
    };
  }
  export const ClientFactory = {
    ..._357,
    ..._358,
    ..._359
  };
}