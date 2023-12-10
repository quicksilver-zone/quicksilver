import { render, screen, waitFor } from '@testing-library/react';
//@ts-ignore
import { test, expect, mock } from 'bun:test';
import { useState } from 'react';

import { NetworkSelect } from '@/components';

import Staking from '../pages/staking';

const networks = [
  {
    value: 'ATOM',
    logo: '/quicksilver-app-v2/img/networks/atom.svg',
    qlogo: '/quicksilver-app-v2/img/networks/q-atom.svg',
    name: 'Cosmos Hub',
    chainName: 'cosmoshub',
    chainId: 'cosmoshub-4',
  },
  {
    value: 'OSMO',
    logo: '/quicksilver-app-v2/img/networks/osmosis.svg',
    qlogo: '/quicksilver-app-v2/img/networks/qosmo.svg',
    name: 'Osmosis',
    chainName: 'osmosis',
    chainId: 'osmosis-1',
  },
  {
    value: 'STARS',
    logo: '/quicksilver-app-v2/img/networks/stargaze.svg',
    qlogo: '/quicksilver-app-v2/img/networks/stargaze-2.png',
    name: 'Stargaze',
    chainName: 'stargaze',
    chainId: 'stargaze-1',
  },
  {
    value: 'REGEN',
    logo: '/quicksilver-app-v2/img/networks/regen.svg',
    qlogo: '/quicksilver-app-v2/img/networks/regen.svg',
    name: 'Regen',
    chainName: 'regen',
    chainId: 'regen-1',
  },
  {
    value: 'SOMM',
    logo: '/quicksilver-app-v2/img/networks/sommelier.png',
    qlogo: '/quicksilver-app-v2/img/networks/sommelier.png',
    name: 'Sommelier',
    chainName: 'sommelier',
    chainId: 'sommelier-3',
  },
];

function MockNetwork() {
  const [selectedNetwork, setSelectedNetwork] = useState(networks[0]);

  const DummyNetworkSelect = mock(() => {
    return <NetworkSelect selectedOption={selectedNetwork} setSelectedNetwork={setSelectedNetwork} />;
  });
  return DummyNetworkSelect();
}

test('Staking Page renders without crashing', async () => {
  render(
    <>
      <MockNetwork />
      <Staking />
    </>,
  );

  await waitFor(() => {
    expect(screen.getByText('NetworkSelect')).pass();
  });
});
