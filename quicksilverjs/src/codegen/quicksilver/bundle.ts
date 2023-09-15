import * as _177 from "./airdrop/v1/airdrop";
import * as _178 from "./airdrop/v1/genesis";
import * as _179 from "./airdrop/v1/messages";
import * as _180 from "./airdrop/v1/params";
import * as _181 from "./airdrop/v1/proposals";
import * as _182 from "./airdrop/v1/query";
import * as _183 from "./claimsmanager/v1/claimsmanager";
import * as _184 from "./claimsmanager/v1/genesis";
import * as _185 from "./claimsmanager/v1/messages";
import * as _186 from "./claimsmanager/v1/query";
import * as _187 from "./epochs/v1/genesis";
import * as _188 from "./epochs/v1/query";
import * as _189 from "./interchainquery/v1/genesis";
import * as _190 from "./interchainquery/v1/interchainquery";
import * as _191 from "./interchainquery/v1/messages";
import * as _192 from "./interchainquery/v1/query";
import * as _193 from "./interchainstaking/v1/genesis";
import * as _194 from "./interchainstaking/v1/interchainstaking";
import * as _195 from "./interchainstaking/v1/messages";
import * as _196 from "./interchainstaking/v1/proposals";
import * as _197 from "./interchainstaking/v1/query";
import * as _198 from "./participationrewards/v1/proposals";
import * as _199 from "./mint/v1beta1/genesis";
import * as _200 from "./mint/v1beta1/mint";
import * as _201 from "./mint/v1beta1/query";
import * as _202 from "./participationrewards/v1/genesis";
import * as _203 from "./participationrewards/v1/messages";
import * as _204 from "./participationrewards/v1/participationrewards";
import * as _205 from "./participationrewards/v1/query";
import * as _206 from "./tokenfactory/v1beta1/authorityMetadata";
import * as _207 from "./tokenfactory/v1beta1/genesis";
import * as _208 from "./tokenfactory/v1beta1/params";
import * as _209 from "./tokenfactory/v1beta1/query";
import * as _210 from "./tokenfactory/v1beta1/tx";
import * as _325 from "./airdrop/v1/messages.amino";
import * as _326 from "./interchainquery/v1/messages.amino";
import * as _327 from "./interchainstaking/v1/messages.amino";
import * as _328 from "./participationrewards/v1/messages.amino";
import * as _329 from "./tokenfactory/v1beta1/tx.amino";
import * as _330 from "./airdrop/v1/messages.registry";
import * as _331 from "./interchainquery/v1/messages.registry";
import * as _332 from "./interchainstaking/v1/messages.registry";
import * as _333 from "./participationrewards/v1/messages.registry";
import * as _334 from "./tokenfactory/v1beta1/tx.registry";
import * as _335 from "./airdrop/v1/query.lcd";
import * as _336 from "./claimsmanager/v1/query.lcd";
import * as _337 from "./epochs/v1/query.lcd";
import * as _338 from "./interchainstaking/v1/query.lcd";
import * as _339 from "./mint/v1beta1/query.lcd";
import * as _340 from "./participationrewards/v1/query.lcd";
import * as _341 from "./tokenfactory/v1beta1/query.lcd";
import * as _342 from "./airdrop/v1/query.rpc.Query";
import * as _343 from "./claimsmanager/v1/query.rpc.Query";
import * as _344 from "./epochs/v1/query.rpc.Query";
import * as _345 from "./interchainstaking/v1/query.rpc.Query";
import * as _346 from "./mint/v1beta1/query.rpc.Query";
import * as _347 from "./participationrewards/v1/query.rpc.Query";
import * as _348 from "./tokenfactory/v1beta1/query.rpc.Query";
import * as _349 from "./airdrop/v1/messages.rpc.msg";
import * as _350 from "./interchainquery/v1/messages.rpc.msg";
import * as _351 from "./interchainstaking/v1/messages.rpc.msg";
import * as _352 from "./participationrewards/v1/messages.rpc.msg";
import * as _353 from "./tokenfactory/v1beta1/tx.rpc.msg";
import * as _363 from "./lcd";
import * as _364 from "./rpc.query";
export namespace quicksilver {
  export namespace airdrop {
    export const v1 = {
      ..._177,
      ..._178,
      ..._179,
      ..._180,
      ..._181,
      ..._182,
      ..._325,
      ..._330,
      ..._335,
      ..._342,
      ..._349
    };
  }
  export namespace claimsmanager {
    export const v1 = {
      ..._183,
      ..._184,
      ..._185,
      ..._186,
      ..._336,
      ..._343
    };
  }
  export namespace epochs {
    export const v1 = {
      ..._187,
      ..._188,
      ..._337,
      ..._344
    };
  }
  export namespace interchainquery {
    export const v1 = {
      ..._189,
      ..._190,
      ..._191,
      ..._192,
      ..._326,
      ..._331,
      ..._350
    };
  }
  export namespace interchainstaking {
    export const v1 = {
      ..._193,
      ..._194,
      ..._195,
      ..._196,
      ..._197,
      ..._198,
      ..._327,
      ..._332,
      ..._338,
      ..._345,
      ..._351
    };
  }
  export namespace mint {
    export const v1beta1 = {
      ..._199,
      ..._200,
      ..._201,
      ..._339,
      ..._346
    };
  }
  export namespace participationrewards {
    export const v1 = {
      ..._202,
      ..._203,
      ..._204,
      ..._205,
      ..._328,
      ..._333,
      ..._340,
      ..._347,
      ..._352
    };
  }
  export namespace tokenfactory {
    export const v1beta1 = {
      ..._206,
      ..._207,
      ..._208,
      ..._209,
      ..._210,
      ..._329,
      ..._334,
      ..._341,
      ..._348,
      ..._353
    };
  }
  export const ClientFactory = {
    ..._363,
    ..._364,
  };
}