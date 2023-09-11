import { OfflineSigner, GeneratedType, Registry } from "@cosmjs/proto-signing";
import { defaultRegistryTypes, AminoTypes, SigningStargateClient } from "@cosmjs/stargate";
import * as quicksilverAirdropV1MessagesRegistry from "./airdrop/v1/messages.registry";
import * as quicksilverInterchainqueryV1MessagesRegistry from "./interchainquery/v1/messages.registry";
import * as quicksilverInterchainstakingV1MessagesRegistry from "./interchainstaking/v1/messages.registry";
import * as quicksilverParticipationrewardsV1MessagesRegistry from "./participationrewards/v1/messages.registry";
import * as quicksilverTokenfactoryV1beta1TxRegistry from "./tokenfactory/v1beta1/tx.registry";
import * as quicksilverAirdropV1MessagesAmino from "./airdrop/v1/messages.amino";
import * as quicksilverInterchainqueryV1MessagesAmino from "./interchainquery/v1/messages.amino";
import * as quicksilverInterchainstakingV1MessagesAmino from "./interchainstaking/v1/messages.amino";
import * as quicksilverParticipationrewardsV1MessagesAmino from "./participationrewards/v1/messages.amino";
import * as quicksilverTokenfactoryV1beta1TxAmino from "./tokenfactory/v1beta1/tx.amino";
export const quicksilverAminoConverters = { ...quicksilverAirdropV1MessagesAmino.AminoConverter,
  ...quicksilverInterchainqueryV1MessagesAmino.AminoConverter,
  ...quicksilverInterchainstakingV1MessagesAmino.AminoConverter,
  ...quicksilverParticipationrewardsV1MessagesAmino.AminoConverter,
  ...quicksilverTokenfactoryV1beta1TxAmino.AminoConverter
};
export const quicksilverProtoRegistry: ReadonlyArray<[string, GeneratedType]> = [...quicksilverAirdropV1MessagesRegistry.registry, ...quicksilverInterchainqueryV1MessagesRegistry.registry, ...quicksilverInterchainstakingV1MessagesRegistry.registry, ...quicksilverParticipationrewardsV1MessagesRegistry.registry, ...quicksilverTokenfactoryV1beta1TxRegistry.registry];
export const getSigningQuicksilverClientOptions = ({
  defaultTypes = defaultRegistryTypes
}: {
  defaultTypes?: ReadonlyArray<[string, GeneratedType]>;
} = {}): {
  registry: Registry;
  aminoTypes: AminoTypes;
} => {
  const registry = new Registry([...defaultTypes, ...quicksilverProtoRegistry]);
  const aminoTypes = new AminoTypes({ ...quicksilverAminoConverters
  });
  return {
    registry,
    aminoTypes
  };
};
export const getSigningQuicksilverClient = async ({
  rpcEndpoint,
  signer,
  defaultTypes = defaultRegistryTypes
}: {
  rpcEndpoint: string;
  signer: OfflineSigner;
  defaultTypes?: ReadonlyArray<[string, GeneratedType]>;
}) => {
  const {
    registry,
    aminoTypes
  } = getSigningQuicksilverClientOptions({
    defaultTypes
  });
  const client = await SigningStargateClient.connectWithSigner(rpcEndpoint, signer, {
    registry,
    aminoTypes
  });
  return client;
};