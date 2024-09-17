import { ENVTYPES } from './chains';

export * from './theme';
export * from './defaults';
export * from './chains';

export const env: string = process.env.NEXT_PUBLIC_CHAIN_ENV ?? ENVTYPES.PROD;