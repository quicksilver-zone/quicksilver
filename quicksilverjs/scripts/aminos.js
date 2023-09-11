export const AMINO_MAP = {
  // PUT YOUR AMINO names here...
  '/quicksilver.participationrewards.v1.MsgSubmitClaim': {
    aminoType: 'quicksilver/MsgSubmitClaim'
  },
  '/quicksilver.airdrop.v1.MsgClaim': {
    aminoType: 'quicksilver/MsgClaim'
  },
  '/quicksilver.interchainstaking.v1.MsgRequestRedemption': {
    aminoType: 'quicksilver/MsgRequestRedemption'
  },
  '/quicksilver.interchainstaking.v1.MsgSignalIntent': {
    aminoType: 'quicksilver/MsgSignalIntent'
  },

  // Staking
  '/cosmos.staking.v1beta1.MsgCreateValidator': {
    aminoType: 'cosmos-sdk/MsgCreateValidator'
  },
  '/cosmos.staking.v1beta1.MsgEditValidator': {
    aminoType: 'cosmos-sdk/MsgEditValidator'
  },
  '/cosmos.staking.v1beta1.MsgDelegate': {
    aminoType: 'cosmos-sdk/MsgDelegate'
  },
  '/cosmos.staking.v1beta1.MsgUndelegate': {
    aminoType: 'cosmos-sdk/MsgUndelegate'
  },
  '/cosmos.staking.v1beta1.MsgBeginRedelegate': {
    aminoType: 'cosmos-sdk/MsgBeginRedelegate'
  },
  '/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation': {
    aminoType: 'cosmos-sdk/MsgCancelUnbondingDelegation'
  },
  '/cosmos.staking.v1beta1.MsgUpdateParams': {
    aminoType: 'cosmos-sdk/x/staking/MsgUpdateParams'
  },

  // IBC
  '/ibc.applications.transfer.v1.MsgTransfer': {
    aminoType: 'cosmos-sdk/MsgTransfer'
  }
};
