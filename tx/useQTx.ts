import { quicksilver } from 'quicksilverjs';

const { requestRedemption, signalIntent } =
  quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
