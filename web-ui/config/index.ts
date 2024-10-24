import { ENVTYPES, local_chain } from './chains';

export * from './theme';
export * from './chains';

export const env: string = process.env.NEXT_PUBLIC_CHAIN_ENV ?? ENVTYPES.PROD;

export const defaultChainName: string = local_chain.get(env)?.chain_name ?? '';