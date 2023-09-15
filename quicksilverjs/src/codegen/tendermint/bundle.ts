import * as _211 from "./abci/types";
import * as _212 from "./crypto/keys";
import * as _213 from "./crypto/proof";
import * as _214 from "./libs/bits/types";
import * as _215 from "./p2p/types";
import * as _216 from "./types/block";
import * as _217 from "./types/evidence";
import * as _218 from "./types/params";
import * as _219 from "./types/types";
import * as _220 from "./types/validator";
import * as _221 from "./version/types";
export namespace tendermint {
  export const abci = {
    ..._211
  };
  export const crypto = {
    ..._212,
    ..._213
  };
  export namespace libs {
    export const bits = {
      ..._214
    };
  }
  export const p2p = {
    ..._215
  };
  export const types = {
    ..._216,
    ..._217,
    ..._218,
    ..._219,
    ..._220
  };
  export const version = {
    ..._221
  };
}