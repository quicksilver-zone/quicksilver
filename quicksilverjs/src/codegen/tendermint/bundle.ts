import * as _117 from "./abci/types";
import * as _118 from "./crypto/keys";
import * as _119 from "./crypto/proof";
import * as _120 from "./libs/bits/types";
import * as _121 from "./p2p/types";
import * as _122 from "./types/block";
import * as _123 from "./types/evidence";
import * as _124 from "./types/params";
import * as _125 from "./types/types";
import * as _126 from "./types/validator";
import * as _127 from "./version/types";
export namespace tendermint {
  export const abci = { ..._117
  };
  export const crypto = { ..._118,
    ..._119
  };
  export namespace libs {
    export const bits = { ..._120
    };
  }
  export const p2p = { ..._121
  };
  export const types = { ..._122,
    ..._123,
    ..._124,
    ..._125,
    ..._126
  };
  export const version = { ..._127
  };
}