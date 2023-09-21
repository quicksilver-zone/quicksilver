import { render, screen, waitFor } from '@testing-library/react';
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
  },
  {
    value: 'OSMO',
    logo: '/quicksilver-app-v2/img/networks/osmosis.svg',
    qlogo: '/quicksilver-app-v2/img/networks/qosmo.svg',
    name: 'Osmosis',
    chainName: 'osmosis',
  },
  {
    value: 'STARS',
    logo: '/quicksilver-app-v2/img/networks/stargaze.svg',
    qlogo: '/quicksilver-app-v2/img/networks/stargaze-2.png',
    name: 'Stargaze',
    chainName: 'stargaze',
  },
  {
    value: 'REGEN',
    logo: '/quicksilver-app-v2/img/networks/regen.svg',
    qlogo: '/quicksilver-app-v2/img/networks/regen.svg',
    name: 'Regen',
    chainName: 'regen',
  },
  {
    value: 'SOMM',
    logo: '/quicksilver-app-v2/img/networks/sommelier.png',
    qlogo: '/quicksilver-app-v2/img/networks/sommelier.png',
    name: 'Sommelier',
    chainName: 'sommelier',
  },
];

function MockNetwork() {
  const [selectedNetwork, setSelectedNetwork] = useState(networks[0]);

  const DummyNetworkSelect = mock(() => {
    return (
      <NetworkSelect
        selectedOption={selectedNetwork}
        setSelectedNetwork={setSelectedNetwork}
      />
    );
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
    expect(screen.getByText('NetworkSelect')).toBeInTheDocument();
  });
});
