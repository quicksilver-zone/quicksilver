import { render, screen, waitFor } from '@testing-library/react';
//@ts-ignore
import { test, expect, mock } from 'bun:test';
import { useState } from 'react';

import Staking from '../pages/staking';

import { NetworkSelect } from '@/components';

const networks = [
  {
    value: 'ATOM',
    logo: '/img/networks/atom.svg',
    qlogo: '/img/networks/q-atom.svg',
    name: 'Cosmos Hub',
    chainName: 'cosmoshub',
    chainId: 'cosmoshub-4',
  },
  {
    value: 'OSMO',
    logo: '/img/networks/osmosis.svg',
    qlogo: '/img/networks/qosmo.svg',
    name: 'Osmosis',
    chainName: 'osmosis',
    chainId: 'osmosis-1',
  },
  {
    value: 'STARS',
    logo: '/img/networks/stargaze.svg',
    qlogo: '/img/networks/stargaze-2.png',
    name: 'Stargaze',
    chainName: 'stargaze',
    chainId: 'stargaze-1',
  },
  {
    value: 'REGEN',
    logo: '/img/networks/regen.svg',
    qlogo: '/img/networks/regen.svg',
    name: 'Regen',
    chainName: 'regen',
    chainId: 'regen-1',
  },
  {
    value: 'SOMM',
    logo: '/img/networks/sommelier.png',
    qlogo: '/img/networks/sommelier.png',
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
