import * as _1 from "./auth/v1beta1/auth";
import * as _2 from "./auth/v1beta1/genesis";
import * as _3 from "./auth/v1beta1/query";
import * as _4 from "./authz/v1beta1/authz";
import * as _5 from "./authz/v1beta1/event";
import * as _6 from "./authz/v1beta1/genesis";
import * as _7 from "./authz/v1beta1/query";
import * as _8 from "./authz/v1beta1/tx";
import * as _9 from "./bank/v1beta1/authz";
import * as _10 from "./bank/v1beta1/bank";
import * as _11 from "./bank/v1beta1/genesis";
import * as _12 from "./bank/v1beta1/query";
import * as _13 from "./bank/v1beta1/tx";
import * as _14 from "./base/abci/v1beta1/abci";
import * as _15 from "./base/kv/v1beta1/kv";
import * as _16 from "./base/query/v1beta1/pagination";
import * as _17 from "./base/reflection/v1beta1/reflection";
import * as _18 from "./base/reflection/v2alpha1/reflection";
import * as _19 from "./base/snapshots/v1beta1/snapshot";
import * as _20 from "./base/store/v1beta1/commit_info";
import * as _21 from "./base/store/v1beta1/listening";
import * as _22 from "./base/store/v1beta1/snapshot";
import * as _23 from "./base/tendermint/v1beta1/query";
import * as _24 from "./base/v1beta1/coin";
import * as _25 from "./capability/v1beta1/capability";
import * as _26 from "./capability/v1beta1/genesis";
import * as _27 from "./crisis/v1beta1/genesis";
import * as _28 from "./crisis/v1beta1/tx";
import * as _29 from "./crypto/ed25519/keys";
import * as _30 from "./crypto/multisig/keys";
import * as _31 from "./crypto/secp256k1/keys";
import * as _32 from "./crypto/secp256r1/keys";
import * as _33 from "./distribution/v1beta1/distribution";
import * as _34 from "./distribution/v1beta1/genesis";
import * as _35 from "./distribution/v1beta1/query";
import * as _36 from "./distribution/v1beta1/tx";
import * as _37 from "./evidence/v1beta1/evidence";
import * as _38 from "./evidence/v1beta1/genesis";
import * as _39 from "./evidence/v1beta1/query";
import * as _40 from "./evidence/v1beta1/tx";
import * as _41 from "./feegrant/v1beta1/feegrant";
import * as _42 from "./feegrant/v1beta1/genesis";
import * as _43 from "./feegrant/v1beta1/query";
import * as _44 from "./feegrant/v1beta1/tx";
import * as _45 from "./genutil/v1beta1/genesis";
import * as _46 from "./gov/v1beta1/genesis";
import * as _47 from "./gov/v1beta1/gov";
import * as _48 from "./gov/v1beta1/query";
import * as _49 from "./gov/v1beta1/tx";
import * as _50 from "./mint/v1beta1/genesis";
import * as _51 from "./mint/v1beta1/mint";
import * as _52 from "./mint/v1beta1/query";
import * as _53 from "./params/v1beta1/params";
import * as _54 from "./params/v1beta1/query";
import * as _55 from "./slashing/v1beta1/genesis";
import * as _56 from "./slashing/v1beta1/query";
import * as _57 from "./slashing/v1beta1/slashing";
import * as _58 from "./slashing/v1beta1/tx";
import * as _59 from "./staking/v1beta1/authz";
import * as _60 from "./staking/v1beta1/genesis";
import * as _61 from "./staking/v1beta1/query";
import * as _62 from "./staking/v1beta1/staking";
import * as _63 from "./staking/v1beta1/tx";
import * as _64 from "./tx/signing/v1beta1/signing";
import * as _65 from "./tx/v1beta1/service";
import * as _66 from "./tx/v1beta1/tx";
import * as _67 from "./upgrade/v1beta1/query";
import * as _68 from "./upgrade/v1beta1/upgrade";
import * as _69 from "./vesting/v1beta1/tx";
import * as _70 from "./vesting/v1beta1/vesting";
import * as _128 from "./authz/v1beta1/tx.amino";
import * as _129 from "./bank/v1beta1/tx.amino";
import * as _130 from "./crisis/v1beta1/tx.amino";
import * as _131 from "./distribution/v1beta1/tx.amino";
import * as _132 from "./evidence/v1beta1/tx.amino";
import * as _133 from "./feegrant/v1beta1/tx.amino";
import * as _134 from "./gov/v1beta1/tx.amino";
import * as _135 from "./slashing/v1beta1/tx.amino";
import * as _136 from "./staking/v1beta1/tx.amino";
import * as _137 from "./vesting/v1beta1/tx.amino";
import * as _138 from "./authz/v1beta1/tx.registry";
import * as _139 from "./bank/v1beta1/tx.registry";
import * as _140 from "./crisis/v1beta1/tx.registry";
import * as _141 from "./distribution/v1beta1/tx.registry";
import * as _142 from "./evidence/v1beta1/tx.registry";
import * as _143 from "./feegrant/v1beta1/tx.registry";
import * as _144 from "./gov/v1beta1/tx.registry";
import * as _145 from "./slashing/v1beta1/tx.registry";
import * as _146 from "./staking/v1beta1/tx.registry";
import * as _147 from "./vesting/v1beta1/tx.registry";
import * as _148 from "./auth/v1beta1/query.rpc.Query";
import * as _149 from "./authz/v1beta1/query.rpc.Query";
import * as _150 from "./bank/v1beta1/query.rpc.Query";
import * as _151 from "./base/tendermint/v1beta1/query.rpc.Service";
import * as _152 from "./distribution/v1beta1/query.rpc.Query";
import * as _153 from "./evidence/v1beta1/query.rpc.Query";
import * as _154 from "./feegrant/v1beta1/query.rpc.Query";
import * as _155 from "./gov/v1beta1/query.rpc.Query";
import * as _156 from "./mint/v1beta1/query.rpc.Query";
import * as _157 from "./params/v1beta1/query.rpc.Query";
import * as _158 from "./slashing/v1beta1/query.rpc.Query";
import * as _159 from "./staking/v1beta1/query.rpc.Query";
import * as _160 from "./tx/v1beta1/service.rpc.Service";
import * as _161 from "./upgrade/v1beta1/query.rpc.Query";
import * as _162 from "./authz/v1beta1/tx.rpc.msg";
import * as _163 from "./bank/v1beta1/tx.rpc.msg";
import * as _164 from "./crisis/v1beta1/tx.rpc.msg";
import * as _165 from "./distribution/v1beta1/tx.rpc.msg";
import * as _166 from "./evidence/v1beta1/tx.rpc.msg";
import * as _167 from "./feegrant/v1beta1/tx.rpc.msg";
import * as _168 from "./gov/v1beta1/tx.rpc.msg";
import * as _169 from "./slashing/v1beta1/tx.rpc.msg";
import * as _170 from "./staking/v1beta1/tx.rpc.msg";
import * as _171 from "./vesting/v1beta1/tx.rpc.msg";
import * as _194 from "./rpc.query";
import * as _195 from "./rpc.tx";
export namespace cosmos {
  export namespace auth {
    export const v1beta1 = { ..._1,
      ..._2,
      ..._3,
      ..._148
    };
  }
  export namespace authz {
    export const v1beta1 = { ..._4,
      ..._5,
      ..._6,
      ..._7,
      ..._8,
      ..._128,
      ..._138,
      ..._149,
      ..._162
    };
  }
  export namespace bank {
    export const v1beta1 = { ..._9,
      ..._10,
      ..._11,
      ..._12,
      ..._13,
      ..._129,
      ..._139,
      ..._150,
      ..._163
    };
  }
  export namespace base {
    export namespace abci {
      export const v1beta1 = { ..._14
      };
    }
    export namespace kv {
      export const v1beta1 = { ..._15
      };
    }
    export namespace query {
      export const v1beta1 = { ..._16
      };
    }
    export namespace reflection {
      export const v1beta1 = { ..._17
      };
      export const v2alpha1 = { ..._18
      };
    }
    export namespace snapshots {
      export const v1beta1 = { ..._19
      };
    }
    export namespace store {
      export const v1beta1 = { ..._20,
        ..._21,
        ..._22
      };
    }
    export namespace tendermint {
      export const v1beta1 = { ..._23,
        ..._151
      };
    }
    export const v1beta1 = { ..._24
    };
  }
  export namespace capability {
    export const v1beta1 = { ..._25,
      ..._26
    };
  }
  export namespace crisis {
    export const v1beta1 = { ..._27,
      ..._28,
      ..._130,
      ..._140,
      ..._164
    };
  }
  export namespace crypto {
    export const ed25519 = { ..._29
    };
    export const multisig = { ..._30
    };
    export const secp256k1 = { ..._31
    };
    export const secp256r1 = { ..._32
    };
  }
  export namespace distribution {
    export const v1beta1 = { ..._33,
      ..._34,
      ..._35,
      ..._36,
      ..._131,
      ..._141,
      ..._152,
      ..._165
    };
  }
  export namespace evidence {
    export const v1beta1 = { ..._37,
      ..._38,
      ..._39,
      ..._40,
      ..._132,
      ..._142,
      ..._153,
      ..._166
    };
  }
  export namespace feegrant {
    export const v1beta1 = { ..._41,
      ..._42,
      ..._43,
      ..._44,
      ..._133,
      ..._143,
      ..._154,
      ..._167
    };
  }
  export namespace genutil {
    export const v1beta1 = { ..._45
    };
  }
  export namespace gov {
    export const v1beta1 = { ..._46,
      ..._47,
      ..._48,
      ..._49,
      ..._134,
      ..._144,
      ..._155,
      ..._168
    };
  }
  export namespace mint {
    export const v1beta1 = { ..._50,
      ..._51,
      ..._52,
      ..._156
    };
  }
  export namespace params {
    export const v1beta1 = { ..._53,
      ..._54,
      ..._157
    };
  }
  export namespace slashing {
    export const v1beta1 = { ..._55,
      ..._56,
      ..._57,
      ..._58,
      ..._135,
      ..._145,
      ..._158,
      ..._169
    };
  }
  export namespace staking {
    export const v1beta1 = { ..._59,
      ..._60,
      ..._61,
      ..._62,
      ..._63,
      ..._136,
      ..._146,
      ..._159,
      ..._170
    };
  }
  export namespace tx {
    export namespace signing {
      export const v1beta1 = { ..._64
      };
    }
    export const v1beta1 = { ..._65,
      ..._66,
      ..._160
    };
  }
  export namespace upgrade {
    export const v1beta1 = { ..._67,
      ..._68,
      ..._161
    };
  }
  export namespace vesting {
    export const v1beta1 = { ..._69,
      ..._70,
      ..._137,
      ..._147,
      ..._171
    };
  }
  export const ClientFactory = { ..._194,
    ..._195
  };
}